// Use gin during development : // go get github.com/codegangsta/gin
// Run this command in the folder: gin -p 5000 -a 4000 -i run

package main

import (
	"github.com/ardanlabs/kit/examples/web/routes"

	"github.com/ardanlabs/kit/web/app"
)

func main() {
	// Look at /kit/config for a set of possible config settings.

	app.Run(":4000", routes.API())
}
