package cmdauth

import (
	"github.com/ardanlabs/kit/auth"
	"github.com/ardanlabs/kit/db"

	"github.com/spf13/cobra"
)

var statusLong = `Use status to change the status of a user to active or disabled.

Note: Not including the status results in the user being disabled.

Example:
  ./kit auth status -e "bill@ardanlabs.com" -a true
`

// status contains the state for this command.
var status struct {
	pid    string
	email  string
	active bool
}

//==============================================================================

// addStatus handles the deletion of user records.
func addStatus() {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Changes the status of a user.",
		Long:  statusLong,
		Run:   runStatus,
	}

	cmd.Flags().StringVarP(&status.pid, "public_id", "p", "", "Public Id of the user.")
	cmd.Flags().StringVarP(&status.email, "email", "e", "", "Email of the user.")
	cmd.Flags().BoolVarP(&status.active, "active", "a", false, "Use `true` if active or `false` if not.")

	userCmd.AddCommand(cmd)
}

// runStatus is the code that implements the status command.
func runStatus(cmd *cobra.Command, args []string) {
	cmd.Printf("Status User : Pid[%s] Email[%s] Active[%v]\n", status.pid, status.email, status.active)

	if status.pid == "" && status.email == "" {
		cmd.Help()
		return
	}

	db := db.NewMGO()
	defer db.CloseMGO()

	var publicID string
	if status.pid != "" {
		publicID = status.pid
	} else {
		u, err := auth.GetUserByEmail("", db, status.email)
		if err != nil {
			cmd.Println("Status User : ", err)
			return
		}
		publicID = u.PublicID
	}

	st := auth.StatusDisabled
	if status.active {
		st = auth.StatusActive
	}

	if err := auth.UpdateUserStatus("", db, publicID, st); err != nil {
		cmd.Println("Status User : ", err)
		return
	}

	cmd.Println("Status User : Updated")
}
