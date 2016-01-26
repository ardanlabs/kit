package auth_test

import (
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/ardanlabs/kit/auth"
	"github.com/ardanlabs/kit/auth/session"
	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/db/mongo"
	"github.com/ardanlabs/kit/tests"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func init() {
	os.Setenv("KIT_LOGGING_LEVEL", "1")

	cfg := mongo.Config{
		Host:     "ds027155.mongolab.com:27155",
		AuthDB:   "kit",
		DB:       "kit",
		User:     "kit",
		Password: "community",
	}

	tests.Init("KIT")
	tests.InitMongo(cfg)

	ensureIndexes()
}

//==============================================================================

// TestModelInvalidation tests things that can fail with model validation.
func TestModelInvalidation(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	t.Log("Given the need to validate models will invalidate with bad data.")
	{
		t.Log("\tWhen using a test user.")
		{
			var nu auth.NUser

			if err := nu.Validate(); err == nil {
				t.Errorf("\t%s\tShould Not be able to validate NUser value.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould Not be able to validate NUser value.", tests.Success)
			}

			if _, err := auth.NewUser(nu); err == nil {
				t.Errorf("\t%s\tShould Not be able to create a new user value.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould Not be able to create a new user value.", tests.Success)
			}

			var u auth.User

			if _, err := u.Pwd(); err == nil {
				t.Errorf("\t%s\tShould Not be able to call Pwd is empty password.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould Not be able to call Pwd is empty password.", tests.Success)
			}

			if _, err := u.Salt(); err == nil {
				t.Errorf("\t%s\tShould Not be able to call Salt with bad values.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould Not be able to call Salt with bad values.", tests.Success)
			}

			if _, err := u.WebToken(""); err == nil {
				t.Errorf("\t%s\tShould Not be able to call WebToken with bad values.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould Not be able to call WebToken with bad values.", tests.Success)
			}

			if ok := u.IsPasswordValid(""); ok {
				t.Errorf("\t%s\tShould Not be able to call IsPasswordValid with empty password.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould Not be able to call IsPasswordValid with empty password.", tests.Success)
			}

			u.Password = "123"

			if ok := u.IsPasswordValid(""); ok {
				t.Errorf("\t%s\tShould Not be able to call IsPasswordValid with bad password.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould Not be able to call IsPasswordValid with bad password.", tests.Success)
			}
		}
	}
}

// TestCreateUser tests the creation of a user.
func TestCreateUser(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	db, err := db.NewMGO(tests.Context, tests.TestSession)
	if err != nil {
		t.Fatalf("\t%s\tShould be able to get a Mongo session : %v", tests.Failed, err)
	}
	defer db.CloseMGO(tests.Context)

	var publicID string
	defer func() {
		if err := removeUser(db, publicID); err != nil {
			t.Fatalf("\t%s\tShould be able to remove the test user : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to remove the test user.", tests.Success)
	}()

	t.Log("Given the need to create users in the DB.")
	{
		t.Log("\tWhen using a test user.")
		{
			u1, err := auth.NewUser(auth.NUser{
				Status:   auth.StatusActive,
				FullName: "Test Kennedy",
				Email:    "bill@ardanlabs.com",
				Password: "_Password124",
			})
			if err != nil {
				t.Fatalf("\t%s\tShould be able to build a new user : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to build a new user.", tests.Success)

			if err := auth.CreateUser(tests.Context, db, u1); err != nil {
				t.Fatalf("\t%s\tShould be able to create a user : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a user.", tests.Success)

			// We need to do this so we can clean up after.
			publicID = u1.PublicID

			u2, err := auth.GetUserByPublicID(tests.Context, db, u1.PublicID, true)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve the user by PublicID : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to retrieve the user by PublicID.", tests.Success)

			// Remove the objectid to be able to compare the values.
			u2.ID = ""

			// Need to remove the nanoseconds to be able to compare the values.
			u1.DateModified = u1.DateModified.Add(-time.Duration(u1.DateModified.Nanosecond()))
			u1.DateCreated = u1.DateCreated.Add(-time.Duration(u1.DateCreated.Nanosecond()))
			u2.DateModified = u2.DateModified.Add(-time.Duration(u2.DateModified.Nanosecond()))
			u2.DateCreated = u2.DateCreated.Add(-time.Duration(u2.DateCreated.Nanosecond()))

			if !reflect.DeepEqual(*u1, *u2) {
				t.Logf("\t%+v", *u1)
				t.Logf("\t%+v", *u2)
				t.Fatalf("\t%s\tShould be able to get back the same user.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould be able to get back the same user.", tests.Success)
			}

			u3, err := auth.GetUserByEmail(tests.Context, db, u1.Email, true)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve the user by Email : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to retrieve the user by Email.", tests.Success)

			// Remove the objectid to be able to compare the values.
			u3.ID = ""

			// Need to remove the nanoseconds to be able to compare the values.
			u3.DateModified = u3.DateModified.Add(-time.Duration(u3.DateModified.Nanosecond()))
			u3.DateCreated = u3.DateCreated.Add(-time.Duration(u3.DateCreated.Nanosecond()))

			if !reflect.DeepEqual(*u1, *u3) {
				t.Logf("\t%+v", *u1)
				t.Logf("\t%+v", *u3)
				t.Fatalf("\t%s\tShould be able to get back the same user.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould be able to get back the same user.", tests.Success)
			}
		}
	}
}

// TestCreateUserTwice tests the creation of the same user fails. This test
// requires an index on the collection.
func TestCreateUserTwice(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	db, err := db.NewMGO(tests.Context, tests.TestSession)
	if err != nil {
		t.Fatalf("\t%s\tShould be able to get a Mongo session : %v", tests.Failed, err)
	}
	defer db.CloseMGO(tests.Context)

	var publicID string
	defer func() {
		if err := removeUser(db, publicID); err != nil {
			t.Fatalf("\t%s\tShould be able to remove the test user : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to remove the test user.", tests.Success)
	}()

	t.Log("Given the need to make sure the same user can't be created twice.")
	{
		t.Log("\tWhen using a test user.")
		{
			u1, err := auth.NewUser(auth.NUser{
				Status:   auth.StatusActive,
				FullName: "Test Kennedy",
				Email:    "bill@ardanlabs.com",
				Password: "_Password124",
			})
			if err != nil {
				t.Fatalf("\t%s\tShould be able to build a new user : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to build a new user.", tests.Success)

			if err := auth.CreateUser(tests.Context, db, u1); err != nil {
				t.Fatalf("\t%s\tShould be able to create a user : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a user.", tests.Success)

			// We need to do this so we can clean up after.
			publicID = u1.PublicID

			if err := auth.CreateUser(tests.Context, db, u1); err == nil {
				t.Fatalf("\t%s\tShould Not be able to create a user", tests.Failed)
			}
			t.Logf("\t%s\tShould Not be able to create a user.", tests.Success)
		}
	}
}

// TestCreateUserValidation tests the creation of a user that is not valid.
func TestCreateUserValidation(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	db, err := db.NewMGO(tests.Context, tests.TestSession)
	if err != nil {
		t.Fatalf("\t%s\tShould be able to get a Mongo session : %v", tests.Failed, err)
	}
	defer db.CloseMGO(tests.Context)

	var publicID string
	defer func() {
		if err := removeUser(db, publicID); err == nil {
			t.Fatalf("\t%s\tShould Not be able to remove the test user", tests.Failed)
		}
		t.Logf("\t%s\tShould Not be able to remove the test user.", tests.Success)
	}()

	t.Log("Given the need to make sure only valid users are created in the DB.")
	{
		t.Log("\tWhen using a test user.")
		{
			u, _ := auth.NewUser(auth.NUser{
				Status:   auth.StatusActive,
				FullName: "Test Kennedy",
				Email:    "bill@ardanlabs.com",
				Password: "_Password124",
			})

			u.Status = 0

			if err := auth.CreateUser(tests.Context, db, u); err == nil {
				t.Errorf("\t%s\tShould Not be able to create a user with invalid Status", tests.Failed)
			} else {
				t.Logf("\t%s\tShould Not be able to create a user with invalid Status.", tests.Success)
			}

			u, _ = auth.NewUser(auth.NUser{
				Status:   auth.StatusActive,
				FullName: "Test Kennedy",
				Email:    "bill@ardanlabs.com",
				Password: "_Password124",
			})

			u.FullName = "1234567"

			if err := auth.CreateUser(tests.Context, db, u); err == nil {
				t.Errorf("\t%s\tShould Not be able to create a user with invalid FullName", tests.Failed)
			} else {
				t.Logf("\t%s\tShould Not be able to create a user with invalid FullName.", tests.Success)
			}

			u, _ = auth.NewUser(auth.NUser{
				Status:   auth.StatusActive,
				FullName: "Test Kennedy",
				Email:    "bill@ardanlabs.com",
				Password: "_Password124",
			})

			u.Email = "bill"

			if err := auth.CreateUser(tests.Context, db, u); err == nil {
				t.Errorf("\t%s\tShould Not be able to create a user with invalid Email", tests.Failed)
			} else {
				t.Logf("\t%s\tShould Not be able to create a user with invalid Email.", tests.Success)
			}

			u, _ = auth.NewUser(auth.NUser{
				Status:   auth.StatusActive,
				FullName: "Test Kennedy",
				Email:    "bill@ardanlabs.com",
				Password: "_Password124",
			})

			u.Password = "1234567"

			if err := auth.CreateUser(tests.Context, db, u); err == nil {
				t.Errorf("\t%s\tShould Not be able to create a user with invalid Password", tests.Failed)
			} else {
				t.Logf("\t%s\tShould Not be able to create a user with invalid Password.", tests.Success)
			}
		}
	}
}

// TestUpdateUser tests we can update user information.
func TestUpdateUser(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	db, err := db.NewMGO(tests.Context, tests.TestSession)
	if err != nil {
		t.Fatalf("\t%s\tShould be able to get a Mongo session : %v", tests.Failed, err)
	}
	defer db.CloseMGO(tests.Context)

	var publicID string
	defer func() {
		if err := removeUser(db, publicID); err != nil {
			t.Fatalf("\t%s\tShould be able to remove the test user : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to remove the test user.", tests.Success)
	}()

	t.Log("Given the need to update a user.")
	{
		t.Log("\tWhen using an existing user.")
		{
			u1, err := auth.NewUser(auth.NUser{
				Status:   auth.StatusActive,
				FullName: "Test Kennedy",
				Email:    "bill@ardanlabs.com",
				Password: "_Password124",
			})
			if err != nil {
				t.Fatalf("\t%s\tShould be able to build a new user : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to build a new user.", tests.Success)

			if err := auth.CreateUser(tests.Context, db, u1); err != nil {
				t.Fatalf("\t%s\tShould be able to create a user : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a user.", tests.Success)

			// We need to do this so we can clean up after.
			publicID = u1.PublicID

			uu := auth.UpdUser{
				PublicID: publicID,
				Status:   auth.StatusActive,
				FullName: "Update Kennedy",
				Email:    "upt@ardanlabs.com",
			}

			if err := auth.UpdateUser(tests.Context, db, uu); err != nil {
				t.Fatalf("\t%s\tShould be able to update a user : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to update a user.", tests.Success)

			u2, err := auth.GetUserByPublicID(tests.Context, db, u1.PublicID, true)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve the user by PublicID : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to retrieve the user by PublicID.", tests.Success)

			// Remove the objectid to be able to compare the values.
			u2.ID = ""

			// Need to remove the nanoseconds to be able to compare the values.
			u1.DateModified = u1.DateModified.Add(-time.Duration(u1.DateModified.Nanosecond()))
			u1.DateCreated = u1.DateCreated.Add(-time.Duration(u1.DateCreated.Nanosecond()))
			u2.DateModified = u2.DateModified.Add(-time.Duration(u2.DateModified.Nanosecond()))
			u2.DateCreated = u2.DateCreated.Add(-time.Duration(u2.DateCreated.Nanosecond()))

			// Update the fields that changed
			u1.Status = u2.Status
			u1.FullName = u2.FullName
			u1.Email = u2.Email

			if !reflect.DeepEqual(*u1, *u2) {
				t.Logf("\t%+v", *u1)
				t.Logf("\t%+v", *u2)
				t.Errorf("\t%s\tShould be able to get back the same user with changes.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould be able to get back the same user with changes.", tests.Success)
			}
		}
	}
}

// TestUpdateUserValidation tests the update of a user that is not valid.
func TestUpdateUserValidation(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	db, err := db.NewMGO(tests.Context, tests.TestSession)
	if err != nil {
		t.Fatalf("\t%s\tShould be able to get a Mongo session : %v", tests.Failed, err)
	}
	defer db.CloseMGO(tests.Context)

	var publicID string
	defer func() {
		if err := removeUser(db, publicID); err == nil {
			t.Fatalf("\t%s\tShould Not be able to remove the test user", tests.Failed)
		}
		t.Logf("\t%s\tShould Not be able to remove the test user.", tests.Success)
	}()

	t.Log("Given the need to make sure only valid users are created in the DB.")
	{
		t.Log("\tWhen using a test user.")
		{
			uu := auth.UpdUser{
				PublicID: "asdasdasd",
				Status:   auth.StatusActive,
				FullName: "Test Kennedy",
				Email:    "bill@ardanlabs.com",
			}

			if err := auth.UpdateUser(tests.Context, db, uu); err == nil {
				t.Errorf("\t%s\tShould Not be able to update a user with invalid PublicID", tests.Failed)
			} else {
				t.Logf("\t%s\tShould Not be able to update a user with invalid PublicID.", tests.Success)
			}

			uu = auth.UpdUser{
				PublicID: "6dcda2da-92c3-11e5-8994-feff819cdc9f",
				Status:   0,
				FullName: "Test Kennedy",
				Email:    "bill@ardanlabs.com",
			}

			if err := auth.UpdateUser(tests.Context, db, uu); err == nil {
				t.Errorf("\t%s\tShould Not be able to update a user with invalid Status", tests.Failed)
			} else {
				t.Logf("\t%s\tShould Not be able to update a user with invalid Status.", tests.Success)
			}

			uu = auth.UpdUser{
				PublicID: "6dcda2da-92c3-11e5-8994-feff819cdc9f",
				Status:   auth.StatusActive,
				FullName: "1234567",
				Email:    "bill@ardanlabs.com",
			}

			if err := auth.UpdateUser(tests.Context, db, uu); err == nil {
				t.Errorf("\t%s\tShould Not be able to update a user with invalid FullName", tests.Failed)
			} else {
				t.Logf("\t%s\tShould Not be able to update a user with invalid FullName.", tests.Success)
			}

			uu = auth.UpdUser{
				PublicID: "6dcda2da-92c3-11e5-8994-feff819cdc9f",
				Status:   auth.StatusActive,
				FullName: "Test Kennedy",
				Email:    "ardanlabs.com",
			}

			if err := auth.UpdateUser(tests.Context, db, uu); err == nil {
				t.Errorf("\t%s\tShould Not be able to update a user with invalid Email", tests.Failed)
			} else {
				t.Logf("\t%s\tShould Not be able to update a user with invalid Email.", tests.Success)
			}

			uu = auth.UpdUser{
				PublicID: "6dcda2da-92c3-11e5-8994-feff819cdc9f",
				Status:   auth.StatusActive,
				FullName: "Test Kennedy",
				Email:    "bill@ardanlabs.com",
			}
		}
	}
}

// TestUpdateUserPassword tests we can update user password.
func TestUpdateUserPassword(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	db, err := db.NewMGO(tests.Context, tests.TestSession)
	if err != nil {
		t.Fatalf("\t%s\tShould be able to get a Mongo session : %v", tests.Failed, err)
	}
	defer db.CloseMGO(tests.Context)

	var publicID string
	defer func() {
		if err := removeUser(db, publicID); err != nil {
			t.Fatalf("\t%s\tShould be able to remove the test user : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to remove the test user.", tests.Success)
	}()

	t.Log("Given the need to update a user.")
	{
		t.Log("\tWhen using an existing user.")
		{
			u1, err := auth.NewUser(auth.NUser{
				Status:   auth.StatusActive,
				FullName: "Test Kennedy",
				Email:    "bill@ardanlabs.com",
				Password: "_Password124",
			})
			if err != nil {
				t.Fatalf("\t%s\tShould be able to build a new user : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to build a new user.", tests.Success)

			if err := auth.CreateUser(tests.Context, db, u1); err != nil {
				t.Fatalf("\t%s\tShould be able to create a user : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a user.", tests.Success)

			// We need to do this so we can clean up after.
			publicID = u1.PublicID

			webTok, err := auth.CreateWebToken(tests.Context, db, u1, 5*time.Second)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to create a web token : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a web token.", tests.Success)

			if err := auth.UpdateUserPassword(tests.Context, db, u1, "_Password567"); err != nil {
				t.Fatalf("\t%s\tShould be able to update a user : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to update a user.", tests.Success)

			if _, err := auth.ValidateWebToken(tests.Context, db, webTok); err == nil {
				t.Fatalf("\t%s\tShould Not be able to validate the org web token : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould Not be able to validate the new org token.", tests.Success)

			u2, err := auth.GetUserByPublicID(tests.Context, db, u1.PublicID, true)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve the user by PublicID : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to retrieve the user by PublicID.", tests.Success)

			webTok2, err := auth.CreateWebToken(tests.Context, db, u2, 5*time.Second)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to create a new web token : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a new web token.", tests.Success)

			if webTok == webTok2 {
				t.Fatalf("\t%s\tShould have different web tokens after the update.", tests.Failed)
			}
			t.Logf("\t%s\tShould have different web tokens after the update.", tests.Success)

			u3, err := auth.ValidateWebToken(tests.Context, db, webTok2)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to validate the new web token : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to validate the new web token.", tests.Success)

			if u1.PublicID != u3.PublicID {
				t.Log(u2.PublicID)
				t.Log(u3.PublicID)
				t.Fatalf("\t%s\tShould have the right user for the new token.", tests.Failed)
			}
			t.Logf("\t%s\tShould have the right user for the new token.", tests.Success)
		}
	}
}

// TestUpdateInvalidUserPassword tests we can't update user password.
func TestUpdateInvalidUserPassword(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	db, err := db.NewMGO(tests.Context, tests.TestSession)
	if err != nil {
		t.Fatalf("\t%s\tShould be able to get a Mongo session : %v", tests.Failed, err)
	}
	defer db.CloseMGO(tests.Context)

	t.Log("Given the need to validate an invalid update to a user.")
	{
		t.Log("\tWhen using an existing user.")
		{
			u1, err := auth.NewUser(auth.NUser{
				Status:   auth.StatusActive,
				FullName: "Test Kennedy",
				Email:    "bill@ardanlabs.com",
				Password: "_Password124",
			})
			if err != nil {
				t.Fatalf("\t%s\tShould be able to build a new user : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to build a new user.", tests.Success)

			if err := auth.UpdateUserPassword(tests.Context, db, u1, "_Pass"); err == nil {
				t.Errorf("\t%s\tShould Not be able to update a user with bad password.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould Not be able to update a user with bad password.", tests.Success)
			}

			u1.Status = auth.StatusDisabled

			if err := auth.UpdateUserPassword(tests.Context, db, u1, "_Password789"); err == nil {
				t.Errorf("\t%s\tShould Not be able to update a user with bad user value.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould Not be able to update a user with bad user value.", tests.Success)
			}
		}
	}
}

// TestDisableUser test the disabling of a user.
func TestDisableUser(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	db, err := db.NewMGO(tests.Context, tests.TestSession)
	if err != nil {
		t.Fatalf("\t%s\tShould be able to get a Mongo session : %v", tests.Failed, err)
	}
	defer db.CloseMGO(tests.Context)

	var publicID string
	defer func() {
		if err := removeUser(db, publicID); err != nil {
			t.Fatalf("\t%s\tShould be able to remove the test user : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to remove the test user.", tests.Success)
	}()

	t.Log("Given the need to update a user.")
	{
		t.Log("\tWhen using an existing user.")
		{
			u1, err := auth.NewUser(auth.NUser{
				Status:   auth.StatusActive,
				FullName: "Test Kennedy",
				Email:    "bill@ardanlabs.com",
				Password: "_Password124",
			})
			if err != nil {
				t.Fatalf("\t%s\tShould be able to build a new user : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to build a new user.", tests.Success)

			if err := auth.CreateUser(tests.Context, db, u1); err != nil {
				t.Fatalf("\t%s\tShould be able to create a user : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a user.", tests.Success)

			// We need to do this so we can clean up after.
			publicID = u1.PublicID

			u2, err := auth.GetUserByPublicID(tests.Context, db, u1.PublicID, true)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve the user by PublicID : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to retrieve the user by PublicID.", tests.Success)

			if err := auth.UpdateUserStatus(tests.Context, db, u2.PublicID, auth.StatusDisabled); err != nil {
				t.Fatalf("\t%s\tShould be able to disable the user : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to disable the user.", tests.Success)

			if _, err := auth.GetUserByPublicID(tests.Context, db, u1.PublicID, true); err == nil {
				t.Fatalf("\t%s\tShould Not be able to retrieve the user by PublicID.", tests.Failed)
			}
			t.Logf("\t%s\tShould Not be able to retrieve the user by PublicID.", tests.Success)
		}
	}
}

// TestCreateWebToken tests create a web token and a pairing session.
func TestCreateWebToken(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	db, err := db.NewMGO(tests.Context, tests.TestSession)
	if err != nil {
		t.Fatalf("\t%s\tShould be able to get a Mongo session : %v", tests.Failed, err)
	}
	defer db.CloseMGO(tests.Context)

	var publicID string
	defer func() {
		if err := removeUser(db, publicID); err != nil {
			t.Fatalf("\t%s\tShould be able to remove the test user : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to remove the test user.", tests.Success)
	}()

	t.Log("Given the need to create a web token.")
	{
		t.Log("\tWhen using a new user.")
		{
			u1, err := auth.NewUser(auth.NUser{
				Status:   auth.StatusActive,
				FullName: "Test Kennedy",
				Email:    "bill@ardanlabs.com",
				Password: "_Password124",
			})
			if err != nil {
				t.Fatalf("\t%s\tShould be able to build a new user : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to build a new user.", tests.Success)

			if err := auth.CreateUser(tests.Context, db, u1); err != nil {
				t.Fatalf("\t%s\tShould be able to create a user : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a user.", tests.Success)

			// We need to do this so we can clean up after.
			publicID = u1.PublicID

			webTok, err := auth.CreateWebToken(tests.Context, db, u1, time.Second)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to create a web token : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a web token.", tests.Success)

			sId, _, err := auth.DecodeWebToken(tests.Context, webTok)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to decode the web token : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to decode the web token.", tests.Success)

			s2, err := session.GetBySessionID(tests.Context, db, sId)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve the session : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to retrieve the session.", tests.Success)

			u2, err := auth.GetUserByPublicID(tests.Context, db, u1.PublicID, true)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve the user by PublicID : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to retrieve the user by PublicID.", tests.Success)

			if u2.PublicID != s2.PublicID {
				t.Fatalf("\t%s\tShould have the right session for user.", tests.Failed)
				t.Log(u2.PublicID)
				t.Log(s2.PublicID)
			}
			t.Logf("\t%s\tShould have the right session for user.", tests.Success)

			webTok2, err := u2.WebToken(sId)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to create a new web token : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a web new token.", tests.Success)

			if webTok != webTok2 {
				t.Log(webTok)
				t.Log(webTok2)
				t.Fatalf("\t%s\tShould be able to create the same web token.", tests.Failed)
			}
			t.Logf("\t%s\tShould be able to create the same web token.", tests.Success)

			u3, err := auth.ValidateWebToken(tests.Context, db, webTok2)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to validate the new web token : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to validate the new web token.", tests.Success)

			if u1.PublicID != u3.PublicID {
				t.Log(u1.PublicID)
				t.Log(u3.PublicID)
				t.Fatalf("\t%s\tShould have the right user for the token.", tests.Failed)
			}
			t.Logf("\t%s\tShould have the right user for the token.", tests.Success)

			webTok3, err := auth.GetUserWebToken(tests.Context, db, u2.PublicID)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to get the web token : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to get the web token.", tests.Success)

			if webTok3 != webTok2 {
				t.Log(webTok3)
				t.Log(webTok2)
				t.Fatalf("\t%s\tShould match existing tokens.", tests.Failed)
			}
			t.Logf("\t%s\tShould match existing tokens.", tests.Success)
		}
	}
}

// TestExpiredWebToken tests create a web token and tests when it expires.
func TestExpiredWebToken(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	db, err := db.NewMGO(tests.Context, tests.TestSession)
	if err != nil {
		t.Fatalf("\t%s\tShould be able to get a Mongo session : %v", tests.Failed, err)
	}
	defer db.CloseMGO(tests.Context)

	var publicID string
	defer func() {
		if err := removeUser(db, publicID); err != nil {
			t.Fatalf("\t%s\tShould be able to remove the test user : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to remove the test user.", tests.Success)
	}()

	t.Log("Given the need to validate web tokens expire.")
	{
		t.Log("\tWhen using a new user.")
		{
			u1, err := auth.NewUser(auth.NUser{
				Status:   auth.StatusActive,
				FullName: "Test Kennedy",
				Email:    "bill@ardanlabs.com",
				Password: "_Password124",
			})
			if err != nil {
				t.Fatalf("\t%s\tShould be able to build a new user : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to build a new user.", tests.Success)

			if err := auth.CreateUser(tests.Context, db, u1); err != nil {
				t.Fatalf("\t%s\tShould be able to create a user : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a user.", tests.Success)

			// We need to do this so we can clean up after.
			publicID = u1.PublicID

			webTok, err := auth.CreateWebToken(tests.Context, db, u1, 1*time.Millisecond)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to create a web token : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a web token.", tests.Success)

			if _, err := auth.ValidateWebToken(tests.Context, db, webTok); err == nil {
				t.Fatalf("\t%s\tShould Not be able to validate the web token : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould Not be able to validate the web token.", tests.Success)
		}
	}
}

// TestInvalidWebTokens tests create an invalid web token and tests it fails.
func TestInvalidWebTokens(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	db, err := db.NewMGO(tests.Context, tests.TestSession)
	if err != nil {
		t.Fatalf("\t%s\tShould be able to get a Mongo session : %v", tests.Failed, err)
	}
	defer db.CloseMGO(tests.Context)

	tokens := []string{
		"",
		"6dcda2da-92c3-11e5-8994-feff819cdc9f",
		"OGY4OGI3YWQtZjc5Ny00ODI1LWI0MmUtMjIwZTY5ZDQxYjMzOmFKT2U1b0pFZlZ4cWUrR0JONEl0WlhmQTY0K3JsN2VGcmM2MVNQMkV1WVE9",
	}

	t.Log("Given the need to validate bad web tokens don't validate.")
	{
		for _, token := range tokens {
			t.Logf("\tWhen using token [%s]", token)
			{
				if _, err := auth.ValidateWebToken(tests.Context, db, token); err == nil {
					t.Errorf("\t%s\tShould Not be able to validate the web token : %v", tests.Failed, err)
				} else {
					t.Logf("\t%s\tShould Not be able to validate the web token.", tests.Success)
				}
			}
		}
	}
}

// TestInvalidWebTokenUpdateEmail tests a token becomes invalid after an update.
func TestInvalidWebTokenUpdateEmail(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	db, err := db.NewMGO(tests.Context, tests.TestSession)
	if err != nil {
		t.Fatalf("\t%s\tShould be able to get a Mongo session : %v", tests.Failed, err)
	}
	defer db.CloseMGO(tests.Context)

	var publicID string
	defer func() {
		if err := removeUser(db, publicID); err != nil {
			t.Fatalf("\t%s\tShould be able to remove the test user : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to remove the test user.", tests.Success)
	}()

	t.Log("Given the need to validate web tokens don't work after user update.")
	{
		t.Log("\tWhen using a new user.")
		{
			u1, err := auth.NewUser(auth.NUser{
				Status:   auth.StatusActive,
				FullName: "Test Kennedy",
				Email:    "bill@ardanlabs.com",
				Password: "_Password124",
			})
			if err != nil {
				t.Fatalf("\t%s\tShould be able to build a new user : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to build a new user.", tests.Success)

			if err := auth.CreateUser(tests.Context, db, u1); err != nil {
				t.Fatalf("\t%s\tShould be able to create a user : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a user.", tests.Success)

			// We need to do this so we can clean up after.
			publicID = u1.PublicID

			webTok, err := auth.CreateWebToken(tests.Context, db, u1, 5*time.Second)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to create a web token : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a web token.", tests.Success)

			if _, err := auth.ValidateWebToken(tests.Context, db, webTok); err != nil {
				t.Fatalf("\t%s\tShould be able to validate the web token : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to validate the web token.", tests.Success)

			uu := auth.UpdUser{
				PublicID: publicID,
				Status:   auth.StatusActive,
				FullName: "Update Kennedy",
				Email:    "change@ardanlabs.com",
			}

			if err := auth.UpdateUser(tests.Context, db, uu); err != nil {
				t.Fatalf("\t%s\tShould be able to update a user : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to update a user.", tests.Success)

			if _, err := auth.ValidateWebToken(tests.Context, db, webTok); err == nil {
				t.Fatalf("\t%s\tShould Not be able to validate the org web token.", tests.Failed)
			}
			t.Logf("\t%s\tShould Not be able to validate the org web token.", tests.Success)
		}
	}
}

// TestLoginUser validates a user can login and not after changes.
func TestLoginUser(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	db, err := db.NewMGO(tests.Context, tests.TestSession)
	if err != nil {
		t.Fatalf("\t%s\tShould be able to get a Mongo session : %v", tests.Failed, err)
	}
	defer db.CloseMGO(tests.Context)

	var publicID string
	defer func() {
		if err := removeUser(db, publicID); err != nil {
			t.Fatalf("\t%s\tShould be able to remove the test user : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to remove the test user.", tests.Success)
	}()

	t.Log("Given the need to test user login.")
	{
		t.Log("\tWhen using a new user")
		{
			u1, err := auth.NewUser(auth.NUser{
				Status:   auth.StatusActive,
				FullName: "Test Kennedy",
				Email:    "bill@ardanlabs.com",
				Password: "_Password124",
			})
			if err != nil {
				t.Fatalf("\t%s\tShould be able to build a new user : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to build a new user.", tests.Success)

			if err := auth.CreateUser(tests.Context, db, u1); err != nil {
				t.Fatalf("\t%s\tShould be able to create a user : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a user.", tests.Success)

			// We need to do this so we can clean up after.
			publicID = u1.PublicID

			if _, err := auth.LoginUser(tests.Context, db, u1.Email, "_Password124"); err != nil {
				t.Errorf("\t%s\tShould be able to login the user : %v", tests.Failed, err)
			} else {
				t.Logf("\t%s\tShould be able to login the user.", tests.Success)
			}

			if err := auth.UpdateUserPassword(tests.Context, db, u1, "password890"); err != nil {
				t.Fatalf("\t%s\tShould be able to update the user password : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to update the user password.", tests.Success)

			if _, err := auth.LoginUser(tests.Context, db, u1.Email, "_Password124"); err == nil {
				t.Errorf("\t%s\tShould Not be able to login the user.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould Not be able to login the user.", tests.Success)
			}
		}
	}
}

// TestNoSession tests when a nil session is used.
func TestNoSession(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	t.Log("Given the need to test calls with a bad session.")
	{
		t.Log("\tWhen using a nil session")
		{
			u1, err := auth.NewUser(auth.NUser{
				Status:   auth.StatusActive,
				FullName: "Test Kennedy",
				Email:    "bill@ardanlabs.com",
				Password: "_Password124",
			})
			if err != nil {
				t.Fatalf("\t%s\tShould be able to build a new user : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to build a new user.", tests.Success)

			if err := auth.CreateUser(tests.Context, nil, u1); err == nil {
				t.Errorf("\t%s\tShould Not be able to create a user.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould Not be able to create a user.", tests.Success)
			}

			if _, err := auth.CreateWebToken(tests.Context, nil, u1, time.Second); err == nil {
				t.Errorf("\t%s\tShould Not be able to create a web token.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould Not be able to create a web token.", tests.Success)
			}

			webTok := "OGY4OGI3YWQtZjc5Ny00ODI1LWI0MmUtMjIwZTY5ZDQxYjMzOmFKT2U1b0pFZlZ4cWUrR0JONEl0WlhmQTY0K3JsN2VGcmM2MVNQMkV1WVE9"

			if _, err := auth.ValidateWebToken(tests.Context, nil, webTok); err == nil {
				t.Errorf("\t%s\tShould Not be able to validate a web token.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould Not be able to validate a web token.", tests.Success)
			}

			if _, err := auth.GetUserByPublicID(tests.Context, nil, "6dcda2da-92c3-11e5-8994-feff819cdc9f", true); err == nil {
				t.Errorf("\t%s\tShould Not be able to get a user by PublicID.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould Not be able to get a user by PublicID.", tests.Success)
			}

			if _, err := auth.GetUserByEmail(tests.Context, nil, "bill@ardanlabs.com", true); err == nil {
				t.Errorf("\t%s\tShould Not be able to get a user by Email.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould Not be able to get a user by Email.", tests.Success)
			}

			uu := auth.UpdUser{
				PublicID: "6dcda2da-92c3-11e5-8994-feff819cdc9f",
				Status:   auth.StatusActive,
				FullName: "Update Kennedy",
				Email:    "upt@ardanlabs.com",
			}

			if err := auth.UpdateUser(tests.Context, nil, uu); err == nil {
				t.Errorf("\t%s\tShould Not be able to update a user.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould Not be able to update a user.", tests.Success)
			}

			if err := auth.UpdateUserPassword(tests.Context, nil, u1, "password890"); err == nil {
				t.Errorf("\t%s\tShould Not be able to update a user pasword.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould Not be able to update a user password.", tests.Success)
			}

			if err := auth.UpdateUserStatus(tests.Context, nil, "6dcda2da-92c3-11e5-8994-feff819cdc9f", auth.StatusDisabled); err == nil {
				t.Errorf("\t%s\tShould Not be able to disable a user.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould Not be able to disable a user.", tests.Success)
			}

			if _, err := auth.LoginUser(tests.Context, nil, "bill@email.com", "_pass"); err == nil {
				t.Errorf("\t%s\tShould Not be able to login a user.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould Not be able to login a user.", tests.Success)
			}

			if _, err := auth.GetUserWebToken(tests.Context, nil, "6dcda2da-92c3-11e5-8994-feff819cdc9f"); err == nil {
				t.Errorf("\t%s\tShould Not be able to get user web token.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould Not be able to get user web token.", tests.Success)
			}
		}
	}
}

//==============================================================================

func ensureIndexes() {
	db, err := db.NewMGO(tests.Context, tests.TestSession)
	if err != nil {
		fmt.Printf("Should be able to get a Mongo session : %v", err)
		os.Exit(1)
	}
	defer db.CloseMGO(tests.Context)

	index := mgo.Index{
		Key:    []string{"public_id"},
		Unique: true,
	}

	col, err := db.CollectionMGO(tests.Context, auth.Collection)
	if err != nil {
		fmt.Printf("Should be able to get a Mongo session : %v", err)
		os.Exit(1)
	}

	col.EnsureIndex(index)
}

// removeUser is used to clear out all the test user from the collection.
func removeUser(db *db.DB, publicID string) error {
	f := func(c *mgo.Collection) error {
		q := bson.M{"public_id": publicID}
		return c.Remove(q)
	}

	if err := db.ExecuteMGO(tests.Context, auth.Collection, f); err != nil {
		return err
	}

	f = func(c *mgo.Collection) error {
		q := bson.M{"public_id": publicID}
		_, err := c.RemoveAll(q)
		return err
	}

	if err := db.ExecuteMGO(tests.Context, session.Collection, f); err != nil {
		return err
	}

	return nil
}
