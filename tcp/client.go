package tcp

import (
	"bytes"
	"context"
	"io"
	"net"
	"strconv"
	"sync"
	"time"
)

// client represents a single networked connection.
type client struct {
	t         *TCP
	conn      net.Conn
	ipAddress string
	isIPv6    bool
	reader    io.Reader
	writer    io.Writer
	wg        sync.WaitGroup

	timeConn time.Time
	lastAct  time.Time
	nReads   int
	nWrites  int
}

// newClient creates a new client for an incoming connection.
func newClient(t *TCP, conn net.Conn) *client {
	now := time.Now().UTC()
	ipAddress := conn.RemoteAddr().String()

	// Ask the user to bind the reader and writer they want to
	// use for this connection.
	r, w := t.ConnHandler.Bind(conn)

	c := client{
		t:         t,
		conn:      conn,
		ipAddress: ipAddress,
		reader:    r,
		writer:    w,
		timeConn:  now,
		lastAct:   now,
	}

	// Check to see if this connection is ipv6.
	if raddr := conn.RemoteAddr().(*net.TCPAddr); raddr.IP.To4() == nil {
		c.isIPv6 = true
	}

	// Launch a goroutine for this connection.
	c.wg.Add(1)
	go c.read()

	return &c
}

// drop closes the client connection and read operation.
func (c *client) drop() {

	// Close the connection.
	c.conn.Close()
	c.wg.Wait()

	c.t.Event(EvtDrop, TypInfo, c.ipAddress, "connect dropped")
}

// read waits for a message and sends it to the user for procesing.
func (c *client) read() {
	c.t.Event(EvtRead, TypTrigger, c.ipAddress, "ready")

close:
	for {

		// Wait for a message to arrive.
		data, length, err := c.t.ReqHandler.Read(c.ipAddress, c.reader)
		c.lastAct = time.Now().UTC()
		c.nReads++

		if err != nil {

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
			IsIPv6:  c.isIPv6,
			ReadAt:  c.lastAct,
			Context: context.Background(),
			Data:    data,
			Length:  length,
		}

		// Process the request on this goroutine that is
		// handling the socket connection.
		c.t.ReqHandler.Process(&r)
	}

	// Remove from the list of connections and report we are done.
	c.t.remove(c.conn)
	c.wg.Done()
	c.t.Event(EvtDrop, TypTrigger, c.ipAddress, "dropped connection")
}
