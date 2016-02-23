package tcp

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ardanlabs/kit/log"
	"github.com/ardanlabs/kit/pool"
)

// Set of error variables for start up.
var (
	ErrInvalidConfiguration     = errors.New("Invalid Configuration")
	ErrInvalidNetType           = errors.New("Invalid NetType Configuration")
	ErrInvalidConnHandler       = errors.New("Invalid Connection Handler Configuration")
	ErrInvalidReqHandler        = errors.New("Invalid Request Handler Configuration")
	ErrInvalidRespHandler       = errors.New("Invalid Response Handler Configuration")
	ErrInvalidPoolConfiguration = errors.New("Invalid Pool Configuration")
)

//==============================================================================

// TCP contains a set of networked client connections.
type TCP struct {
	Config
	Name string

	ipAddress string
	port      int
	tcpAddr   *net.TCPAddr

	listener   *net.TCPListener
	listenerMu sync.Mutex

	clients   map[string]*client
	clientsMu sync.Mutex

	recv      *pool.Pool
	send      *pool.Pool
	userPools bool

	wg sync.WaitGroup

	dropConns    int32
	shuttingDown int32

	lastAcceptedConnection time.Time
}

// New creates a new manager to service clients.
func New(context interface{}, name string, cfg Config) (*TCP, error) {
	log.Dev(context, "New", "Started : Name[%s] NetType[%s] Addr[%s] RecvMaxPoolSize[%d] SendMaxPoolSize[%d]", name, cfg.NetType, cfg.Addr, cfg.RecvMaxPoolSize, cfg.SendMaxPoolSize)

	// Validate the configuration.
	if err := cfg.Validate(); err != nil {
		log.Error(context, "New", err, "Completed")
		return nil, err
	}

	// Resolve the addr that is provided.
	tcpAddr, err := net.ResolveTCPAddr(cfg.NetType, cfg.Addr)
	if err != nil {
		log.Error(context, "New", err, "Completed")
		return nil, err
	}

	log.Dev(context, "New", "Address[ %s ] Zone[%s]", join(tcpAddr.IP.String(), tcpAddr.Port), tcpAddr.Zone)

	// Need a work pool to handle the received messages.
	var recv *pool.Pool
	if cfg.RecvPool != nil {
		recv = cfg.RecvPool
	} else {
		recvCfg := pool.Config{
			MinRoutines: cfg.RecvMinPoolSize,
			MaxRoutines: cfg.RecvMaxPoolSize,
		}

		var err error
		if recv, err = pool.New(context, name+"-Recv", recvCfg); err != nil {
			log.Error(context, "New", err, "Completed")
			return nil, err
		}
	}

	// Need a work pool to handle the messages to send.
	var send *pool.Pool
	if cfg.SendPool != nil {
		send = cfg.SendPool
	} else {
		sendCfg := pool.Config{
			MinRoutines: cfg.SendMinPoolSize,
			MaxRoutines: cfg.SendMaxPoolSize,
		}

		var err error
		if send, err = pool.New(context, name+"-Send", sendCfg); err != nil {
			log.Error(context, "New", err, "Completed")
			return nil, err
		}
	}

	// Are we using user provided work pools. Validation is helping us
	// only have to check one of the two configuration options for this.
	var userPools bool
	if cfg.RecvPool != nil {
		log.Dev(context, "New", "Using User Pools")
		userPools = true
	}

	// Create a TCP for this ipaddress and port.
	t := TCP{
		Config: cfg,
		Name:   name,

		ipAddress: tcpAddr.IP.String(),
		port:      tcpAddr.Port,
		tcpAddr:   tcpAddr,

		clients: make(map[string]*client),

		recv:      recv,
		send:      send,
		userPools: userPools,
	}

	log.Dev(context, "New", "Completed")
	return &t, nil
}

// join takes an IP and port values and creates a cleaner string.
func join(ip string, port int) string {
	return net.JoinHostPort(ip, strconv.Itoa(port))
}

