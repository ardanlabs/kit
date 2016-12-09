// Package runner provide support for writing tasks that must complete
// within a certain duration or they must be killed. It also provides
// support for notifying the task the shutdown using a <control> C.
package runner

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"time"
)

// Error variables for the different states.
var (
	ErrTimeout  = errors.New("Timeout")
	ErrSignaled = errors.New("Signaled")
)

// Jobber defines an interface for providing the implementation details for
// processing a user job.
type Jobber interface {
	Job(traceID string) error
}

// Runner maintains state for the running process.
type Runner struct {
	shutdown chan struct{}
	sigChan  chan os.Signal
	kill     <-chan time.Time
	complete chan error
}

// New returns a new Runner value for use.
func New(timeout time.Duration) *Runner {
	return &Runner{
		shutdown: make(chan struct{}),
		sigChan:  make(chan os.Signal, 1),
		kill:     time.After(timeout),
		complete: make(chan error),
	}
}

// Run performs the execution of the specified job.
func (r *Runner) Run(traceID string, job Jobber) error {

	// We want to receive all interrupt based signals.
	signal.Notify(r.sigChan, os.Interrupt)

	// Launch the processor.
	go r.processor(traceID, job)

	for {
		select {
		case <-r.sigChan:

			// Close the channel to signal to the processor
			// it needs to shutdown.
			close(r.shutdown)

			// No need to process anymore events.
			signal.Stop(r.sigChan)

		case <-r.kill:

			// We have taken too much time. Kill the app hard.
			return ErrTimeout

		case err := <-r.complete:

			// Return the potential error.
			return err
		}
	}
}

// CheckShutdown can be used to check if a shutdown request has been issued.
func (r *Runner) CheckShutdown() bool {
	select {
	case <-r.shutdown:

		// We have been asked to shutdown.
		return true

	default:

		// We have not been asked to shutdown.
		return false
	}
}

// processor provides the main program logic for the program.
func (r *Runner) processor(traceID string, job Jobber) {

	// Variable to store any error that occurs.
	var err error

	// Defer the send on the channel so it happens
	// regardless of how this function terminates.
	defer func() {

		// Capture any potential panic.
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}

		// Signal the goroutine we have shutdown.
		r.complete <- err
	}()

	// Run the job.
	err = job.Job(traceID)
}
