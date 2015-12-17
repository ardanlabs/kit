// Package runner provide support for writing tasks that must complete
// within a certain duration or they must be killed. It also provides
// support for notifying the task the shutdown using a <control> C.
package runner

import (
	"errors"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/ardanlabs/kit/log"
)

// ErrTimeout is returned when the task timesout.
var ErrTimeout = errors.New("Timeout")

// Jobber defines an interface for providing the implementation details for
// processing a user job.
type Jobber interface {
	Job(context interface{}) error
}

// runner maintains state for the running process.
var runner struct {
	sync.Mutex
	shutdown chan struct{}
}

// Run performs the execution of the specified job.
func Run(context interface{}, timeout time.Duration, job Jobber) error {
	log.User(context, "Run", "Started : Timeout[%v]", timeout)

	// Had Run already been called.
	running := true
	runner.Lock()
	{
		if runner.shutdown == nil {
			runner.shutdown = make(chan struct{})
			running = false
		}
	}
	runner.Unlock()

	if running {
		return errors.New("Already running")
	}

	// Initialize the local channels.
	var (
		sigChan  = make(chan os.Signal, 1)
		kill     = time.After(timeout)
		complete = make(chan error)
	)

	// We want to receive all interrupt based signals.
	signal.Notify(sigChan, os.Interrupt)

	// Launch the processor.
	go processor(context, job, complete)

	for {
		select {
		case <-sigChan:
			// Interrupt event signaled by the operation system.
			log.User(context, "Run", "Interrupt Received")

			// Close the channel to signal to the processor
			// it needs to shutdown.
			close(runner.shutdown)

			// No need to process anymore events.
			signal.Stop(sigChan)

		case <-kill:
			// We have taken too much time. Kill the app hard.
			log.User(context, "Run", "Completed : Timedout")
			return ErrTimeout

		case err := <-complete:
			// Everything completed within the time given.
			log.User(context, "Run", "Completed : Task Result : %s", err)
			return err
		}
	}
}

// CheckShutdown checks the shutdown flag to determine
// if we have been asked to interrupt processing.
func CheckShutdown(context interface{}) bool {
	select {
	case <-runner.shutdown:
		// We have been asked to shutdown.
		log.User(context, "CheckShutdown", "Shutdown Early")
		return true

	default:
		// We have not been asked to shutdown.
		return false
	}
}

// processor provides the main program logic for the program.
func processor(context interface{}, job Jobber, complete chan<- error) {
	log.User(context, "processor", "Started")

	// Variable to store any error that occurs.
	var err error

	// Defer the send on the channel so it happens
	// regardless of how this function terminates.
	defer func() {
		// Capture any potential panic.
		if r := recover(); r != nil {
			log.User(context, "processor", "Panic : %v", r)
		}

		// Signal the goroutine we have shutdown.
		complete <- err
	}()

	// Run the job.
	err = job.Job(context)

	log.User(context, "processor", "Completed")
}