// Start creates the accept routine and begins to accept connections.
func (t *TCP) Start(context interface{}) error {
	log.Dev(context, "Start", "IPAddress[ %s ]", join(t.ipAddress, t.port))

	t.listenerMu.Lock()
	{
		// If the listener has been started already, return an error.
		if t.listener != nil {
			err := errors.New("This TCP has already been started")
			log.Error(context, "Start", err, "Completed")
			t.listenerMu.Unlock()
			return err
		}
	}
	t.listenerMu.Unlock()

	t.wg.Add(1)

	// We need to wait for the goroutine to initialize itself.
	var waitStart sync.WaitGroup
	waitStart.Add(1)

	// Start the connection accept routine.
	go func() {
		var listener *net.TCPListener

		for {
			t.listenerMu.Lock()
			{
				// Start a listener for the specified addr and port is one
				// does not exist.
				if t.listener == nil {
					var err error
					listener, err = net.ListenTCP(t.NetType, t.tcpAddr)
					if err != nil {
						log.Error(context, "Start", err, "Completed")
						panic(err)
					}

					t.listener = listener

					waitStart.Done()

					log.Dev(context, "accept-routine", "Waiting For Connections : IPAddress[ %s ]", join(t.ipAddress, t.port))
				}
			}
			t.listenerMu.Unlock()

			// Listen for new connections.
			conn, err := listener.Accept()
			if err != nil {
				shutdown := atomic.LoadInt32(&t.shuttingDown)

				if shutdown == 0 {
					log.Error(context, "accept-routine", err, "Completed")
				} else {
					t.listenerMu.Lock()
					{
						t.listener = nil
					}
					t.listenerMu.Unlock()
					break
				}

				// temporary is declared to test for the existence of
				// the method coming from the net package.
				type temporary interface {
					Temporary() bool
				}

				if e, ok := err.(temporary); ok && !e.Temporary() {
					t.listenerMu.Lock()
					{
						t.listener.Close()
						t.listener = nil
					}
					t.listenerMu.Unlock()

					// Don't want to add a flag. So setting this back to
					// 1 so when the listener is re-established, the call
					// to Done does not fail.
					waitStart.Add(1)
				}

				continue
			}

			// Check if we are being asked to drop all new connections.
			if drop := atomic.LoadInt32(&t.dropConns); drop == 1 {
				log.Dev(context, "accept-routine", "*******> DROPPING CONNECTION")
				conn.Close()
				continue
			}

			// Check if rate limit is enabled.
			if t.RateLimit != nil {
				now := time.Now()

				// We will only accept 1 connection per duration. Anything
				// connection above that must be dropped.
				if t.lastAcceptedConnection.Add(t.RateLimit()).After(now) {
					log.Dev(context, "accept-routine", "*******> DROPPING CONNECTION Local[ %v ] Remote[ %v ] DUE TO RATE LIMIT %v", conn.LocalAddr(), conn.RemoteAddr(), t.RateLimit())
					conn.Close()
					continue
				}

				// Since we accepted connection, mark the time.
				t.lastAcceptedConnection = now
			}

			// Add this new connection to the manager map.
			t.join(context, conn)
		}

		// Shutting down the routine.
		t.wg.Done()
		log.Dev(context, "accept-routine", "Shutdown : IPAddress[ %s ]", join(t.ipAddress, t.port))
	}()

	// Wait for the goroutine to initialize itself.
	waitStart.Wait()

	log.Dev(context, "Start", "Completed")
	return nil
}

// Stop shuts down the manager and closes all connections.
func (t *TCP) Stop(context interface{}) error {
	log.Dev(context, "Stop", "Started : IPAddress[ %s ]", join(t.ipAddress, t.port))

	t.listenerMu.Lock()
	{
		// If the listener has been stopped already, return an error.
		if t.listener == nil {
			err := errors.New("This TCP has already been stopped")
			log.Error(context, "Stop", err, "Completed")
			t.listenerMu.Unlock()
			return err
		}
	}
	t.listenerMu.Unlock()

	// Mark that we are shutting down.
	atomic.StoreInt32(&t.shuttingDown, 1)

	// Don't accept anymore client connections.
	t.listenerMu.Lock()
	{
		t.listener.Close()
	}
	t.listenerMu.Unlock()

	// Stop processing all the work.
	if !t.userPools {
		t.recv.Shutdown(context)
		t.send.Shutdown(context)
	} else {
		log.Dev(context, "Stop", "User pools being used, not shutting them down.")
	}

	// Make a copy of all the connections. We need to do this
	// since we have to lock the map to read it. Dropping a
	// connection requires locks as well.
	var clients map[string]*client
	t.clientsMu.Lock()
	{
		clients = make(map[string]*client)
		for k, v := range t.clients {
			clients[k] = v
		}
	}
	t.clientsMu.Unlock()

	// Drop all the existing connections.
	for _, c := range clients {
		// This waits for each routine to terminate.
		c.drop()
	}

	// Wait for the accept routine to terminate.
	t.wg.Wait()

	log.Dev(context, "Stop", "Completed")
	return nil
}

