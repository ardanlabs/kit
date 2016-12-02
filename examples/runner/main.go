// Sample program to show how to use the runner package to build tasks
// that must run within a well defined duration.
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/log"
	"github.com/ardanlabs/kit/runner"
)

// Configuation settings.
const (
	configKey       = "KIT"
	cfgLoggingLevel = "LOGGING_LEVEL"
)

func init() {

	// This is being added to showcase configuration.
	os.Setenv("KIT_LOGGING_LEVEL", "1")

	// Init the configuration system.
	if err := cfg.Init(cfg.EnvProvider{Namespace: configKey}); err != nil {
		fmt.Println("Error initalizing configuration system", err)
		os.Exit(1)
	}

	// Init the log system.
	logLevel := func() int {
		ll, err := cfg.Int(cfgLoggingLevel)
		if err != nil {
			return log.USER
		}
		return ll
	}
	log.Init(os.Stderr, logLevel, log.Ldefault)

	// Log all the configuration options
	log.User("startup", "init", "\n\nConfig Settings: %s\n%s\n", configKey, cfg.Log())
}

//==============================================================================

// Task represents a task we need to run.
type Task struct {
	Name string
}

// Job implements the Jobber interface so task can be managed.
func (t *Task) Job(ctx interface{}) error {
	log.User(ctx, "Job", "Started : **********")

	time.Sleep(time.Second)

	log.User(ctx, "Job", "Completed : **********")
	return nil
}

//==============================================================================

func main() {
	const ctx = "main"

	// Create a task value for execution.
	t := Task{
		Name: "test task",
	}

	rn := runner.New(time.Second)

	// Start the job running with a specified duration.
	if err := rn.Run(ctx, &t); err != nil {
		switch err {
		case runner.ErrTimeout:

			// The task did not finish within the specified duration.
			log.Error(ctx, "main", err, "Task timeout")

		case runner.ErrSignaled:

			// The user hit <control> c and we shutdown early.
			log.Error(ctx, "main", err, "Shutdown early")

		default:

			// An error occurred in the processing of the task.
			log.Error(ctx, "main", err, "Processing error")
		}

		os.Exit(1)
	}

	log.User(ctx, "main", "Completed")
}
