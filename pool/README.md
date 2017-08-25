

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
	    Work(ctx context.Context, id int)
	}

The Worker interface is how you can provide work to the pool. A user-defined type
implements this interface, then values of that type can be passed into the Do
function.

### Sample Application
<a href="https://github.com/ardanlabs/kit/blob/master/examples/pool/main.go">https://github.com/ardanlabs/kit/blob/master/examples/pool/main.go</a>




## <a name="pkg-index">Index</a>
* [Variables](#pkg-variables)
* [type Config](#Config)
  * [func (cfg *Config) Event(ctx context.Context, event string, format string, a ...interface{})](#Config.Event)
* [type OptEvent](#OptEvent)
* [type Pool](#Pool)
  * [func New(name string, cfg Config) (*Pool, error)](#New)
  * [func (p *Pool) Do(ctx context.Context, work Worker)](#Pool.Do)
  * [func (p *Pool) DoCancel(ctx context.Context, work Worker) error](#Pool.DoCancel)
  * [func (p *Pool) Shutdown()](#Pool.Shutdown)
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




## <a name="Config">type</a> [Config](/src/target/pool.go?s=1564:2003#L45)
``` go
type Config struct {
    MinRoutines func() int // Initial and minimum number of routines always in the pool.
    MaxRoutines func() int // Maximum number of routines we will ever grow the pool to.

    OptEvent
}
```
Config provides configuration for the pool.










### <a name="Config.Event">func</a> (\*Config) [Event](/src/target/pool.go?s=2066:2158#L57)
``` go
func (cfg *Config) Event(ctx context.Context, event string, format string, a ...interface{})
```
Event fires events back to the user for important events.




## <a name="OptEvent">type</a> [OptEvent](/src/target/pool.go?s=1411:1515#L40)
``` go
type OptEvent struct {
    Event func(ctx context.Context, event string, format string, a ...interface{})
}
```
OptEvent defines an handler used to provide events.










## <a name="Pool">type</a> [Pool](/src/target/pool.go?s=2339:3346#L65)
``` go
type Pool struct {
    Config
    Name string // Name of this pool.
    // contains filtered or unexported fields
}
```
Pool provides a pool of routines that can execute any Worker
tasks that are submitted.







### <a name="New">func</a> [New](/src/target/pool.go?s=3375:3423#L88)
``` go
func New(name string, cfg Config) (*Pool, error)
```
New creates a new Pool.





### <a name="Pool.Do">func</a> (\*Pool) [Do](/src/target/pool.go?s=4248:4299#L131)
``` go
func (p *Pool) Do(ctx context.Context, work Worker)
```
Do waits for the goroutine pool to take the work to be executed.




### <a name="Pool.DoCancel">func</a> (\*Pool) [DoCancel](/src/target/pool.go?s=4630:4693#L147)
``` go
func (p *Pool) DoCancel(ctx context.Context, work Worker) error
```
DoCancel waits for the goroutine pool to take the work to be executed
or gives up if the Context is cancelled. Only use when you want to throw
away work and not push back.




### <a name="Pool.Shutdown">func</a> (\*Pool) [Shutdown](/src/target/pool.go?s=3984:4009#L120)
``` go
func (p *Pool) Shutdown()
```
Shutdown waits for all the workers to finish.




### <a name="Pool.Stats">func</a> (\*Pool) [Stats](/src/target/pool.go?s=5042:5069#L169)
``` go
func (p *Pool) Stats() Stat
```
Stats returns the current snapshot of the pool stats.




## <a name="Stat">type</a> [Stat](/src/target/pool.go?s=1016:1354#L31)
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










## <a name="Worker">type</a> [Worker](/src/target/pool.go?s=796:856#L20)
``` go
type Worker interface {
    Work(ctx context.Context, id int)
}
```
Worker must be implemented by types that want to use
this worker processes.














- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
