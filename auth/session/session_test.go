package session_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/ardanlabs/kit/auth/session"
	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/db/mongo"
	"github.com/ardanlabs/kit/tests"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const publicID = "6dcda2da-92c3-11e5-8994-feff819cdc9f"

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

func ensureIndexes() {
	db, err := db.NewMGO(tests.Context, tests.TestSession)
	if err != nil {
		fmt.Printf("Should be able to get a Mongo session : %v", err)
		os.Exit(1)
	}
	defer db.CloseMGO(tests.Context)

	index := mgo.Index{
		Key:    []string{"public_id"},
		Unique: false,
	}

	col, err := db.CollectionMGO(tests.Context, session.Collection)
	if err != nil {
		fmt.Printf("Should be able to get a Mongo session : %v", err)
		os.Exit(1)
	}

	col.EnsureIndex(index)
}

// removeSessions is used to clear out all the test sessions that are
// created from tests.
func removeSessions(db *db.DB) error {
	f := func(c *mgo.Collection) error {
		q := bson.M{"public_id": publicID}
		_, err := c.RemoveAll(q)
		return err
	}

	if err := db.ExecuteMGO(tests.Context, session.Collection, f); err != nil {
		return err
	}

	return nil
}

// TestIsExpired tests that we can properly identify when a
// session is expired or not.
func TestIsExpired(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	t.Log("Given the need to validate a session has expired.")
	{
		s := session.Session{
			DateExpires: time.Now().Add(-time.Hour),
		}

		t.Log("\tWhen using an expired session.")
		{
			if !s.IsExpired(tests.Context) {
				t.Fatalf("\t%s\tShould be expired.", tests.Failed)
			}
			t.Logf("\t%s\tShould be expired", tests.Success)
		}

		s = session.Session{
			DateExpires: time.Now().Add(time.Hour),
		}

		t.Log("\tWhen using an valid session")
		{
			if s.IsExpired(tests.Context) {
				t.Fatalf("\t%s\tShould Not be expired.", tests.Failed)
			}
			t.Logf("\t%s\tShould Not be expired", tests.Success)
		}
	}
}

// TestCreate tests the creation of sessions.
func TestCreate(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	db, err := db.NewMGO(tests.Context, tests.TestSession)
	if err != nil {
		t.Fatalf("\t%s\tShould be able to get a Mongo session : %v", tests.Failed, err)
	}
	defer db.CloseMGO(tests.Context)

	defer func() {
		if err := removeSessions(db); err != nil {
			t.Errorf("\t%s\tShould be able to remove all sessions : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to remove all sessions.", tests.Success)
	}()

	t.Log("Given the need to create sessions in the DB.")
	{
		t.Logf("\tWhen using PublicID %s", publicID)
		{
			if err := removeSessions(db); err != nil {
				t.Fatalf("\t%s\tShould be able to remove all sessions : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to remove all sessions.", tests.Success)

			s1, err := session.Create(tests.Context, db, publicID, 10*time.Second)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to create a session : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a session.", tests.Success)

			s2, err := session.GetBySessionID(tests.Context, db, s1.SessionID)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve the session : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to retrieve the session.", tests.Success)

			if s1.SessionID != s2.SessionID {
				t.Fatalf("\t%s\tShould be able to get back the same session.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould be able to get back the same session.", tests.Success)
			}

			if s1.PublicID != s2.PublicID {
				t.Fatalf("\t%s\tShould be able to get back the same user.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould be able to get back the same user.", tests.Success)
			}
		}
	}
}

// TestGetLatest tests the retrieval of the latest session.
func TestGetLatest(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	db, err := db.NewMGO(tests.Context, tests.TestSession)
	if err != nil {
		t.Fatalf("\t%s\tShould be able to get a Mongo session : %v", tests.Failed, err)
	}
	defer db.CloseMGO(tests.Context)

	defer func() {
		if err := removeSessions(db); err != nil {
			t.Errorf("\t%s\tShould be able to remove all sessions : %v", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to remove all sessions.", tests.Success)
	}()

	t.Log("Given the need to get the latest sessions in the DB.")
	{
		t.Logf("\tWhen using PublicID %s", publicID)
		{
			if err := removeSessions(db); err != nil {
				t.Fatalf("\t%s\tShould be able to remove all sessions : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to remove all sessions.", tests.Success)

			if _, err := session.Create(tests.Context, db, publicID, 10*time.Second); err != nil {
				t.Fatalf("\t%s\tShould be able to create a session : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a session.", tests.Success)

			time.Sleep(time.Second)

			s2, err := session.Create(tests.Context, db, publicID, 10*time.Second)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to create another session : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create another session.", tests.Success)

			s3, err := session.GetByLatest(tests.Context, db, publicID)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve the latest session : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to retrieve the latest session.", tests.Success)

			if s2.SessionID != s3.SessionID {
				t.Errorf("\t%s\tShould be able to get back the latest session.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould be able to get back the latest session.", tests.Success)
			}
		}
	}
}

// TestGetNotFound tests when a session is not found.
func TestGetNotFound(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	db, err := db.NewMGO(tests.Context, tests.TestSession)
	if err != nil {
		t.Fatalf("\t%s\tShould be able to get a Mongo session : %v", tests.Failed, err)
	}
	defer db.CloseMGO(tests.Context)

	t.Log("Given the need to test finding a session and it is not found.")
	{
		t.Logf("\tWhen using SessionID %s", "NOT EXISTS")
		{
			if _, err := session.GetBySessionID(tests.Context, db, "NOT EXISTS"); err == nil {
				t.Fatalf("\t%s\tShould Not be able to retrieve the session.", tests.Failed)
			}
			t.Logf("\t%s\tShould Not be able to retrieve the session.", tests.Success)
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
			if _, err := session.Create(tests.Context, nil, publicID, 10*time.Second); err == nil {
				t.Errorf("\t%s\tShould Not be able to create a session.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould Not be able to create a session.", tests.Success)
			}

			if _, err := session.GetBySessionID(tests.Context, nil, "NOT EXISTS"); err == nil {
				t.Errorf("\t%s\tShould Not be able to retrieve the session.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould Not be able to retrieve the session.", tests.Success)
			}

			if _, err := session.GetByLatest(tests.Context, nil, publicID); err == nil {
				t.Errorf("\t%s\tShould Not be able to retrieve the session.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould Not be able to retrieve the session.", tests.Success)
			}
		}
	}
}
