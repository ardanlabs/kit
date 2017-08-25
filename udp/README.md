

# udp
`import "github.com/ardanlabs/kit/udp"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)

## <a name="pkg-overview">Overview</a>
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




## <a name="pkg-index">Index</a>
* [Variables](#pkg-variables)
* [type Config](#Config)
  * [func (cfg *Config) Event(event string, format string, a ...interface{})](#Config.Event)
  * [func (cfg *Config) Validate() error](#Config.Validate)
* [type ConnHandler](#ConnHandler)
* [type OptEvent](#OptEvent)
* [type ReqHandler](#ReqHandler)
* [type Request](#Request)
* [type RespHandler](#RespHandler)
* [type Response](#Response)
* [type UDP](#UDP)
  * [func New(name string, cfg Config) (*UDP, error)](#New)
  * [func (d *UDP) Addr() net.Addr](#UDP.Addr)
  * [func (d *UDP) Send(r *Response) error](#UDP.Send)
  * [func (d *UDP) Start() error](#UDP.Start)
  * [func (d *UDP) Stop() error](#UDP.Stop)


#### <a name="pkg-files">Package files</a>
[doc.go](/src/github.com/ardanlabs/kit/udp/doc.go) [handlers.go](/src/github.com/ardanlabs/kit/udp/handlers.go) [udp.go](/src/github.com/ardanlabs/kit/udp/udp.go) [udp_config.go](/src/github.com/ardanlabs/kit/udp/udp_config.go) 



## <a name="pkg-variables">Variables</a>
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




## <a name="Config">type</a> [Config](/src/target/udp_config.go?s=227:845#L1)
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










### <a name="Config.Event">func</a> (\*Config) [Event](/src/target/udp_config.go?s=1369:1440#L40)
``` go
func (cfg *Config) Event(event string, format string, a ...interface{})
```
Event fires events back to the user for important events.




### <a name="Config.Validate">func</a> (\*Config) [Validate](/src/target/udp_config.go?s=903:938#L15)
``` go
func (cfg *Config) Validate() error
```
Validate checks the configuration to required items.




## <a name="ConnHandler">type</a> [ConnHandler](/src/target/handlers.go?s=447:579#L18)
``` go
type ConnHandler interface {

    // Bind is called to set the reader and writer.
    Bind(listener *net.UDPConn) (io.Reader, io.Writer)
}
```
ConnHandler is implemented by the user to bind the listener
to a reader and writer for processing.










## <a name="OptEvent">type</a> [OptEvent](/src/target/udp_config.go?s=68:151#L1)
``` go
type OptEvent struct {
    Event func(event string, format string, a ...interface{})
}
```
OptEvent defines an handler used to provide events.










## <a name="ReqHandler">type</a> [ReqHandler](/src/target/handlers.go?s=690:1109#L26)
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










## <a name="Request">type</a> [Request](/src/target/handlers.go?s=96:217#L1)
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










## <a name="RespHandler">type</a> [RespHandler](/src/target/handlers.go?s=1224:1368#L40)
``` go
type RespHandler interface {

    // Write is provided the user-defined writer and the data to write.
    Write(r *Response, writer io.Writer) error
}
```
RespHandler is implemented by the user to implement the processing
of the response messages to the client.










## <a name="Response">type</a> [Response](/src/target/handlers.go?s=265:340#L10)
``` go
type Response struct {
    UDPAddr *net.UDPAddr
    Data    []byte
    Length  int
}
```
Response is message to send to the client.










## <a name="UDP">type</a> [UDP](/src/target/udp.go?s=717:953#L19)
``` go
type UDP struct {
    Config
    Name string
    // contains filtered or unexported fields
}
```
UDP manages message to a specific ip address and port.







### <a name="New">func</a> [New](/src/target/udp.go?s=1004:1051#L38)
``` go
func New(name string, cfg Config) (*UDP, error)
```
New creates a new manager to service clients.





### <a name="UDP.Addr">func</a> (\*UDP) [Addr](/src/target/udp.go?s=4944:4973#L216)
``` go
func (d *UDP) Addr() net.Addr
```
Addr returns the local listening network address.




### <a name="UDP.Send">func</a> (\*UDP) [Send](/src/target/udp.go?s=4807:4844#L211)
``` go
func (d *UDP) Send(r *Response) error
```
Send will deliver the response back to the client.




### <a name="UDP.Start">func</a> (\*UDP) [Start](/src/target/udp.go?s=1675:1702#L70)
``` go
func (d *UDP) Start() error
```
Start begins to accept data.




### <a name="UDP.Stop">func</a> (\*UDP) [Stop](/src/target/udp.go?s=4236:4262#L183)
``` go
func (d *UDP) Stop() error
```
Stop shuts down the manager and closes all connections.








- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
