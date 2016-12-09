package runner_test

import (
	"errors"
	"runtime"
	"syscall"
	"testing"
	"time"

	"github.com/ardanlabs/kit/runner"
)

// Success and failure markers.
var (
	success = "\u2713"
	failed  = "\u2717"
)

//==============================================================================

// task represents a test task.
type task struct {
	kill chan bool
	err  error
}

// Job is the implementation of the Jobber interface.
func (t *task) Job(traceID string) error {
	// Pretend you are doing work for the specified
	// amount of time.
	<-t.kill

	// Report we received the signal to keep things in
	// sync between test functions.
	t.kill <- true

	return t.err
}

// Kill will kill the Job method immediately.
func (t *task) Kill() {
	select {
	case t.kill <- true:
		// If we were able to send the message, wait
		// for the response to keep things in sync.
		<-t.kill
	default:
	}
}

// KillAfter will kill the Job method after the specified duration.
func (t *task) KillAfter(dur time.Duration) {
	t.kill = make(chan bool)

	go func() {
		time.Sleep(dur)
		t.Kill()
	}()

	runtime.Gosched()
}

//==============================================================================

// TestCompleted tests when jobs complete properly.
func TestCompleted(t *testing.T) {
	t.Log("Given the need to test a successful task run.")
	{
		t.Log("\tWhen using a task that will complete in time.")
		{
			var job task
			job.KillAfter(time.Millisecond)

			r := runner.New(time.Second)

			if err := r.Run("traceID", &job); err != nil {
				t.Fatalf("\t%s\tShould not receive an error : %v", failed, err)
			}
			t.Logf("\t%s\tShould not receive an error.", success)
		}
	}
}

// TestError tests when jobs complete properly but with errors.
func TestError(t *testing.T) {
	t.Log("Given the need to test a successful task run with error.")
	{
		t.Log("\tWhen using a task that will complete in time.")
		{
			job := task{
				err: errors.New("Error"),
			}
			job.KillAfter(time.Millisecond)

			r := runner.New(time.Second)

			if err := r.Run("traceID", &job); err == nil {
				t.Fatalf("\t%s\tShould receive our error : %v", failed, err)
			}
			t.Logf("\t%s\tShould receive our error.", success)
		}
	}
}

// TestTimeout tests when jobs timeout.
func TestTimeout(t *testing.T) {
	t.Log("Given the need to test a task that timesout.")
	{
		t.Log("\tWhen using a task that will timeout.")
		{
			var job task
			job.KillAfter(time.Second)

			// Need the job method to quit as soon as we are done.
			defer job.Kill()

			r := runner.New(time.Millisecond)

			if err := r.Run("traceID", &job); err != runner.ErrTimeout {
				t.Fatalf("\t%s\tShould receive a timeout error : %v", failed, err)
			}
			t.Logf("\t%s\tShould receive a timeout error.", success)
		}
	}
}

// TestSignaled tests when jobs is requested to shutdown.
func TestSignaled(t *testing.T) {
	t.Log("Given the need to test a task that is requested to shutdown.")
	{
		t.Log("\tWhen using a task that should see the signal.")
		{
			var job task
			job.KillAfter(100 * time.Millisecond)

			// Need the job method to quit as soon as we are done.
			defer job.Kill()

			go func() {
				time.Sleep(50 * time.Millisecond)
				syscall.Kill(syscall.Getpid(), syscall.SIGINT)
			}()

			r := runner.New(3 * time.Second)

			if err := r.Run("traceID", &job); err != nil {
				t.Errorf("\t%s\tShould receive no error : %v", failed, err)
			} else {
				t.Logf("\t%s\tShould receive no error.", success)
			}

			if !r.CheckShutdown() {
				t.Errorf("\t%s\tShould show the check shutdown flag is set.", failed)
			} else {
				t.Logf("\t%s\tShould show the check shutdown flag is set.", success)
			}
		}
	}
}
