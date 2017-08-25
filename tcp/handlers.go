package tcp

import (
	"context"
	"io"
	"net"
	"time"
)

// Request is the message received by the client.
type Request struct {
	TCP     *TCP
	TCPAddr *net.TCPAddr
	IsIPv6  bool
	ReadAt  time.Time
	Context context.Context
	Data    []byte
	Length  int
}

// Response is message to send to the client.
type Response struct {
	TCPAddr *net.TCPAddr
	Data    []byte
	Length  int
}

// ConnHandler is implemented by the user to bind the connection
// to a reader and writer for processing.
type ConnHandler interface {

	// Bind is called to set the reader and writer.
	Bind(conn net.Conn) (io.Reader, io.Writer)
}

// ReqHandler is implemented by the user to implement the processing
// of request messages from the client.
type ReqHandler interface {

	// Read is provided a request and a user-defined reader for each client
	// connection on its own routine. Read must read a full request and return
	// the populated request value.
	// Returning io.EOF or a non temporary error will show down the connection.

	// Read is provided an ipaddress and the user-defined reader and must return
	// the data read off the wire and the length. Returning io.EOF or a non
	// temporary error will show down the listener.
	Read(ipAddress string, reader io.Reader) ([]byte, int, error)

	// Process is used to handle the processing of the request.
	Process(r *Request)
}

// RespHandler is implemented by the user to implement the processing
// of the response messages to the client.
type RespHandler interface {

	// Write is provided the response to write and the user-defined writer.
	Write(r *Response, writer io.Writer) error
}
