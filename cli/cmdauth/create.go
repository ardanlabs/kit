package cmdauth

import (
	"time"

	"github.com/ardanlabs/kit/auth"
	"github.com/ardanlabs/kit/db"

	"github.com/spf13/cobra"
)

var createLong = `Use create to add a new user to the system. The user email
must be unique for every user.

Example:
  ./kit auth create -n "Bill Kennedy" -e "bill@ardanlabs.com" -p "yefc*7fdf92"
`

// create contains the state for this command.
var create struct {
	name  string
	pass  string
	email string
}

// addCreate handles the creation of users.
func addCreate() {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Add a new user to the system.",
		Long:  createLong,
		Run:   runCreate,
	}

	cmd.Flags().StringVarP(&create.name, "name", "n", "", "Full name for the user.")
	cmd.Flags().StringVarP(&create.email, "email", "e", "", "Email for the user.")
	cmd.Flags().StringVarP(&create.pass, "pass", "p", "", "Password for the user.")

	userCmd.AddCommand(cmd)
}

// runCreate is the code that implements the create command.
func runCreate(cmd *cobra.Command, args []string) {
	cmd.Printf("Creating User : Name[%s] Email[%s] Pass[%s]\n", create.name, create.email, create.pass)

	if create.name == "" && create.email == "" && create.pass == "" {
		cmd.Help()
		return
	}

	u, err := auth.NewUser(auth.NUser{
		Status:   auth.StatusActive,
		FullName: create.name,
		Email:    create.email,
		Password: create.pass,
	})
	if err != nil {
		cmd.Println("Creating User : ", err)
		return
	}

	db := db.NewMGO()
	defer db.CloseMGO()

	if err := auth.CreateUser("", db, u); err != nil {
		cmd.Println("Creating User : ", err)
		return
	}

	webTok, err := auth.CreateWebToken("", db, u, 24*365*time.Hour)
	if err != nil {
		cmd.Println("Creating User : ", err)
		return
	}

	cmd.Printf("\nToken: %s\n\n", webTok)
}
