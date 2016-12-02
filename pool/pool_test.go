package pool_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/ardanlabs/kit/pool"
)

// Success and failure markers.
var (
	success = "\u2713"
	failed  = "\u2717"
)

// theWork is the customer work type for using the pool.
type theWork struct {
	privateID int
}

// Work implements the DoWorker interface.
func (p *theWork) Work(ctx interface{}, id int) {
	fmt.Printf("Performing Work with privateID %d\n", p.privateID)
}

// ExampleNew provides a basic example for using a pool.
func ExampleNew() {
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

// TestPool tests the pool is functional.
func TestPool(t *testing.T) {
	t.Log("Given the need to validate the work pool functions.")
	{
		cfg := pool.Config{
			MinRoutines: func() int { return 100 },
			MaxRoutines: func() int { return 5000 },
		}

		p, err := pool.New("TestPool", "Pool1", cfg)
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
