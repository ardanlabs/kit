package db

import (
	"errors"
	"sync"

	"github.com/ardanlabs/kit/db/mongo"

	"gopkg.in/mgo.v2"
)

// mgoDB maintains a master session for a given database.
type mgoDB struct {
	dbName string
	ses    *mgo.Session
}

// masterMGO manages a set of different MongoDB master sessions.
var masterMGO = struct {
	sync.RWMutex
	ses map[string]mgoDB
}{
	ses: make(map[string]mgoDB),
}

// RegMasterSession adds a new master session to the set.
func RegMasterSession(context interface{}, name string, cfg mongo.Config) error {
	masterMGO.Lock()
	defer masterMGO.Unlock()

	if _, exists := masterMGO.ses[name]; exists {
		return errors.New("Master session already exists")
	}

	ses, err := mongo.New(cfg)
	if err != nil {
		return err
	}

	masterMGO.ses[name] = mgoDB{cfg.DB, ses}

	return nil
}
