package tcp

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

// Set of error variables for start up.
var (
	ErrInvalidConfiguration = errors.New("invalid configuration")
	ErrInvalidNetType       = errors.New("invalid net type configuration")
	ErrInvalidConnHandler   = errors.New("invalid connection handler configuration")
	ErrInvalidReqHandler    = errors.New("invalid request handler configuration")
	ErrInvalidRespHandler   = errors.New("invalid response handler configuration")
)

// Set of event types.
const (
	EvtAccept = iota + 1
	EvtJoin
	EvtRead
	EvtRemove
	EvtDrop
	EvtGroom
)

// Set of event sub types.
const (
	TypError = iota + 1
	TypInfo
	TypTrigger
)

// CltError provides support for multi client operations that might error.
type CltError []error

// Error implments the error interface for CltError.
func (ce CltError) Error() string {
	var b bytes.Buffer
	for _, err := range ce {
		b.WriteString(err.Error())
		b.WriteString("\n")
	}
	return b.String()
}

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

	wg sync.WaitGroup

	dropConns    int32
	shuttingDown int32

	lastAcceptedConnection time.Time
}

// New creates a new manager to service clients.
func New(name string, cfg Config) (*TCP, error) {

	// Validate the configuration.
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	// Resolve the addr that is provided.
	tcpAddr, err := net.ResolveTCPAddr(cfg.NetType, cfg.Addr)
	if err != nil {
		return nil, err
	}

	// Create a TCP for this ipaddress and port.
	t := TCP{
		Config: cfg,
		Name:   name,

		ipAddress: tcpAddr.IP.String(),
		port:      tcpAddr.Port,
		tcpAddr:   tcpAddr,

		clients: make(map[string]*client),
	}

	return &t, nil
}

// join takes an IP and port values and creates a cleaner string.
func join(ip string, port int) string {
	return net.JoinHostPort(ip, strconv.Itoa(port))
}

// Start creates the accept routine and begins to accept connections.
func (t *TCP) Start() error {
	t.listenerMu.Lock()
	{
		// If the listener has been started already, return an error.
		if t.listener != nil {
			t.listenerMu.Unlock()
			return errors.New("this TCP has already been started")
		}
	}
	t.listenerMu.Unlock()

	// We need to wait for the goroutine we are about to
	// create to initialize itself.
	var waitStart sync.WaitGroup
	waitStart.Add(1)

	// Start the connection accept routine.
	t.wg.Add(1)
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
						panic(err)
					}

					t.listener = listener
					waitStart.Done()

					t.Event(EvtAccept, TypInfo, join(t.ipAddress, t.port), "waiting")
				}
			}
			t.listenerMu.Unlock()

			// Listen for new connections.
			conn, err := listener.Accept()
			if err != nil {
				shutdown := atomic.LoadInt32(&t.shuttingDown)

				if shutdown == 0 {
					t.Event(EvtAccept, TypError, conn.RemoteAddr().String(), err.Error())
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
				t.Event(EvtAccept, TypInfo, "", "dropping new connection")
				conn.Close()
				continue
			}

			// Check if rate limit is enabled.
			if t.RateLimit != nil {
				now := time.Now().UTC()

				// We will only accept 1 connection per duration. Anything
				// connection above that must be dropped.
				if t.lastAcceptedConnection.Add(t.RateLimit()).After(now) {
					t.Event(EvtAccept, TypError, conn.RemoteAddr().String(), "rate limit drop : Local[ %v ] Limit[ %v ]", conn.LocalAddr(), t.RateLimit())
					conn.Close()
					continue
				}

				// Since we accepted connection, mark the time.
				t.lastAcceptedConnection = now
			}

			// Add this new connection to the manager map.
			t.join(conn)
		}

		// Shutting down the routine.
		t.wg.Done()
		t.Event(EvtAccept, TypError, join(t.ipAddress, t.port), "shutdown")
	}()

	// Wait for the goroutine to initialize itself.
	waitStart.Wait()

	return nil
}

