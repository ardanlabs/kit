
# tcp
    import "github.com/ardanlabs/kit/tcp"

Package tcp provides the boilerpale code for working with TCP based data. The package
allows you to establish a TCP listener that can accept client connections on a specified IP address
and port. It also provides a function to send data back to the client. The processing
of received data and sending data happens on a configured routine pool, so concurrency
is handled.

There are three interfaces that need to be implemented to use the package. These
interfaces provide the API for processing data.

ConnHandler


	type ConnHandler interface {
	    Bind(context string, conn net.Conn) (io.Reader, io.Writer)
	}

The ConnHandler interface is implemented by the user to bind the client connection
to a reader and writer for processing.

ReqHandler


	type ReqHandler interface {
	    Read(context string, ipAddress string, reader io.Reader) ([]byte, int, error)
	    Process(context string, r *Request)
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
	    Write(context string, r *Response, writer io.Writer)
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
	    log.Startf("TEST", "main", "Starting Test App")
	
	    cfg := tcp.Config{
	        NetType:      "tcp4",
	        Addr:         ":9000",
	        WorkRoutines: 2,
	        WorkStats:    time.Minute,
	        ConnHandler:  tcpConnHandler{},
	        ReqHandler:   udpReqHandler{},
	        RespHandler:  udpRespHandler{},
	    }
	
	    t, err := tcp.New("TEST", &cfg)
	    if err != nil {
	        log.ErrFatal(err, "TEST", "main")
	    }
	
	    if err := t.Start("TEST"); err != nil {
	        log.ErrFatal(err, "TEST", "main")
	    }
	
	    // Wait for a signal to shutdown.
	    sigChan := make(chan os.Signal, 1)
	    signal.Notify(sigChan, os.Interrupt)
	    <-sigChan
	
	    t.Stop("TEST")
	
	    log.Complete("TEST", "main")
	}





## Variables
``` go
var (
    ErrInvalidConfiguration     = errors.New("Invalid Configuration")
    ErrInvalidNetType           = errors.New("Invalid NetType Configuration")
    ErrInvalidConnHandler       = errors.New("Invalid Connection Handler Configuration")
    ErrInvalidReqHandler        = errors.New("Invalid Request Handler Configuration")
    ErrInvalidRespHandler       = errors.New("Invalid Response Handler Configuration")
    ErrInvalidPoolConfiguration = errors.New("Invalid Pool Configuration")
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

    OptUserPool
    OptIntPool

    OptRateLimit
}
```
Config provides a data structure of required configuration parameters.











### func (\*Config) Validate
``` go
func (cfg *Config) Validate() error
```
Validate checks the configuration to required items.



## type ConnHandler
``` go
type ConnHandler interface {
    // Bind is called to set the reader and writer.
    Bind(context interface{}, conn net.Conn) (io.Reader, io.Writer)
}
```
ConnHandler is implemented by the user to bind the connection
to a reader and writer for processing.











## type OptIntPool
``` go
type OptIntPool struct {
    RecvMinPoolSize func() int // Min number of routines the recv pool must have.
    RecvMaxPoolSize func() int // Max number of routines the recv pool can have.
    SendMinPoolSize func() int // Min number of routines the send pool must have.
    SendMaxPoolSize func() int // Max number of routines the send pool can have.
}
```
OptIntPool declares fields for the user to provide configuration
for an internally configured pool.











## type OptRateLimit
``` go
type OptRateLimit struct {
    RateLimit func() time.Duration // Connection rate limit per single connection.
}
```
OptRateLimit declares fields for the user to provide configuration
for connection rate limit.











## type OptUserPool
``` go
type OptUserPool struct {
    RecvPool *pool.Pool // User provided work pool for the receive work.
    SendPool *pool.Pool // User provided work pool for the send work.
}
```
OptUserPool declares fields for the user to pass their own
work pools for configuration.











## type ReqHandler
``` go
type ReqHandler interface {

    // Read is provided an ipaddress and the user-defined reader and must return
    // the data read off the wire and the length. Returning io.EOF or a non
    // temporary error will show down the listener.
    Read(context interface{}, ipAddress string, reader io.Reader) ([]byte, int, error)

    // Process is used to handle the processing of the request. This method
    // is called on a routine from a pool of routines.
    Process(context interface{}, r *Request)
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











### func (\*Request) Work
``` go
func (r *Request) Work(context interface{}, id int)
```
Work implements the worker interface for processing received messages.
This is called from a routine in the work pool.



## type RespHandler
``` go
type RespHandler interface {
    // Write is provided the response to write and the user-defined writer.
    Write(context interface{}, r *Response, writer io.Writer)
}
```
RespHandler is implemented by the user to implement the processing
of the response messages to the client.











## type Response
``` go
type Response struct {
    TCPAddr  *net.TCPAddr
    Data     []byte
    Length   int
    Complete func(r *Response)
    // contains filtered or unexported fields
}
```
Response is message to send to the client.











### func (\*Response) Work
``` go
func (r *Response) Work(context interface{}, id int)
```
Work implements the worker interface for sending messages to the client.
This is called from a routine in the work pool.



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
func New(context interface{}, name string, cfg Config) (*TCP, error)
```
New creates a new manager to service clients.




### func (\*TCP) Addr
``` go
func (t *TCP) Addr() net.Addr
```
Addr returns the listener's network address. This may be different than the values
provided in the configuration, for example if configuration port value is 0.



### func (\*TCP) Do
``` go
func (t *TCP) Do(context interface{}, r *Response) error
```
Do will post the request to be sent by the client worker pool.



### func (\*TCP) DropConnections
``` go
func (t *TCP) DropConnections(context interface{}, drop bool)
```
DropConnections sets a flag to tell the accept routine to immediately
drop connections that come in.



### func (\*TCP) Start
``` go
func (t *TCP) Start(context interface{}) error
```
Start creates the accept routine and begins to accept connections.



### func (\*TCP) StatsRecv
``` go
func (t *TCP) StatsRecv() pool.Stat
```
StatsRecv returns the current snapshot of the recv pool stats.



### func (\*TCP) StatsSend
``` go
func (t *TCP) StatsSend() pool.Stat
```
StatsSend returns the current snapshot of the send pool stats.



### func (\*TCP) Stop
``` go
func (t *TCP) Stop(context interface{}) error
```
Stop shuts down the manager and closes all connections.









- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)