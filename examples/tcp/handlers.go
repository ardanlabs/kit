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
func (tcpConnHandler) Bind(conn net.Conn) (io.Reader, io.Writer) {
	return bufio.NewReader(conn), bufio.NewWriter(conn)
}

// tcpReqHandler is required to process client messages.
type tcpReqHandler struct{}

// Read implements the tcp.ReqHandler interface. It is provided a request
// value to populate and a io.Reader that was created in the Bind above.
func (tcpReqHandler) Read(ipAddress string, reader io.Reader) ([]byte, int, error) {
	log.Dev("handler", "Read", "Started : Waiting For Data")

	bufReader := reader.(*bufio.Reader)

	// Read a small string to keep the code simple.
	line, err := bufReader.ReadString('\n')
	if err != nil {
		log.Error("handler", "Read", err, "Completed")
		return nil, 0, err
	}

	log.Dev("handler", "Read", "Completed : IP[%s] Length[%d] Data[%s]", ipAddress, len(line), line)
	return []byte(line), len(line), nil
}

// Process is used to handle the processing of the message. This method
// is called on a routine from a pool of routines.
func (tcpReqHandler) Process(r *tcp.Request) {
	log.Dev("handler", "Process", "Started : IP[%s] Length[%d] ReadAt[%v]", r.TCPAddr.String(), r.Length, r.ReadAt)
	log.User("handler", "Process", "Data : %s", string(r.Data))

	resp := tcp.Response{
		TCPAddr: r.TCPAddr,
		Data:    []byte("GOT IT\n"),
		Length:  7,
	}

	r.TCP.Send(&resp)

	log.Dev("handler", "Process", "Completed")
}

type tcpRespHandler struct{}

// Write is provided the user-defined writer and the data to write.
func (tcpRespHandler) Write(r *tcp.Response, writer io.Writer) error {
	log.Dev("handler", "Write", "Started : Length[%d]", len(r.Data))
	defer log.Dev("handler", "Write", "Completed")

	bufWriter := writer.(*bufio.Writer)
	if _, err := bufWriter.WriteString(string(r.Data)); err != nil {
		return err
	}
	return bufWriter.Flush()
}
