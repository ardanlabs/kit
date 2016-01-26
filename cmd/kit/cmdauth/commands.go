package cmdauth

import "github.com/spf13/cobra"

// userCmd represents the parent for all cli commands.
var userCmd = &cobra.Command{
	Use:   "auth",
	Short: "auth provides managing user records.",
}

// mgoSession holds the master session for the DB access.
var mgoSession string

// GetCommands returns the user commands.
func GetCommands(mgoSes string) *cobra.Command {
	mgoSession = mgoSes

	addCreate()
	addGet()
	addStatus()
	return userCmd
}
