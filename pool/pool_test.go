package pool_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/ardanlabs/kit/log"
	"github.com/ardanlabs/kit/pool"
	"github.com/ardanlabs/kit/tests"
)

func init() {
	tests.Init("KIT")
}

//==============================================================================

// theWork is the customer work type for using the pool.
type theWork struct {
	privateID int
}

// Work implements the DoWorker interface.
func (p *theWork) Work(context interface{}, id int) {
	log.Dev(context, "Work", "Performing Work with privateID %d\n", p.privateID)
}

// ExampleNewDoPool provides a basic example for using a DoPool.
func ExampleNewDoPool() {
	tests.ResetLog()
	defer tests.DisplayLog()

	// Create a configuration.
	cfg := pool.Config{
		MinRoutines: func() int { return 3 },
		MaxRoutines: func() int { return 4 },
	}

	// Create a new pool.
	p, err := pool.New("TEST", "TheWork", cfg)
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
	tests.ResetLog()
	defer tests.DisplayLog()

	// Create a configuration.
	cfg := pool.Config{
		MinRoutines: func() int { return 3 },
		MaxRoutines: func() int { return 4 },
	}

	// Create a new pool.
	p, err := pool.New("TEST", "TheWork", cfg)
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
	tests.ResetLog()
	defer tests.DisplayLog()

	t.Log("Given the need to validate the work pool functions.")
	{
		cfg := pool.Config{
			MinRoutines: func() int { return 100 },
			MaxRoutines: func() int { return 5000 },
		}

		p, err := pool.New("TestPool", "Pool1", cfg)
		if err != nil {
			t.Fatal("\tShould not get error creating pool.", tests.Failed, err)
		}
		t.Log("\tShould not get error creating pool.", tests.Success)

		for i := 0; i < 100; i++ {
			p.Do("TestPool", &theWork{privateID: i})
		}

		time.Sleep(100 * time.Millisecond)

		p.Shutdown("TestPool")
	}
}
