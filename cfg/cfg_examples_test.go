package cfg_test

import (
	"fmt"
	"os"
	"time"

	"github.com/coralproject/shelf/pkg/cfg"
)

// ExampleDev shows how to use the config package.
func ExampleDev() {
	// Set up some basic environment variables.
	os.Setenv("DOCK_IP", "40.23.233.10")
	os.Setenv("DOCK_PORT", "4044")
	os.Setenv("DOCK_INIT_STAMP", time.Date(2009, time.November,
		10, 15, 0, 0, 0, time.UTC).UTC().Format(time.UnixDate))
	os.Setenv("DOCK_FLAG", "true")

	// Init() must be called only once with the given namespace to load.
	cfg.Init("DOCK")

	// To get the ip.
	fmt.Println(cfg.MustString("IP"))

	// To get the port number.
	fmt.Println(cfg.MustInt("PORT"))

	// To get the timestamp.
	fmt.Println(cfg.MustTime("INIT_STAMP"))

	// To get the flag.
	fmt.Println(cfg.MustBool("FLAG"))

	// Output:
	// 40.23.233.10
	// 4044
	// 2009-11-10 15:00:00 +0000 UTC
	// true
}