// Do will post the request to be sent by the client worker pool.
func (t *TCP) Do(context interface{}, r *Response) error {
	log.Dev(context, "Do", "Started : Local[ %s ] Remote[ %s ]", join(t.ipAddress, t.port), r.TCPAddr.String())

	// Find the client connection for this IPAddress.
	var c *client
	t.clientsMu.Lock()
	{
		// If this ipaddress and socket does not exist, report an error.
		var ok bool
		if c, ok = t.clients[r.TCPAddr.String()]; !ok {
			err := fmt.Errorf("IP Address disconnected [ %s ]", r.TCPAddr.String())
			log.Error(context, "Do", err, "Completed")

			t.clientsMu.Unlock()
			return err
		}
	}
	t.clientsMu.Unlock()

	// Set the unexported fields.
	r.tcp = t
	r.client = c
	r.context = context

	// Send this to the client work pool for processing.
	t.send.Do(context, r)

	log.Dev(context, "Do", "Completed")
	return nil
}

// DropConnections sets a flag to tell the accept routine to immediately
// drop connections that come in.
func (t *TCP) DropConnections(context interface{}, drop bool) {
	log.Dev(context, "DropConnections", "Started : drop[%v]", drop)

	if drop {
		atomic.StoreInt32(&t.dropConns, 1)
		return
	}

	atomic.StoreInt32(&t.dropConns, 0)
	log.Dev(context, "DropConnections", "Completed")
}

// StatsRecv returns the current snapshot of the recv pool stats.
func (t *TCP) StatsRecv() pool.Stat {
	return t.recv.Stats()
}

// StatsSend returns the current snapshot of the send pool stats.
func (t *TCP) StatsSend() pool.Stat {
	return t.send.Stats()
}

// Addr returns the listener's network address. This may be different than the values
// provided in the configuration, for example if configuration port value is 0.
func (t *TCP) Addr() net.Addr {
	// We are aware this read is not safe with the
	// goroutine accepting connections.
	if t.listener == nil {
		return nil
	}
	return t.listener.Addr()
}

// join takes a new connection and adds it to the manager.
func (t *TCP) join(context interface{}, conn net.Conn) {
	ipAddress := conn.RemoteAddr().String()
	cntx := fmt.Sprintf("%s-%s", context, ipAddress)
	log.Dev(cntx, "join", "Started : Remote IPAddress[ %s ], Local IPAddress[ %v ]", ipAddress, conn.LocalAddr())

	t.clientsMu.Lock()
	{
		// If this ipaddress and socket alread exist, we have a problet.
		if _, ok := t.clients[ipAddress]; ok {
			err := fmt.Errorf("IP Address already connected [ %s ]", ipAddress)
			log.Error(context, "Do", err, "Completed")
			conn.Close()

			t.clientsMu.Unlock()
			return
		}

		// Add the new client connection.
		t.clients[ipAddress] = newClient(cntx, t, conn)
	}
	t.clientsMu.Unlock()

	log.Dev(context, "join", "Completed")
}

// remove deletes a connection from the manager.
func (t *TCP) remove(context interface{}, conn net.Conn) {
	ipAddress := conn.RemoteAddr().String()
	log.Dev(context, "remove", "Started : IPAddress[ %s ]", ipAddress)

	t.clientsMu.Lock()
	{
		// If this ipaddress and socket does not exist, we have a probler.
		if _, ok := t.clients[ipAddress]; !ok {
			err := fmt.Errorf("IP Address already removed [ %s ]", ipAddress)
			log.Error(context, "remove", err, "Completed")

			t.clientsMu.Unlock()
			return
		}

		// Remove the client connection from the map.
		delete(t.clients, ipAddress)
	}
	t.clientsMu.Unlock()

	// Close the connection for safe keeping.
	conn.Close()

	log.Dev(context, "remove", "Completed")
}
