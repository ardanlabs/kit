package cmdauth

import "github.com/spf13/cobra"

// userCmd represents the parent for all cli commands.
var userCmd = &cobra.Command{
	Use:   "auth",
	Short: "auth provides managing user records.",
}

// GetCommands returns the user commands.
func GetCommands() *cobra.Command {
	addCreate()
	addGet()
	addStatus()
	return userCmd
}
