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
	Bind(traceID string, listener *net.UDPConn) (io.Reader, io.Writer)
}

//==============================================================================

// ReqHandler is implemented by the user to implement the processing
// of request messages from the client.
type ReqHandler interface {
	// Read is provided the user-defined reader and must return the data read
	// off the wire and the length. Returning io.EOF or a non temporary error
	// will show down the listener.
	Read(traceID string, reader io.Reader) (*net.UDPAddr, []byte, int, error)

	// Process is used to handle the processing of the request. This method
	// is called on a routine from a pool of routines.
	Process(traceID string, r *Request)
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

// traceID returns a string to use for the logging traceID.
func (r *Request) traceID(traceID string) string {
	return fmt.Sprintf("%s-%s", traceID, r.UDPAddr)
}

// Work implements the worker inteface for processing messages. This is called
// from a routine in the work pool.
func (r *Request) Work(traceID string, id int) {
	r.UDP.ReqHandler.Process(traceID, r)
}

//==============================================================================

// RespHandler is implemented by the user to implement the processing
// of the response messages to the client.
type RespHandler interface {
	// Write is provided the user-defined writer and the data to write.
	Write(traceID string, r *Response, writer io.Writer)
}

// Response is message to send to the client.
type Response struct {
	UDPAddr  *net.UDPAddr
	Data     []byte
	Length   int
	Complete func(r *Response)

	udp     *UDP
	traceID string
}

// Work implements the worker interface for sending messages. Called by
// AsyncSend via the d.client.Do(traceID, &resp) method call.
func (r *Response) Work(traceID string, id int) {
	r.udp.RespHandler.Write(traceID, r, r.udp.writer)
	if r.Complete != nil {
		r.Complete(r)
	}
}
