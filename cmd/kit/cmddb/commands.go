package cmddb

import "github.com/spf13/cobra"

// dbCmd represents the parent for all database cli commands.
var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "db will create a kit database and validate everything exists.",
}

// mgoSession holds the master session for the DB access.
var mgoSession string

// GetCommands returns the query commands.
func GetCommands(mgoSes string) *cobra.Command {
	mgoSession = mgoSes

	addCreate()
	return dbCmd
}
