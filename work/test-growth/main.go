// SILLY PROGRAM THAT IS ALLOWING TO GUAGE HOW THINGS ARE PERFORMING.
// I STILL NEED TESTS, IF THAT IS POSSIBLE.
package main

import (
	"fmt"
	"time"

	"github.com/ardanlabs/kit/work"
)

// theWork is the customer work type for using the pool.
type theWork struct{}

// Work implements the DoWorker interface.
func (p *theWork) Work(context interface{}, id int) {
	//	time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
	time.Sleep(5 * time.Second)
}

func main() {
	// Create a configuration.
	config := work.Config{
		MinRoutines: func() int { return 100 },
		MaxRoutines: func() int { return 10000 },
	}

	// Create a new do pool.
	p, err := work.NewPool("TEST", "TheWork", &config)
	if err != nil {
		fmt.Println(err)
		return
	}

	go func() {
		for {
			stat := p.Stats()
			fmt.Printf("Check: %+v\n", stat)
			time.Sleep(time.Second)
		}
	}()

	const rate = 64 * 3 // ssps

	for j := 0; j < 3; j++ {
		for b := 0; b < rate; b++ {
			time.Sleep(time.Second / rate)
			p.Do("TEST", &theWork{})
		}

		// I want the pool to reset and start again.
		for {
			stat := p.Stats()
			if stat.Active == 0 {
				break
			}
			time.Sleep(time.Second)
		}
	}

	time.Sleep(time.Hour)
}
