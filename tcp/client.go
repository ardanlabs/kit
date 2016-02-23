package tcp

import (
	"bytes"
	"io"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ardanlabs/kit/log"
)

// client represents a single networked connection.
type client struct {
	context   string
	t         *TCP
	conn      net.Conn
	ipAddress string
	isIPv6    bool
	reader    io.Reader
	writer    io.Writer
	wg        sync.WaitGroup
}

// newClient creates a new client for an incoming connection.
func newClient(context string, t *TCP, conn net.Conn) *client {
	ipAddress := conn.RemoteAddr().String()
	log.Dev(context, "newClient", "Started : IPAddress[%s]", ipAddress)

	// Ask the user to bind the reader and writer they want to
	// use for this connection.
	r, w := t.ConnHandler.Bind(context, conn)

	c := client{
		context:   context,
		t:         t,
		conn:      conn,
		ipAddress: ipAddress,
		reader:    r,
		writer:    w,
	}

	// Check to see if this connection is ipv6.
	if raddr := conn.RemoteAddr().(*net.TCPAddr); raddr.IP.To4() == nil {
		c.isIPv6 = true
	}

	// Launch a goroutine for this connection.
	c.wg.Add(1)
	go c.read()

	log.Dev(context, "newClient", "Completed")
	return &c
}

// drop closes the client connection and read operation.
func (c *client) drop() {
	log.Dev(c.context, "drop", "Started")

	// Close the connection.
	c.conn.Close()
	c.wg.Wait()

	log.Dev(c.context, "drop", "Completed")
}

// read waits for a message and sends it to the user for procesing.
func (c *client) read() {
	log.Dev(c.context, "read", "Started : Read Processing")

close:
	for {
		// Wait for a message to arrive.
		data, length, err := c.t.ReqHandler.Read(c.context, c.ipAddress, c.reader)
		timeRead := time.Now()

		if err != nil {
			if atomic.LoadInt32(&c.t.shuttingDown) == 0 {
				log.Error(c.context, "read", err, "Completed")
			}

			// temporary is declared to test for the existence of
			// the method coming from the net package.
			type temporary interface {
				Temporary() bool
			}

			if e, ok := err.(temporary); ok {
				if !e.Temporary() {
					break close
				}
			}

			if err == io.EOF {
				break close
			}

			continue
		}

		// Convert the IP:socket for populating TCPAddr value.
		parts := bytes.Split([]byte(c.ipAddress), []byte(":"))
		ipAddress := string(parts[0])
		port, _ := strconv.Atoi(string(parts[1]))

		// Create the request.
		r := Request{
			TCP: c.t,
			TCPAddr: &net.TCPAddr{
				IP:   net.ParseIP(ipAddress),
				Port: port,
				Zone: c.t.tcpAddr.Zone,
			},
			IsIPv6: c.isIPv6,
			ReadAt: timeRead,
			Data:   data,
			Length: length,
		}

		// Send this to the user work pool for processing.
		c.t.recv.Do(c.context, &r)
	}

	log.Dev(c.context, "read", "Shutting Down Client Routine")

	// Remove from the list of connections.
	c.t.remove(c.context, c.conn)

	c.wg.Done()
	log.Dev(c.context, "read", "Completed")
	return
}
