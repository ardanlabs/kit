// Use gin during development : // go get github.com/codegangsta/gin
// Run this command in the folder: gin -p 5000 -a 4000 -i run

package main

import (
	"time"

	"github.com/ardanlabs/kit/examples/web/routes"
	"github.com/ardanlabs/kit/log"
	"github.com/ardanlabs/kit/web"
)

func main() {

	// Look at /kit/config for a set of possible config settings.

	err := web.Run(":4000", routes.API(), 10*time.Second, 10*time.Second)
	log.User("main", "main", "DOWN: %v", err)
}
