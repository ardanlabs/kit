package cmdauth

import (
	"encoding/json"

	"github.com/ardanlabs/kit/auth"
	"github.com/spf13/cobra"
)

var listLong = `Use list to retrieve multiple user records from the system.

Example:
  ./kit auth list

  ./kit auth list -a
`

// list contains the state for this command.
var list struct {
	activeOnly bool
}

//==============================================================================

// addList handles the retrieval users records, displayed in JSON formatted response.
func addList() {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Retrieves a user record by public_id or email.",
		Long:  listLong,
		Run:   runlist,
	}

	cmd.Flags().BoolVarP(&list.activeOnly, "active", "a", false, "Limit to only active users")

	userCmd.AddCommand(cmd)
}

// runlist is the code that implements the list command.
func runlist(cmd *cobra.Command, args []string) {
	cmd.Printf("Listing Users : Active Only[%t]\n", list.activeOnly)

	users, err := auth.GetUsers("", conn, list.activeOnly)

	if err != nil {
		cmd.Println("List Users : ", err)
		return
	}

	for i, u := range users {
		webTok, err := auth.GetUserWebToken("", conn, u.PublicID)
		if err != nil {
			webTok = err.Error()
		}

		data, err := json.MarshalIndent(&u, "", "    ")
		if err != nil {
			cmd.Println("List Users : ", err)
			return
		}

		cmd.Printf("\nUser %d\n%s\n\nToken: %s\n\n", i, string(data), webTok)
	}
	return
}