// Stop shuts down the manager and closes all connections.
func (t *TCP) Stop() error {
	t.listenerMu.Lock()
	{
		// If the listener has been stopped already, return an error.
		if t.listener == nil {
			t.listenerMu.Unlock()
			return errors.New("this TCP has already been stopped")
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

	return nil
}

// Drop will close the socket connection.
func (t *TCP) Drop(tcpAddr *net.TCPAddr) error {

	// Find the client connection for this IPAddress.
	var c *client
	t.clientsMu.Lock()
	{
		// Validate this ipaddress and socket exists first.
		var ok bool
		if c, ok = t.clients[tcpAddr.String()]; !ok {
			t.clientsMu.Unlock()
			return fmt.Errorf("IP[ %s ] : disconnected", tcpAddr.String())
		}
	}
	t.clientsMu.Unlock()

	// Drop the connection using a goroutine since we are on the
	// socket goroutine most likely.
	go c.drop()
	return nil
}

// Send will deliver the response back to the client.
func (t *TCP) Send(ctx context.Context, r *Response) error {

	// Find the client connection for this IPAddress.
	var c *client
	t.clientsMu.Lock()
	{
		// Validate this ipaddress and socket exists first.
		var ok bool
		if c, ok = t.clients[r.TCPAddr.String()]; !ok {
			t.clientsMu.Unlock()
			return fmt.Errorf("IP[ %s ] : disconnected", r.TCPAddr.String())
		}

		// Increment the number of writes.
		c.nWrites++
	}
	t.clientsMu.Unlock()

	// Send the response.
	return t.RespHandler.Write(r, c.writer)
}

// SendAll will deliver the response back to all connected clients.
func (t *TCP) SendAll(ctx context.Context, r *Response) error {
	var clts []*client
	t.clientsMu.Lock()
	{
		for _, c := range t.clients {
			clts = append(clts, c)
			c.nWrites++
		}
	}
	t.clientsMu.Unlock()

	// TODO: Consider doing this in parallel.
	var errors CltError
	for _, c := range clts {
		if err := t.RespHandler.Write(r, c.writer); err != nil {
			errors = append(errors, err)
		}
	}

	if errors != nil {
		return errors
	}
	return nil
}

// DropConnections sets a flag to tell the accept routine to immediately
// drop connections that come in.
func (t *TCP) DropConnections(drop bool) {
	if drop {
		atomic.StoreInt32(&t.dropConns, 1)
		return
	}

	atomic.StoreInt32(&t.dropConns, 0)
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

// Connections returns the number of client connections.
func (t *TCP) Connections() int {
	var l int

	t.clientsMu.Lock()
	{
		l = len(t.clients)
	}
	t.clientsMu.Unlock()

	return l
}

// Stat represents a client statistic.
type Stat struct {
	IP       string
	Reads    int
	Writes   int
	TimeConn time.Time
	LastAct  time.Time
}

// ClientStats return details for all active clients.
func (t *TCP) ClientStats() []Stat {
	var clts []*client
	t.clientsMu.Lock()
	{
		for _, v := range t.clients {
			clts = append(clts, v)
		}
	}
	t.clientsMu.Unlock()

	stats := make([]Stat, len(clts))
	for i, c := range clts {
		stats[i] = Stat{
			IP:       c.ipAddress,
			Reads:    c.nReads,
			Writes:   c.nWrites,
			TimeConn: c.timeConn,
			LastAct:  c.lastAct,
		}
	}

	return stats
}

// Clients returns the number of active clients connected.
func (t *TCP) Clients() int {
	var count int
	t.clientsMu.Lock()
	{
		count = len(t.clients)
	}
	t.clientsMu.Unlock()

	return count
}

// Groom drops connections that are not active for the specified duration.
func (t *TCP) Groom(d time.Duration) {
	var clts []*client
	t.clientsMu.Lock()
	{
		for _, v := range t.clients {
			clts = append(clts, v)
		}
	}
	t.clientsMu.Unlock()

	now := time.Now().UTC()
	for _, c := range clts {
		sub := now.Sub(c.lastAct)
		if sub >= d {

			// TODO
			// This is a blocking call that waits for the socket goroutine
			// to report its done. This parallel call should work well since
			// there is no error handling needed.
			t.Event(EvtGroom, TypInfo, c.ipAddress, "Last[ %v ] Dur[ %v ]", c.lastAct.Format(time.RFC3339), sub)
			go c.drop()
		}
	}
}

// join takes a new connection and adds it to the manager.
func (t *TCP) join(conn net.Conn) {
	ipAddress := conn.RemoteAddr().String()
	t.Event(EvtJoin, TypTrigger, ipAddress, "new connection")

	t.clientsMu.Lock()
	{
		// Validate this has not been joined already.
		if _, ok := t.clients[ipAddress]; ok {
			t.Event(EvtJoin, TypError, ipAddress, "already connected")
			conn.Close()

			t.clientsMu.Unlock()
			return
		}

		// Add the client connection to the map.
		t.clients[ipAddress] = newClient(t, conn)
	}
	t.clientsMu.Unlock()
}

// remove deletes a connection from the manager.
func (t *TCP) remove(conn net.Conn) {
	ipAddress := conn.RemoteAddr().String()

	t.clientsMu.Lock()
	{
		// Validate this has not been removed already.
		if _, ok := t.clients[ipAddress]; !ok {
			t.Event(EvtRemove, TypError, ipAddress, "already removed")
			t.clientsMu.Unlock()
			return
		}

		// Remove the client connection from the map.
		delete(t.clients, ipAddress)
	}
	t.clientsMu.Unlock()

	// Close the connection for safe keeping.
	conn.Close()
}
