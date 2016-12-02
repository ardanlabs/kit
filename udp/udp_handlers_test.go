package udp_test

import (
	"io"
	"net"
	"sync/atomic"
	"time"

	"github.com/ardanlabs/kit/udp"
)

// udpConnHandler is required to process data.
type udpConnHandler struct{}

// Bind is called to init to reader and writer.
func (udpConnHandler) Bind(ctx interface{}, listener *net.UDPConn) (io.Reader, io.Writer) {
	return listener, listener
}

//==============================================================================

// udpReqHandler is required to process client messages.
type udpReqHandler struct{}

// Read implements the udp.ReqHandler interface. It is provided a request
// value to popular and a io.Reader that was created in the Bind above.
func (udpReqHandler) Read(ctx interface{}, reader io.Reader) (*net.UDPAddr, []byte, int, error) {
	listener := reader.(*net.UDPConn)

	// Each package is 20 bytes in lengrh.
	data := make([]byte, 20)
	length, udpAddr, err := listener.ReadFromUDP(data)
	if err != nil {
		return nil, nil, 0, err
	}

	return udpAddr, data, length, nil
}

var dur int64

// Process is used to handle the processing of the message. This method
// is called on a routine from a pool of routines.
func (udpReqHandler) Process(ctx interface{}, r *udp.Request) {
	if r.Length != 20 {
		return
	}

	// Extract the header from the first 8 bytes.
	// h := struct {
	// 	Raw           []byte
	// 	Length        int
	// 	Version       uint8
	// 	TransactionID uint8
	// 	OpCode        uint8
	// 	StatusCode    uint8
	// 	StreamHandle  uint32
	// }{
	// 	Raw:           r.Data,
	// 	Length:        r.Length,
	// 	Version:       uint8(r.Data[0]),
	// 	TransactionID: uint8(r.Data[1]),
	// 	OpCode:        uint8(r.Data[2]),
	// 	StatusCode:    uint8(r.Data[3]),
	// 	StreamHandle:  uint32(binary.BigEndian.Uint32(r.Data[4:8])),
	// }

	resp := udp.Response{
		UDPAddr: r.UDPAddr,
		Data:    []byte("GOT IT"),
		Length:  6,
		Complete: func(rsp *udp.Response) {
			d := int64(time.Now().Sub(r.ReadAt))
			atomic.StoreInt64(&dur, d)
		},
	}

	r.UDP.Do(ctx, &resp)
}

//==============================================================================

type udpRespHandler struct{}

// Write is provided the user-defined writer and the data to write.
func (udpRespHandler) Write(ctx interface{}, r *udp.Response, writer io.Writer) {
	listener := writer.(*net.UDPConn)
	listener.WriteToUDP(r.Data, r.UDPAddr)
}
