package main

import (
	"bufio"
	"io"
	"net"

	"github.com/ardanlabs/kit/log"
	"github.com/ardanlabs/kit/tcp"
)

// tcpConnHandler is required to process data.
type tcpConnHandler struct{}

// Bind is called to init a reader and writer.
func (tcpConnHandler) Bind(traceID string, conn net.Conn) (io.Reader, io.Writer) {
	return bufio.NewReader(conn), bufio.NewWriter(conn)
}

//==============================================================================

// tcpReqHandler is required to process client messages.
type tcpReqHandler struct{}

// Read implements the tcp.ReqHandler interface. It is provided a request
// value to populate and a io.Reader that was created in the Bind above.
func (tcpReqHandler) Read(traceID string, ipAddress string, reader io.Reader) ([]byte, int, error) {
	log.Dev(traceID, "Read", "Started : Waiting For Data")

	bufReader := reader.(*bufio.Reader)

	// Read a small string to keep the code simple.
	line, err := bufReader.ReadString('\n')
	if err != nil {
		log.Error(traceID, "Read", err, "Completed")
		return nil, 0, err
	}

	log.Dev(traceID, "Read", "Completed : IP[%s] Length[%d] Data[%s]", ipAddress, len(line), line)
	return []byte(line), len(line), nil
}

// Process is used to handle the processing of the message. This method
// is called on a routine from a pool of routines.
func (tcpReqHandler) Process(traceID string, r *tcp.Request) {
	log.Dev(traceID, "Process", "Started : IP[%s] Length[%d] ReadAt[%v]", r.TCPAddr.String(), r.Length, r.ReadAt)

	log.User(traceID, "Process", "Data : %s", string(r.Data))

	resp := tcp.Response{
		TCPAddr: r.TCPAddr,
		Data:    []byte("GOT IT\n"),
		Length:  7,
		Complete: func(rsp *tcp.Response) {
			log.Dev(traceID, "Process", "*****************> %v", rsp)
		},
	}

	r.TCP.Do(traceID, &resp)

	log.Dev(traceID, "Process", "Completed")
}

//==============================================================================

type tcpRespHandler struct{}

// Write is provided the user-defined writer and the data to write.
func (tcpRespHandler) Write(traceID string, r *tcp.Response, writer io.Writer) {
	log.Dev(traceID, "Write", "Started : Length[%d]", len(r.Data))

	bufWriter := writer.(*bufio.Writer)
	bufWriter.WriteString(string(r.Data))
	bufWriter.Flush()

	log.Dev(traceID, "Write", "Completed")
}
