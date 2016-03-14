package pool

import (
	"errors"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

const (
	addRoutine = 1
	rmvRoutine = 2
)

// Set of error variables for start up.
var (
	ErrNilMinRoutines        = errors.New("Invalid (nil) minimum number of routines")
	ErrNilMaxRoutines        = errors.New("Invalid (nil) maximum number of routines")
	ErrInvalidMinRoutines    = errors.New("Invalid minimum number of routines")
	ErrInvalidMaxRoutines    = errors.New("Invalid maximum number of routines")
	ErrInvalidAdd            = errors.New("Invalid number of routines to add")
	ErrInvalidMetricHandler  = errors.New("Invalid metric handler")
	ErrInvalidMetricInterval = errors.New("Invalid metric interval")
)

//==============================================================================

// Worker must be implemented by types that want to use
// this worker processes.
type Worker interface {
	Work(context interface{}, id int)
}

// doWork is used internally to route work to the pool.
type doWork struct {
	context interface{}
	do      Worker
}

// Stat contains information about the pool.
type Stat struct {
	Routines    int64 // Current number of routines.
	Pending     int64 // Pending number of routines waiting to submit work.
	Active      int64 // Active number of routines in the work pool.
	Executed    int64 // Number of pieces of work executed.
	MaxRoutines int64 // High water mark of routines the pool has been at.
}

//==============================================================================

// OptEvent defines an handler used to provide events.
type OptEvent struct {
	Event func(context interface{}, event string, format string, a ...interface{})
}

// Config provides configuration for the pool.
type Config struct {
	MinRoutines func() int // Initial and minimum number of routines always in the pool.
	MaxRoutines func() int // Maximum number of routines we will ever grow the pool to.

	// *************************************************************************
	// ** Not Required, optional                                              **
	// *************************************************************************

	OptEvent
}

// Event fires events back to the user for important events.
func (cfg *Config) Event(context interface{}, event string, format string, a ...interface{}) {
	if cfg.OptEvent.Event != nil {
		cfg.OptEvent.Event(context, event, format, a...)
	}
}

//==============================================================================

// Pool provides a pool of routines that can execute any Worker
// tasks that are submitted.
type Pool struct {
	Config
	Name string // Name of this pool.

	tasks    chan doWork    // Unbuffered channel that work is sent into.
	control  chan int       // Unbuffered channel that work for the manager is send into.
	kill     chan struct{}  // Unbuffered channel to signal for a goroutine to die.
	shutdown chan struct{}  // Closed when the Work pool is being shutdown.
	wg       sync.WaitGroup // Manages the number of routines for shutdown.

	counter       int64 // Maintains a count of goroutines ever created to use as an id.
	updatePending int64 // Used to indicate a change to the pool is pending.

	muHealth sync.Mutex // Mutex used to check the health of the system safely.

	routines    int64 // Current number of routines.
	pending     int64 // Pending number of routines waiting to submit work.
	active      int64 // Active number of routines in the work pool.
	executed    int64 // Number of pieces of work executed.
	maxRoutines int64 // High water mark of routines the pool has been at.
}

// New creates a new Pool.
func New(context interface{}, name string, cfg Config) (*Pool, error) {
	if cfg.MinRoutines == nil {
		return nil, ErrNilMinRoutines
	}
	if cfg.MinRoutines() <= 0 {
		return nil, ErrInvalidMinRoutines
	}

	if cfg.MaxRoutines == nil {
		return nil, ErrNilMaxRoutines
	}
	if cfg.MaxRoutines() < cfg.MinRoutines() {
		return nil, ErrInvalidMaxRoutines
	}

	p := Pool{
		Config: cfg,
		Name:   name,

		tasks:    make(chan doWork),
		control:  make(chan int),
		kill:     make(chan struct{}),
		shutdown: make(chan struct{}),
	}

	p.manager(context)
	p.add(context, cfg.MinRoutines())

	return &p, nil
}

// Shutdown waits for all the workers to finish.
func (p *Pool) Shutdown(context interface{}) {
	// If a reset or change is being made, we need to wait.
	for atomic.LoadInt64(&p.updatePending) > 0 {
		time.Sleep(time.Second)
	}

	close(p.shutdown)
	p.wg.Wait()
}

// Do waits for the goroutine pool to take the work to be executed.
func (p *Pool) Do(context interface{}, work Worker) {
	dw := doWork{
		context: context,
		do:      work,
	}

	p.measureHealth()

	atomic.AddInt64(&p.pending, 1)
	p.tasks <- dw
	atomic.AddInt64(&p.pending, -1)
}

// DoWait waits for the goroutine pool to take the work to be executed or gives
// up after the alloted duration. Only use when you want to throw away work and
// not push back.
func (p *Pool) DoWait(context interface{}, work Worker, duration time.Duration) error {
	dw := doWork{
		context: context,
		do:      work,
	}

	p.measureHealth()

	atomic.AddInt64(&p.pending, 1)

	select {
	case p.tasks <- dw:
		atomic.AddInt64(&p.pending, -1)
		return nil

	case <-time.After(duration):
		atomic.AddInt64(&p.pending, -1)
		return errors.New("Timedout waiting to post work")
	}
}

// Stats returns the current snapshot of the pool stats.
func (p *Pool) Stats() Stat {
	return Stat{
		Routines:    atomic.LoadInt64(&p.routines),
		Pending:     atomic.LoadInt64(&p.pending),
		Active:      atomic.LoadInt64(&p.active),
		Executed:    atomic.LoadInt64(&p.executed),
		MaxRoutines: atomic.LoadInt64(&p.maxRoutines),
	}
}

// add creates routines to process work or sets a count for
// routines to terminate.
// NOTE: since our pools are auto-adjustable, we will not give the user ability
// to add routines.
func (p *Pool) add(context interface{}, routines int) error {
	if routines == 0 {
		return ErrInvalidAdd
	}

	cmd := addRoutine
	if routines < 0 {
		routines = routines * -1
		cmd = rmvRoutine
	}

	// Mark the number of adds or removes we are going to perform.
	atomic.AddInt64(&p.updatePending, int64(routines))

	for i := 0; i < routines; i++ {
		p.control <- cmd
	}

	return nil
}

// Reset re-adjusts the pool to match the specified number of routines.
// NOTE: since our pools are auto-adjustable, we will not give the user ability
// to reset the number of routines.
func (p *Pool) reset(context interface{}, routines int) {
	if routines < 0 {
		routines = 0
	}

	current := int(atomic.LoadInt64(&p.routines))
	p.add(context, routines-current)
}

// work performs the users work and keeps stats.
func (p *Pool) work(id int) {

	// Increment the number of routines.
	value := atomic.AddInt64(&p.routines, 1)

	// We need to check and set the high water mark.
	if value > atomic.LoadInt64(&p.maxRoutines) {
		atomic.StoreInt64(&p.maxRoutines, value)
	}

	// Decrement that the add command is complete.
	atomic.AddInt64(&p.updatePending, -1)

done:
	for {
		select {
		case dw := <-p.tasks:
			atomic.AddInt64(&p.active, 1)

			p.execute(id, dw)

			atomic.AddInt64(&p.active, -1)
			atomic.AddInt64(&p.executed, 1)

		case <-p.kill:
			break done
		}
	}

	// Decrement the number of routines.
	atomic.AddInt64(&p.routines, -1)

	// Decrement that the rmv command is complete.
	atomic.AddInt64(&p.updatePending, -1)

	p.wg.Done()
}

// execute performs the work in a recoverable way.
func (p *Pool) execute(id int, dw doWork) {
	defer func() {
		if r := recover(); r != nil {

			// Capture the stack trace
			buf := make([]byte, 10000)
			runtime.Stack(buf, false)

			p.Event(dw.context, "execute", "ERROR : %s", string(buf))
		}
	}()

	// Perform the work.
	dw.do.Work(dw.context, id)
}

// measureHealth calculates the health of the work pool.
func (p *Pool) measureHealth() {

	// If there are values pending to be updated, just
	// leave. We need those to finish first.
	if atomic.LoadInt64(&p.updatePending) > 0 {
		return
	}

	p.muHealth.Lock()
	defer p.muHealth.Unlock()

	stats := p.Stats()

	// We are not performing any work at all and we have more routines than min.
	if stats.Pending == 0 && stats.Active == 0 && (stats.Routines > int64(p.MinRoutines())) {

		// Reset the pool back to the min value.
		p.reset(p.Name, p.MinRoutines())
		return
	}

	// If we have no available routines at the moment and we have room to grow.
	if (stats.Routines == stats.Active) && (stats.Routines < int64(p.MaxRoutines())) {

		// Calculate the number of goroutines to add.
		add := int(float64(stats.Routines) * .20)

		// Check if we calculated a 0 grow.
		if add == 0 {
			add = 1
		}

		// Check if we will go over max.
		if (int(stats.Routines) + add) > p.MaxRoutines() {
			add = p.MaxRoutines() - int(stats.Routines)
		}

		// Request this number to be added.
		p.add(p.Name, add)
	}
}

// manager controls changes to the work pool including stats and shutting down.
func (p *Pool) manager(context interface{}) {
	p.wg.Add(1)

	go func() {
		for {
			select {
			case <-p.shutdown:

				// Capture the current number of routines.
				routines := int(atomic.LoadInt64(&p.routines))

				// Send a kill to all the existing routines.
				for i := 0; i < routines; i++ {
					p.kill <- struct{}{}
				}

				// Decrement the waitgroup and kill the manager.
				p.wg.Done()
				return

			case c := <-p.control:
				switch c {
				case addRoutine:

					// Capture the number of routines.
					routines := int(atomic.LoadInt64(&p.routines))

					// Is there room to add goroutines.
					if routines == p.MaxRoutines() {
						break
					}

					// Increment the total number of routines ever created.
					counter := atomic.AddInt64(&p.counter, 1)

					// Create the routine.
					p.wg.Add(1)
					go p.work(int(counter))

				case rmvRoutine:

					// Capture the number of routines.
					routines := int(atomic.LoadInt64(&p.routines))

					// Are there routines to remove.
					if routines <= p.MinRoutines() {
						break
					}

					// Send a kill signal to remove a routine.
					p.kill <- struct{}{}
				}
			}
		}
	}()
}
