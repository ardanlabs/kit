

# pool
`import "github.com/ardanlabs/kit/pool"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)
* [Examples](#pkg-examples)

## <a name="pkg-overview">Overview</a>
Package pool manages a pool of routines to perform work. It does so my providing
a Do function that will block when the pool is busy. This also allows the pool
to monitor and report pushback. The pool also supports the dynamic re-sizing
of the number of routines in the pool.

Worker


	type Worker interface {
	    Work(context interface{}, id int)
	}

The Worker interface is how you can provide work to the pool. A user-defined type
implements this interface, then values of that type can be passed into the Do
function.

### Sample Application
The following is a sample application using the work pool.


	// theWork is the customer work type for using the pool.
	type theWork struct{}
	
	// Work implements the DoWorker interface.
	func (*theWork) Work(context string, id int) {
	    fmt.Printf("%s : Performing Work\n", context)
	}
	
	// ExampleNewDoPool provides a basic example for using a DoPool.
	func ExampleNewDoPool() {
	    // Create a new do pool.
	    p, err := pool.New(context, "TheWork", 3, func() time.Duration { return time.Minute })
	    if err != nil {
	        fmt.Println(err)
	        return
	    }
	
	    // Pass in some work to be performed.
	    p.Do("TEST", &theWork{})
	    p.Do("TEST", &theWork{})
	    p.Do("TEST", &theWork{})
	
	    // Wait to the work to be processed.
	    time.Sleep(1 * time.Second)
	
	    // Shutdown the pool.
	    p.Shutdown(context)
	}




## <a name="pkg-index">Index</a>
* [Variables](#pkg-variables)
* [type Config](#Config)
  * [func (cfg *Config) Event(context interface{}, event string, format string, a ...interface{})](#Config.Event)
* [type OptEvent](#OptEvent)
* [type Pool](#Pool)
  * [func New(context interface{}, name string, cfg Config) (*Pool, error)](#New)
  * [func (p *Pool) Do(context interface{}, work Worker)](#Pool.Do)
  * [func (p *Pool) DoWait(context interface{}, work Worker, duration &lt;-chan time.Time) error](#Pool.DoWait)
  * [func (p *Pool) Shutdown(context interface{})](#Pool.Shutdown)
  * [func (p *Pool) Stats() Stat](#Pool.Stats)
* [type Stat](#Stat)
* [type Worker](#Worker)

#### <a name="pkg-examples">Examples</a>
* [New](#example_New)

#### <a name="pkg-files">Package files</a>
[doc.go](/src/github.com/ardanlabs/kit/pool/doc.go) [pool.go](/src/github.com/ardanlabs/kit/pool/pool.go) 



## <a name="pkg-variables">Variables</a>
``` go
var (
    ErrNilMinRoutines        = errors.New("Invalid (nil) minimum number of routines")
    ErrNilMaxRoutines        = errors.New("Invalid (nil) maximum number of routines")
    ErrInvalidMinRoutines    = errors.New("Invalid minimum number of routines")
    ErrInvalidMaxRoutines    = errors.New("Invalid maximum number of routines")
    ErrInvalidAdd            = errors.New("Invalid number of routines to add")
    ErrInvalidMetricHandler  = errors.New("Invalid metric handler")
    ErrInvalidMetricInterval = errors.New("Invalid metric interval")
)
```
Set of error variables for start up.




## <a name="Config">type</a> [Config](/src/target/pool.go?s=1715:2154#L48)
``` go
type Config struct {
    MinRoutines func() int // Initial and minimum number of routines always in the pool.
    MaxRoutines func() int // Maximum number of routines we will ever grow the pool to.

    OptEvent
}
```
Config provides configuration for the pool.










### <a name="Config.Event">func</a> (\*Config) [Event](/src/target/pool.go?s=2217:2309#L60)
``` go
func (cfg *Config) Event(context interface{}, event string, format string, a ...interface{})
```
Event fires events back to the user for important events.




## <a name="OptEvent">type</a> [OptEvent](/src/target/pool.go?s=1562:1666#L43)
``` go
type OptEvent struct {
    Event func(context interface{}, event string, format string, a ...interface{})
}
```
OptEvent defines an handler used to provide events.










## <a name="Pool">type</a> [Pool](/src/target/pool.go?s=2576:3583#L70)
``` go
type Pool struct {
    Config
    Name string // Name of this pool.
    // contains filtered or unexported fields
}
```
Pool provides a pool of routines that can execute any Worker
tasks that are submitted.







### <a name="New">func</a> [New](/src/target/pool.go?s=3612:3681#L93)
``` go
func New(context interface{}, name string, cfg Config) (*Pool, error)
```
New creates a new Pool.





### <a name="Pool.Do">func</a> (\*Pool) [Do](/src/target/pool.go?s=4541:4592#L136)
``` go
func (p *Pool) Do(context interface{}, work Worker)
```
Do waits for the goroutine pool to take the work to be executed.




### <a name="Pool.DoWait">func</a> (\*Pool) [DoWait](/src/target/pool.go?s=4933:5021#L152)
``` go
func (p *Pool) DoWait(context interface{}, work Worker, duration <-chan time.Time) error
```
DoWait waits for the goroutine pool to take the work to be executed or gives
up after the allotted duration. Only use when you want to throw away work and
not push back.




### <a name="Pool.Shutdown">func</a> (\*Pool) [Shutdown](/src/target/pool.go?s=4258:4302#L125)
``` go
func (p *Pool) Shutdown(context interface{})
```
Shutdown waits for all the workers to finish.




### <a name="Pool.Stats">func</a> (\*Pool) [Stats](/src/target/pool.go?s=5380:5407#L174)
``` go
func (p *Pool) Stats() Stat
```
Stats returns the current snapshot of the pool stats.




## <a name="Stat">type</a> [Stat](/src/target/pool.go?s=1085:1423#L32)
``` go
type Stat struct {
    Routines    int64 // Current number of routines.
    Pending     int64 // Pending number of routines waiting to submit work.
    Active      int64 // Active number of routines in the work pool.
    Executed    int64 // Number of pieces of work executed.
    MaxRoutines int64 // High water mark of routines the pool has been at.
}
```
Stat contains information about the pool.










## <a name="Worker">type</a> [Worker](/src/target/pool.go?s=861:921#L21)
``` go
type Worker interface {
    Work(context interface{}, id int)
}
```
Worker must be implemented by types that want to use
this worker processes.














- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
