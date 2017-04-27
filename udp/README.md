
# udp
    import "github.com/ardanlabs/kit/udp"

Package udp provides the boilerpale code for working with UDP based data. The package
allows you to establish a UDP listener that can accept data on a specified IP address
and port. It also provides a function to send data back to the client. The processing
of received data and sending data happens on a configured routine pool, so concurrency
is handled.

There are three interfaces that need to be implemented to use the package. These
interfaces provide the API for processing data.

ConnHandler


	type ConnHandler interface {
	    Bind(logCtx string, listener *net.UDPConn) (io.Reader, io.Writer)
	}

The ConnHandler interface is implemented by the user to bind the listener
to a reader and writer for processing.

ReqHandler


	type ReqHandler interface {
	    Read(logCtx string, reader io.Reader) (*net.UDPAddr, []byte, int, error)
	    Process(logCtx string, r *Request)
	}
	
	type Request struct {
	    UDP     *UDP
	    UDPAddr *net.UDPAddr
	    Data    []byte
	    Length  int
	}

The ReqHandler interface is implemented by the user to implement the processing
of request messages from the client. Read is provided the user-defined reader
and must return the data read off the wire and the length. Returning io.EOF or
a non temporary error will show down the listener.

RespHandler


	type RespHandler interface {
	    Write(logCtx string, r *Response, writer io.Writer)
	}
	
	type Response struct {
	    UDPAddr *net.UDPAddr
	    Data    []byte
	    Length  int
	}

The RespHandler interface is implemented by the user to implement the processing
of the response messages to the client. Write is provided the user-defined
writer and the data to write.

### Sample Application
After implementing the interfaces, the following code is all that is needed to
start processing messages.


	func main() {
	    log.Startf("TEST", "main", "Starting Test App")
	
	    cfg := udp.Config{
	        NetType:      "udp4",
	        Addr:         ":9000",
	        WorkRoutines: 2,
	        WorkStats:    time.Minute,
	        ConnHandler:  udpConnHandler{},
	        ReqHandler:   udpReqHandler{},
	        RespHandler:  udpRespHandler{},
	    }
	
	    u, err := udp.New("TEST", &cfg)
	    if err != nil {
	        log.ErrFatal(err, "TEST", "main")
	    }
	
	    if err := u.Start("TEST"); err != nil {
	        log.ErrFatal(err, "TEST", "main")
	    }
	
	    // Wait for a signal to shutdown.
	    sigChan := make(chan os.Signal, 1)
	    signal.Notify(sigChan, os.Interrupt)
	    <-sigChan
	
	    u.Stop("TEST")
	
	    log.Complete("TEST", "main")
	}





## Variables
``` go
var (
    ErrInvalidConfiguration = errors.New("Invalid Configuration")
    ErrInvalidNetType       = errors.New("Invalid NetType Configuration")
    ErrInvalidConnHandler   = errors.New("Invalid Connection Handler Configuration")
    ErrInvalidReqHandler    = errors.New("Invalid Request Handler Configuration")
    ErrInvalidRespHandler   = errors.New("Invalid Response Handler Configuration")
)
```
Set of error variables for start up.



## type Config
``` go
type Config struct {
    NetType string // "udp", udp4" or "udp6"
    Addr    string // "host:port" or "[ipv6-host%zone]:port"

    ConnHandler ConnHandler // Support for binding new connections to a reader and writer.
    ReqHandler  ReqHandler  // Support for handling the specific request workflow.
    RespHandler RespHandler // Support for handling the specific response workflow.

    OptEvent
}
```
Config provides a data structure of required configuration parameters.











### func (\*Config) Event
``` go
func (cfg *Config) Event(event string, format string, a ...interface{})
```
Event fires events back to the user for important events.



### func (\*Config) Validate
``` go
func (cfg *Config) Validate() error
```
Validate checks the configuration to required items.



## type ConnHandler
``` go
type ConnHandler interface {

    // Bind is called to set the reader and writer.
    Bind(listener *net.UDPConn) (io.Reader, io.Writer)
}
```
ConnHandler is implemented by the user to bind the listener
to a reader and writer for processing.











## type OptEvent
``` go
type OptEvent struct {
    Event func(event string, format string, a ...interface{})
}
```
OptEvent defines an handler used to provide events.











## type ReqHandler
``` go
type ReqHandler interface {

    // Read is provided the user-defined reader and must return the data read
    // off the wire and the length. Returning io.EOF or a non temporary error
    // will show down the listener.
    Read(reader io.Reader) (*net.UDPAddr, []byte, int, error)

    // Process is used to handle the processing of the request. This method
    // is called on a routine from a pool of routines.
    Process(r *Request)
}
```
ReqHandler is implemented by the user to implement the processing
of request messages from the client.











## type Request
``` go
type Request struct {
    UDP     *UDP
    UDPAddr *net.UDPAddr
    IsIPv6  bool
    ReadAt  time.Time
    Data    []byte
    Length  int
}
```
Request is the message received by the client.











## type RespHandler
``` go
type RespHandler interface {

    // Write is provided the user-defined writer and the data to write.
    Write(r *Response, writer io.Writer) error
}
```
RespHandler is implemented by the user to implement the processing
of the response messages to the client.











## type Response
``` go
type Response struct {
    UDPAddr *net.UDPAddr
    Data    []byte
    Length  int
}
```
Response is message to send to the client.











## type UDP
``` go
type UDP struct {
    Config
    Name string
    // contains filtered or unexported fields
}
```
UDP manages message to a specific ip address and port.









### func New
``` go
func New(name string, cfg Config) (*UDP, error)
```
New creates a new manager to service clients.




### func (\*UDP) Addr
``` go
func (d *UDP) Addr() net.Addr
```
Addr returns the local listening network address.



### func (\*UDP) Send
``` go
func (d *UDP) Send(r *Response) error
```
Send will deliver the response back to the client.



### func (\*UDP) Start
``` go
func (d *UDP) Start() error
```
Start begins to accept data.



### func (\*UDP) Stop
``` go
func (d *UDP) Stop() error
```
Stop shuts down the manager and closes all connections.









- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)