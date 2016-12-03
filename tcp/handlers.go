package tcp

import (
	"io"
	"net"
	"time"
)

// ConnHandler is implemented by the user to bind the connection
// to a reader and writer for processing.
type ConnHandler interface {

	// Bind is called to set the reader and writer.
	Bind(logCtx interface{}, conn net.Conn) (io.Reader, io.Writer)
}

//==============================================================================

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
	Read(logCtx interface{}, ipAddress string, reader io.Reader) ([]byte, int, error)

	// Process is used to handle the processing of the request. This method
	// is called on a routine from a pool of routines.
	Process(logCtx interface{}, r *Request)
}

// Request is the message received by the client.
type Request struct {
	TCP     *TCP
	TCPAddr *net.TCPAddr
	IsIPv6  bool
	ReadAt  time.Time
	Data    []byte
	Length  int
}

// Work implements the worker interface for processing received messages.
// This is called from a routine in the work pool.
func (r *Request) Work(logCtx interface{}, id int) {
	r.TCP.ReqHandler.Process(logCtx, r)
}

//==============================================================================

// RespHandler is implemented by the user to implement the processing
// of the response messages to the client.
type RespHandler interface {

	// Write is provided the response to write and the user-defined writer.
	Write(logCtx interface{}, r *Response, writer io.Writer)
}

// Response is message to send to the client.
type Response struct {
	TCPAddr  *net.TCPAddr
	Data     []byte
	Length   int
	Complete func(r *Response)

	tcp    *TCP
	client *client
	logCtx interface{}
}

// Work implements the worker interface for sending messages to the client.
// This is called from a routine in the work pool.
func (r *Response) Work(logCtx interface{}, id int) {
	r.tcp.RespHandler.Write(logCtx, r, r.client.writer)
	if r.Complete != nil {
		r.Complete(r)
	}
}
