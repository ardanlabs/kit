package cmddb

import "github.com/spf13/cobra"

// dbCmd represents the parent for all database cli commands.
var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "db will create a kit database and validate everything exists.",
}

// GetCommands returns the query commands.
func GetCommands() *cobra.Command {
	addCreate()
	return dbCmd
}
