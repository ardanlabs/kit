package udp

import (
	"fmt"
	"io"
	"net"
	"time"
)

// ConnHandler is implemented by the user to bind the listener
// to a reader and writer for processing.
type ConnHandler interface {
	// Bind is called to set the reader and writer.
	Bind(logCtx interface{}, listener *net.UDPConn) (io.Reader, io.Writer)
}

//==============================================================================

// ReqHandler is implemented by the user to implement the processing
// of request messages from the client.
type ReqHandler interface {
	// Read is provided the user-defined reader and must return the data read
	// off the wire and the length. Returning io.EOF or a non temporary error
	// will show down the listener.
	Read(logCtx interface{}, reader io.Reader) (*net.UDPAddr, []byte, int, error)

	// Process is used to handle the processing of the request. This method
	// is called on a routine from a pool of routines.
	Process(logCtx interface{}, r *Request)
}

// Request is the message received by the client.
type Request struct {
	UDP     *UDP
	UDPAddr *net.UDPAddr
	IsIPv6  bool
	ReadAt  time.Time
	Data    []byte
	Length  int
}

// logCtx returns a string to use for the logging logCtx.
func (r *Request) logCtx(logCtx interface{}) string {
	return fmt.Sprintf("%s-%s", logCtx, r.UDPAddr)
}

// Work implements the worker inteface for processing messages. This is called
// from a routine in the work pool.
func (r *Request) Work(logCtx interface{}, id int) {
	r.UDP.ReqHandler.Process(logCtx, r)
}

//==============================================================================

// RespHandler is implemented by the user to implement the processing
// of the response messages to the client.
type RespHandler interface {
	// Write is provided the user-defined writer and the data to write.
	Write(logCtx interface{}, r *Response, writer io.Writer)
}

// Response is message to send to the client.
type Response struct {
	UDPAddr  *net.UDPAddr
	Data     []byte
	Length   int
	Complete func(r *Response)

	udp    *UDP
	logCtx interface{}
}

// Work implements the worker interface for sending messages. Called by
// AsyncSend via the d.client.Do(logCtx, &resp) method call.
func (r *Response) Work(logCtx interface{}, id int) {
	r.udp.RespHandler.Write(logCtx, r, r.udp.writer)
	if r.Complete != nil {
		r.Complete(r)
	}
}
