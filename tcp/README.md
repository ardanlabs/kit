
# tcp
    import "github.com/ardanlabs/kit/tcp"

Package tcp provides the boilerpale code for working with TCP based data. The package
allows you to establish a TCP listener that can accept client connections on a specified IP address
and port. It also provides a function to send data back to the client.

There are three interfaces that need to be implemented to use the package. These
interfaces provide the API for processing data.

ConnHandler


	type ConnHandler interface {
	    Bind(conn net.Conn) (io.Reader, io.Writer)
	}

The ConnHandler interface is implemented by the user to bind the client connection
to a reader and writer for processing.

ReqHandler


	type ReqHandler interface {
	    Read(ipAddress string, reader io.Reader) ([]byte, int, error)
	    Process(r *Request)
	}
	
	type Request struct {
	    TCP       *TCP
	    TCPAddr   *net.TCPAddr
	    Data      []byte
	    Length    int
	}

The ReqHandler interface is implemented by the user to implement the processing
of request messages from the client. Read is provided an ipaddress and the user-defined
reader and must return the data read off the wire and the length. Returning io.EOF or a non
temporary error will show down the listener.

RespHandler


	type RespHandler interface {
	    Write(r *Response, writer io.Writer) error
	}
	
	type Response struct {
	    TCPAddr   *net.TCPAddr
	    Data      []byte
	    Length    int
	}

The RespHandler interface is implemented by the user to implement the processing
of the response messages to the client. Write is provided the user-defined
writer and the data to write.

### Sample Application
After implementing the interfaces, the following code is all that is needed to
start processing messages.


	func main() {
	    log.Println("Starting Test App")
	
	    cfg := tcp.Config{
	        NetType:      "tcp4",
	        Addr:         ":9000",
	        WorkRoutines: 2,
	        WorkStats:    time.Minute,
	        ConnHandler:  tcpConnHandler{},
	        ReqHandler:   udpReqHandler{},
	        RespHandler:  udpRespHandler{},
	    }
	
	    t, err := tcp.New(&cfg)
	    if err != nil {
	        log.Println(err)
	         return
	    }
	
	    if err := t.Start(); err != nil {
	        log.Println(err)
	         return
	    }
	
	    // Wait for a signal to shutdown.
	    sigChan := make(chan os.Signal, 1)
	    signal.Notify(sigChan, os.Interrupt)
	    <-sigChan
	
	    t.Stop()
	    log.Println("down")
	}





## Variables
``` go
var (
    ErrInvalidConfiguration = errors.New("invalid configuration")
    ErrInvalidNetType       = errors.New("invalid net type configuration")
    ErrInvalidConnHandler   = errors.New("invalid connection handler configuration")
    ErrInvalidReqHandler    = errors.New("invalid request handler configuration")
    ErrInvalidRespHandler   = errors.New("invalid response handler configuration")
)
```
Set of error variables for start up.



## type Config
``` go
type Config struct {
    NetType string // "tcp", tcp4" or "tcp6"
    Addr    string // "host:port" or "[ipv6-host%zone]:port"

    ConnHandler ConnHandler // Support for binding new connections to a reader and writer.
    ReqHandler  ReqHandler  // Support for handling the specific request workflow.
    RespHandler RespHandler // Support for handling the specific response workflow.

    OptRateLimit
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
    Bind(conn net.Conn) (io.Reader, io.Writer)
}
```
ConnHandler is implemented by the user to bind the connection
to a reader and writer for processing.











## type OptEvent
``` go
type OptEvent struct {
    Event func(event string, format string, a ...interface{})
}
```
OptEvent defines an handler used to provide events.











## type OptRateLimit
``` go
type OptRateLimit struct {
    RateLimit func() time.Duration // Connection rate limit per single connection.
}
```
OptRateLimit declares fields for the user to provide configuration
for connection rate limit.











## type ReqHandler
``` go
type ReqHandler interface {

    // Read is provided an ipaddress and the user-defined reader and must return
    // the data read off the wire and the length. Returning io.EOF or a non
    // temporary error will show down the listener.
    Read(ipAddress string, reader io.Reader) ([]byte, int, error)

    // Process is used to handle the processing of the request.
    Process(r *Request)
}
```
ReqHandler is implemented by the user to implement the processing
of request messages from the client.











## type Request
``` go
type Request struct {
    TCP     *TCP
    TCPAddr *net.TCPAddr
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

    // Write is provided the response to write and the user-defined writer.
    Write(r *Response, writer io.Writer) error
}
```
RespHandler is implemented by the user to implement the processing
of the response messages to the client.











## type Response
``` go
type Response struct {
    TCPAddr *net.TCPAddr
    Data    []byte
    Length  int
}
```
Response is message to send to the client.











## type TCP
``` go
type TCP struct {
    Config
    Name string
    // contains filtered or unexported fields
}
```
TCP contains a set of networked client connections.









### func New
``` go
func New(name string, cfg Config) (*TCP, error)
```
New creates a new manager to service clients.




### func (\*TCP) Addr
``` go
func (t *TCP) Addr() net.Addr
```
Addr returns the listener's network address. This may be different than the values
provided in the configuration, for example if configuration port value is 0.



### func (\*TCP) DropConnections
``` go
func (t *TCP) DropConnections(drop bool)
```
DropConnections sets a flag to tell the accept routine to immediately
drop connections that come in.



### func (\*TCP) Send
``` go
func (t *TCP) Send(r *Response) error
```
Send will deliver the response back to the client.



### func (\*TCP) Start
``` go
func (t *TCP) Start() error
```
Start creates the accept routine and begins to accept connections.



### func (\*TCP) Stop
``` go
func (t *TCP) Stop() error
```
Stop shuts down the manager and closes all connections.









- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)