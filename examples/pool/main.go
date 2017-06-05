// Sample program to show how to use the pool package to build pools
// of goroutines to get work done.
package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/log"
	"github.com/ardanlabs/kit/pool"
)

// Configuation settings.
const (
	configKey       = "KIT"
	cfgLoggingLevel = "LOGGING_LEVEL"
	cfgMinRoutines  = "MIN_ROUTINES"
	cfgMaxRoutines  = "MAX_ROUTINES"
)

// TraceID is represents the trace id.
type TraceID string

// TraceIDKey is the type of value to use for the key. The key is
// type specific and only values of the same type will match.
type TraceIDKey int

func init() {

	// This is being added to showcase configuration.
	os.Setenv("KIT_LOGGING_LEVEL", "1")
	os.Setenv("KIT_MIN_ROUTINES", "1")
	os.Setenv("KIT_MAX_ROUTINES", "10")

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

// wg represents a WaitGroup so we can control the shutdown
// of the test application.
var wg sync.WaitGroup

// Task represents a task we need to run.
type Task struct {
	Name string
}

// Work implements the Worker interface so task can be executed by the pool.
func (t *Task) Work(ctx context.Context, id int) {
	time.Sleep(time.Second)
	wg.Done()
}

func main() {

	// Create a traceID for this run.
	traceID := TraceID("f47ac10b-58cc-0372-8567-0e02b2c3d479")
	const traceIDKey TraceIDKey = 0
	ctx := context.WithValue(context.Background(), traceIDKey, traceID)

	// Create the configuration.
	cfg := pool.Config{
		MinRoutines: func() int { return cfg.MustInt(cfgMinRoutines) },
		MaxRoutines: func() int { return cfg.MustInt(cfgMaxRoutines) },
	}

	// Create a pool.
	p, err := pool.New("test", cfg)
	if err != nil {
		log.Error(string(traceID), "main", err, "Creating pool")
		return
	}

	// Look at stats for the work.
	go func() {
		for {
			time.Sleep(250 * time.Millisecond)
			log.User(string(traceID), "Stats", "%#v", p.Stats())
		}
	}()

	const totalWork = 100
	wg.Add(totalWork)

	// Perform some work.
	for i := 0; i < totalWork; i++ {
		p.Do(ctx, &Task{Name: strconv.Itoa(i)})
	}

	// Wait until all the work is complete.
	wg.Wait()

	// Shutdown the pool.
	p.Shutdown()
}
