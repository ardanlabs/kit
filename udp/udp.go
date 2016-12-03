package udp

import (
	"errors"
	"io"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ardanlabs/kit/pool"
)

// Set of error variables for start up.
var (
	ErrInvalidConfiguration = errors.New("Invalid Configuration")
	ErrInvalidNetType       = errors.New("Invalid NetType Configuration")
	ErrInvalidConnHandler   = errors.New("Invalid Connection Handler Configuration")
	ErrInvalidReqHandler    = errors.New("Invalid Request Handler Configuration")
	ErrInvalidRespHandler   = errors.New("Invalid Response Handler Configuration")
)

// temporary is declared to test for the existence of the method coming
// from the net package.
type temporary interface {
	Temporary() bool
}

// UDP manages message to a specific ip address and port.
type UDP struct {
	Config
	Name string

	ipAddress string
	port      int
	udpAddr   *net.UDPAddr

	listener   *net.UDPConn
	listenerMu sync.RWMutex

	reader io.Reader
	writer io.Writer

	recv      *pool.Pool
	send      *pool.Pool
	userPools bool

	wg           sync.WaitGroup
	shuttingDown int32
}

// New creates a new manager to service clients.
func New(logCtx interface{}, name string, cfg Config) (*UDP, error) {
	// Validate the configuration.
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	// Resolve the addr that is provided.
	udpAddr, err := net.ResolveUDPAddr(cfg.NetType, cfg.Addr)
	if err != nil {
		return nil, err
	}

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
		if recv, err = pool.New(logCtx, name+"-Recv", recvCfg); err != nil {
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
		if send, err = pool.New(logCtx, name+"-Send", sendCfg); err != nil {
			return nil, err
		}
	}

	// Are we using user provided work pools. Validation is helping us
	// only have to check one of the two configuration options for this.
	var userPools bool
	if cfg.RecvPool != nil {
		userPools = true
	}

	// Create a UDP for this ipaddress and port.
	udp := UDP{
		Config: cfg,
		Name:   name,

		ipAddress: udpAddr.IP.String(),
		port:      udpAddr.Port,
		udpAddr:   udpAddr,

		recv:      recv,
		send:      send,
		userPools: userPools,
	}

	return &udp, nil
}

// join takes an IP and port values and creates a cleaner string.
func join(ip string, port int) string {
	return net.JoinHostPort(ip, strconv.Itoa(port))
}

// Start begins to accept data.
func (d *UDP) Start(logCtx interface{}) error {
	d.listenerMu.Lock()
	{
		// If the listener has been started already, return an error.
		if d.listener != nil {
			d.listenerMu.Unlock()
			return errors.New("This UDP has already been started")
		}
	}
	d.listenerMu.Unlock()

	d.wg.Add(1)

	// We need to wait for the goroutine to initialize itself.
	var waitStart sync.WaitGroup
	waitStart.Add(1)

	// Start the data accept routine.
	go func() {
		for {
			d.listenerMu.Lock()
			{
				// Start a listener for the specified addr and port is one
				// does not exist.
				if d.listener == nil {
					var err error
					d.listener, err = net.ListenUDP(d.NetType, d.udpAddr)
					if err != nil {
						panic(err)
					}

					// Ask the user to bind the reader and writer they want to
					// use for this listener.
					d.reader, d.writer = d.ConnHandler.Bind(logCtx, d.listener)

					waitStart.Done()

					d.Event(logCtx, "accept", "Waiting For Data : IPAddress[ %s ]", join(d.ipAddress, d.port))
				}
			}
			d.listenerMu.Unlock()

			// Wait for a message to arrive.
			udpAddr, data, length, err := d.ReqHandler.Read(logCtx, d.reader)
			timeRead := time.Now()

			if err != nil {
				if atomic.LoadInt32(&d.shuttingDown) == 1 {
					d.listenerMu.Lock()
					{
						d.listener = nil
					}
					d.listenerMu.Unlock()
					break
				}

				d.Event(logCtx, "accept", "ERROR : %v", err)

				if e, ok := err.(temporary); ok && !e.Temporary() {
					d.listenerMu.Lock()
					{
						d.listener.Close()
						d.listener = nil
					}
					d.listenerMu.Unlock()

					// Don't want to add a flag. So setting this back to
					// 1 so when the listener is re-established, the call
					// to Done does not fail.
					waitStart.Add(1)
				}

				continue
			}

			// Check to see if this message is ipv6.
			isIPv6 := true
			if ip4 := udpAddr.IP.To4(); ip4 != nil {
				// Make sure we return an IPv4 address if udpAddr
				// is an IPv4-mapped IPv6 address.  Otherwise we
				// could end up sending an IPv6 response.
				udpAddr.IP = ip4
				isIPv6 = false
			}

			// Create the request.
			req := Request{
				UDP:     d,
				UDPAddr: udpAddr,
				IsIPv6:  isIPv6,
				ReadAt:  timeRead,
				Data:    data,
				Length:  length,
			}

			// Send this to the user work pool for processing.
			d.recv.Do(req.logCtx(logCtx), &req)
		}

		d.wg.Done()
		d.Event(logCtx, "accept", "Shutdown : IPAddress[ %s ]", join(d.ipAddress, d.port))

		return
	}()

	// Wait for the goroutine to initialize itself.
	waitStart.Wait()

	return nil
}

// Stop shuts down the manager and closes all connections.
func (d *UDP) Stop(logCtx interface{}) error {
	d.listenerMu.Lock()
	{
		// If the listener has been stopped already, return an error.
		if d.listener == nil {
			d.listenerMu.Unlock()
			return errors.New("This UDP has already been stopped")
		}
	}
	d.listenerMu.Unlock()

	// Mark that we are shutting down.
	atomic.StoreInt32(&d.shuttingDown, 1)

	// Don't accept anymore client data.
	d.listenerMu.Lock()
	{
		d.listener.Close()
	}
	d.listenerMu.Unlock()

	// Stop processing all the work.
	if !d.userPools {
		d.recv.Shutdown(logCtx)
		d.send.Shutdown(logCtx)
	}

	// Wait for the accept routine to terminate.
	d.wg.Wait()

	return nil
}

// Do will post the request to be sent by the client worker pool.
func (d *UDP) Do(logCtx interface{}, r *Response) error {
	// Set the unexported fields.
	r.udp = d
	r.logCtx = logCtx

	// Send this to the client work pool for processing.
	d.send.Do(logCtx, r)

	return nil
}

// StatsRecv returns the current snapshot of the recv pool stats.
func (d *UDP) StatsRecv() pool.Stat {
	return d.recv.Stats()
}

// StatsSend returns the current snapshot of the send pool stats.
func (d *UDP) StatsSend() pool.Stat {
	return d.send.Stats()
}

// Addr returns the local listening network address.
func (d *UDP) Addr() net.Addr {
	// We are aware this read is not safe with the
	// goroutine accepting connections.
	if d.listener == nil {
		return nil
	}
	return d.listener.LocalAddr()
}
