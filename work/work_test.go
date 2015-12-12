package work_test

import (
	"bytes"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/ardanlabs/kit/work"
)

// Success and failure markers.
var (
	success = "\u2713"
	failed  = "\u2717"
)

// logdash is the central buffer where all logs are stored.
var logdash bytes.Buffer

//==============================================================================

// resetLog resets the contents of logdash.
func resetLog() {
	logdash.Reset()
}

// displayLog writes the logdash data to standand out, if testing in verbose mode
// was turned on.
func displayLog() {
	if !testing.Verbose() {
		return
	}

	logdash.WriteTo(os.Stdout)
}

//==============================================================================

// theWork is the customer work type for using the pool.
type theWork struct {
	privateID int
}

// Work implements the DoWorker interface.
func (p *theWork) Work(context interface{}, id int) {
	logdash.WriteString(fmt.Sprintf("Performing Work with privateID %d\n", p.privateID))
}

// ExampleNewDoPool provides a basic example for using a DoPool.
func ExampleNewDoPool() {
	resetLog()
	defer displayLog()

	// Create a configuration.
	config := work.Config{
		MinRoutines: func() int { return 3 },
		MaxRoutines: func() int { return 4 },
	}

	// Create a new do pool.
	p, err := work.NewPool("TEST", "TheWork", &config)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Pass in some work to be performed.
	p.Do("TEST", &theWork{})
	p.Do("TEST", &theWork{})
	p.Do("TEST", &theWork{})

	// Wait to the work to be processed.
	time.Sleep(100 * time.Millisecond)

	// Shutdown the pool.
	p.Shutdown("TEST")
}

// ExampleMetrics provides an example of the metrics handler being called.
func ExampleMetrics() {
	resetLog()
	defer displayLog()

	// Create a configuration.
	config := work.Config{
		MinRoutines: func() int { return 3 },
		MaxRoutines: func() int { return 4 },
	}

	// Create a new do pool.
	p, err := work.NewPool("TEST", "TheWork", &config)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Pass in one piece of work to be performed.
	p.Do("TEST", &theWork{})

	// For the example, wait for it to display
	time.Sleep(time.Second)

	// Shutdown the pool.
	p.Shutdown("TEST")
}

// TestPool tests the pool is functional.
func TestPool(t *testing.T) {
	resetLog()
	defer displayLog()

	t.Log("Given the need to validate the work pool functions.")
	{
		cfg := work.Config{
			MinRoutines: func() int { return 100 },
			MaxRoutines: func() int { return 5000 },
		}

		p, err := work.NewPool("TestPool", "Pool1", &cfg)
		if err != nil {
			t.Fatal("\tShould not get error creating pool.", failed, err)
		}
		t.Log("\tShould not get error creating pool.", success)

		for i := 0; i < 100; i++ {
			p.Do("TestPool", &theWork{privateID: i})
		}

		time.Sleep(100 * time.Millisecond)

		p.Shutdown("TestPool")
	}
}
