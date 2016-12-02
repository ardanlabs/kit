// Package tcp provides the boilerpale code for working with TCP based data. The package
// allows you to establish a TCP listener that can accept client connections on a specified IP address
// and port. It also provides a function to send data back to the client. The processing
// of received data and sending data happens on a configured routine pool, so concurrency
// is handled.
//
// There are three interfaces that need to be implemented to use the package. These
// interfaces provide the API for processing data.
//
// ConnHandler
//
//     type ConnHandler interface {
//         Bind(ctx string, conn net.Conn) (io.Reader, io.Writer)
//     }
//
// The ConnHandler interface is implemented by the user to bind the client connection
// to a reader and writer for processing.
//
// ReqHandler
//
//     type ReqHandler interface {
//         Read(ctx string, ipAddress string, reader io.Reader) ([]byte, int, error)
//         Process(ctx string, r *Request)
//     }
//
//     type Request struct {
//         TCP       *TCP
//         TCPAddr   *net.TCPAddr
//         Data      []byte
//         Length    int
//     }
//
// The ReqHandler interface is implemented by the user to implement the processing
// of request messages from the client. Read is provided an ipaddress and the user-defined
// reader and must return the data read off the wire and the length. Returning io.EOF or a non
// temporary error will show down the listener.
//
// RespHandler
//
//     type RespHandler interface {
//         Write(ctx string, r *Response, writer io.Writer)
//     }
//
//     type Response struct {
//         TCPAddr   *net.TCPAddr
//         Data      []byte
//         Length    int
//     }
//
// The RespHandler interface is implemented by the user to implement the processing
// of the response messages to the client. Write is provided the user-defined
// writer and the data to write.
//
// Sample Application
//
// After implementing the interfaces, the following code is all that is needed to
// start processing messages.
//
//     func main() {
//         log.Startf("TEST", "main", "Starting Test App")
//
//         cfg := tcp.Config{
//             NetType:      "tcp4",
//             Addr:         ":9000",
//             WorkRoutines: 2,
//             WorkStats:    time.Minute,
//             ConnHandler:  tcpConnHandler{},
//             ReqHandler:   udpReqHandler{},
//             RespHandler:  udpRespHandler{},
//         }
//
//         t, err := tcp.New("TEST", &cfg)
//         if err != nil {
//             log.ErrFatal(err, "TEST", "main")
//         }
//
//         if err := t.Start("TEST"); err != nil {
//             log.ErrFatal(err, "TEST", "main")
//         }
//
//         // Wait for a signal to shutdown.
//         sigChan := make(chan os.Signal, 1)
//         signal.Notify(sigChan, os.Interrupt)
//         <-sigChan
//
//         t.Stop("TEST")
//
//         log.Complete("TEST", "main")
//     }
package tcp
