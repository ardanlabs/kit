package cmdauth

import (
	"github.com/ardanlabs/kit/db"
	"github.com/spf13/cobra"
)

// userCmd represents the parent for all cli commands.
var userCmd = &cobra.Command{
	Use:   "auth",
	Short: "auth provides managing user records.",
}

// conn holds the session for the DB access.
var conn *db.DB

// GetCommands returns the auth commands.
func GetCommands(db *db.DB) *cobra.Command {
	conn = db

	addCreate()
	addGet()
	addList()
	addStatus()
	return userCmd
}
