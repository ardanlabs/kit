// Use gin during development : // go get github.com/codegangsta/gin
// Run this command in the folder: gin -p 5000 -a 4000 -i run

package main

import (
	"github.com/ardanlabs/kit/examples/web/routes"

	"github.com/ardanlabs/kit/web/app"
)

/*
// If using MongoDB set the env variables. ENV_PREFIX can any
// prefix value of your choice.

export ENV_PREFIX_MONGO_HOST=
export ENV_PREFIX_MONGO_USER=
export ENV_PREFIX_MONGO_AUTHDB=
export ENV_PREFIX_MONGO_DB=
export ENV_PREFIX_MONGO_PASS=

// Use this to adjust the logging level
// 0 - None, 1 - Dev, 2 - User

export ENV_PREFIX_LOGGING_LEVEL=1
*/

func main() {
	// If ENV_HOST is not found then :4000 is used.
	app.Run(":4000", routes.API())
}
