package log_test

import (
	"errors"
	"os"

	"github.com/ardanlabs/kit/log"
)

// ExampleDev shows how to use the log package.
func ExampleDev() {

	// Init the log package for stdout. Hardcode the logging level
	// function to use USER level logging.
	log.Init(os.Stdout, func() int { return log.USER }, log.Ldefault)

	// Write a simple log line with no formatting.
	log.User("ctx", "ExampleDev", "This is a simple line with no formatting")

	// Write a simple log line with formatting.
	log.User("ctx", "ExampleDev", "This is a simple line with no formatting %d", 10)

	// Write a message error for the user.
	log.Error("ctx", "ExampleDev", errors.New("A user error"), "testing error")

	// Write a message error for the user with formatting.
	log.Error("ctx", "ExampleDev", errors.New("A user error"), "testing error %s", "value")

	// Write a message error for the developer only.
	log.Dev("ctx", "ExampleDev", "Formatting %v", 42)

	// Write a simple log line with no formatting.
	log.UserOffset("ctx", 3, "ExampleDev", "This is a simple line with no formatting")

	// Write a simple log line with formatting.
	log.UserOffset("ctx", 3, "ExampleDev", "This is a simple line with no formatting %d", 10)

	// Write a message error for the user.
	log.ErrorOffset("ctx", 3, "ExampleDev", errors.New("A user error"), "testing error")

	// Write a message error for the user with formatting.
	log.ErrorOffset("ctx", 3, "ExampleDev", errors.New("A user error"), "testing error %s", "value")

	// Write a message error for the developer only.
	log.DevOffset("ctx", 3, "ExampleDev", "Formatting %v", 42)

}
