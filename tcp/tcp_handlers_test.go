package tcp_test

import (
	"bufio"
	"io"
	"net"
	"sync/atomic"
	"time"

	"github.com/ardanlabs/kit/tcp"
)

// tcpConnHandler is required to process data.
type tcpConnHandler struct{}

// Bind is called to init to reader and writer.
func (tch tcpConnHandler) Bind(ctx interface{}, conn net.Conn) (io.Reader, io.Writer) {
	return bufio.NewReader(conn), bufio.NewWriter(conn)
}

//==============================================================================

// tcpReqHandler is required to process client messages.
type tcpReqHandler struct{}

// Read implements the udp.ReqHandler interface. It is provided a request
// value to popular and a io.Reader that was created in the Bind above.
func (tcpReqHandler) Read(ctx interface{}, ipAddress string, reader io.Reader) ([]byte, int, error) {
	bufReader := reader.(*bufio.Reader)

	// Read a small string to keep the code simple.
	line, err := bufReader.ReadString('\n')
	if err != nil {
		return nil, 0, err
	}

	return []byte(line), len(line), nil
}

var dur int64

// Process is used to handle the processing of the message. This method
// is called on a routine from a pool of routines.
func (tcpReqHandler) Process(ctx interface{}, r *tcp.Request) {
	resp := tcp.Response{
		TCPAddr: r.TCPAddr,
		Data:    []byte("GOT IT\n"),
		Length:  7,
		Complete: func(rsp *tcp.Response) {
			d := int64(time.Now().Sub(r.ReadAt))
			atomic.StoreInt64(&dur, d)
		},
	}

	r.TCP.Do(ctx, &resp)
}

//==============================================================================

type tcpRespHandler struct{}

// Write is provided the user-defined writer and the data to write.
func (tcpRespHandler) Write(ctx interface{}, r *tcp.Response, writer io.Writer) {
	bufWriter := writer.(*bufio.Writer)
	bufWriter.WriteString(string(r.Data))
	bufWriter.Flush()
}
