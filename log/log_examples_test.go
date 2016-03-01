package log_test

import (
	"errors"
	"os"
	"testing"

	"github.com/ardanlabs/kit/log"
)

// ExampleDev shows how to use the log package.
func ExampleDev(t *testing.T) {

	// Init the log package for stdout. Hardcode the logging level
	// function to use USER level logging.
	log.Init(os.Stdout, func() int { return log.USER })

	// Write a simple log line with no formatting.
	log.User("context", "ExampleDev", "This is a simple line with no formatting")

	// Write a simple log line with formatting.
	log.User("context", "ExampleDev", "This is a simple line with no formatting %d", 10)

	// Write a message error for the user.
	log.Error("context", "ExampleDev", errors.New("A user error"), "testing error")

	// Write a message error for the user with formatting.
	log.Error("context", "ExampleDev", errors.New("A user error"), "testing error %s", "value")

	// Write a message error for the developer only.
	log.Dev("context", "ExampleDev", "Formatting %v", 42)

	// Write a simple log line with no formatting.
	log.UserOffset("context", 1, "ExampleDev", "This is a simple line with no formatting")

	// Write a simple log line with formatting.
	log.UserOffset("context", 1, "ExampleDev", "This is a simple line with no formatting %d", 10)

	// Write a message error for the user.
	log.ErrorOffset("context", 1, "ExampleDev", errors.New("A user error"), "testing error")

	// Write a message error for the user with formatting.
	log.ErrorOffset("context", 1, "ExampleDev", errors.New("A user error"), "testing error %s", "value")

	// Write a message error for the developer only.
	log.DevOffset("context", 1, "ExampleDev", "Formatting %v", 42)

}
