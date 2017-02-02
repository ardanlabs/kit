package udp

import (
	"io"
	"net"
	"time"
)

// Request is the message received by the client.
type Request struct {
	UDP     *UDP
	UDPAddr *net.UDPAddr
	IsIPv6  bool
	ReadAt  time.Time
	Data    []byte
	Length  int
}

// Response is message to send to the client.
type Response struct {
	UDPAddr *net.UDPAddr
	Data    []byte
	Length  int
}

// ConnHandler is implemented by the user to bind the listener
// to a reader and writer for processing.
type ConnHandler interface {

	// Bind is called to set the reader and writer.
	Bind(listener *net.UDPConn) (io.Reader, io.Writer)
}

// ReqHandler is implemented by the user to implement the processing
// of request messages from the client.
type ReqHandler interface {

	// Read is provided the user-defined reader and must return the data read
	// off the wire and the length. Returning io.EOF or a non temporary error
	// will show down the listener.
	Read(reader io.Reader) (*net.UDPAddr, []byte, int, error)

	// Process is used to handle the processing of the request. This method
	// is called on a routine from a pool of routines.
	Process(r *Request)
}

// RespHandler is implemented by the user to implement the processing
// of the response messages to the client.
type RespHandler interface {

	// Write is provided the user-defined writer and the data to write.
	Write(r *Response, writer io.Writer) error
}
