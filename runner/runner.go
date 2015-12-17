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

// Error variables for the different states.
var (
	ErrTimeout  = errors.New("Timeout")
	ErrSignaled = errors.New("Signaled")
)

// Jobber defines an interface for providing the implementation details for
// processing a user job.
type Jobber interface {
	Job(context interface{}) error
}

// runner maintains state for the running process.
var runner struct {
	sync.Mutex

	shutdown chan struct{}
	sigChan  chan os.Signal
	kill     <-chan time.Time
	complete chan error

	recvShutdown bool
}

// Run performs the execution of the specified job.
func Run(context interface{}, timeout time.Duration, job Jobber) error {
	log.User(context, "Run", "Started : Timeout[%v]", timeout)

	// Init the runner for use.
	if initRunner(timeout) {
		return errors.New("Already running")
	}

	// When the task is done reset everything.
	defer resetRunner()

	// We want to receive all interrupt based signals.
	signal.Notify(runner.sigChan, os.Interrupt)

	// Launch the processor.
	go processor(context, job)

	for {
		select {
		case <-runner.sigChan:
			// Interrupt event signaled by the operation system.
			log.User(context, "Run", "Interrupt Received")

			// Flag we received the signal.
			runner.recvShutdown = true

			// Close the channel to signal to the processor
			// it needs to shutdown.
			close(runner.shutdown)

			// No need to process anymore events.
			signal.Stop(runner.sigChan)

		case <-runner.kill:
			// We have taken too much time. Kill the app hard.
			log.User(context, "Run", "Completed : Timedout")
			return ErrTimeout

		case err := <-runner.complete:
			// Everything completed within the time given.
			log.User(context, "Run", "Completed : Task Result : %v", err)

			if runner.recvShutdown {
				return ErrSignaled
			}

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

// resetRunner allows the runner to run a new task.
func resetRunner() {
	runner.Lock()
	{
		runner.shutdown = nil
		runner.sigChan = nil
		runner.kill = nil
		runner.complete = nil
		runner.recvShutdown = false
	}
	runner.Unlock()
}

// initRunner will check if a task is already running. If not
// it will initialize the runner to run a task.
func initRunner(timeout time.Duration) bool {
	runner.Lock()
	defer runner.Unlock()

	if runner.shutdown == nil {
		runner.shutdown = make(chan struct{})
		runner.sigChan = make(chan os.Signal, 1)
		runner.kill = time.After(timeout)
		runner.complete = make(chan error)
		return false
	}

	return true
}

// processor provides the main program logic for the program.
func processor(context interface{}, job Jobber) {
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
		runner.complete <- err
	}()

	// Run the job.
	err = job.Job(context)

	log.User(context, "processor", "Completed")
}
