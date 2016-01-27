package cmddb

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/ardanlabs/kit/db"

	"github.com/spf13/cobra"
	"gopkg.in/mgo.v2"
)

var createLong = `Creates a new database with collections and indexes.

Example:

	db create -f ./scripts/database.json
`

// create contains the state for this command.
var create struct {
	file string
}

//==============================================================================

// addCreate handles the creation of users.
func addCreate() {
	cmd := &cobra.Command{
		Use:   "create [-f file]",
		Short: "Creates a new database from a script file",
		Long:  createLong,
		Run:   runCreate,
	}

	cmd.Flags().StringVarP(&create.file, "file", "f", "", "file path of script json file")

	dbCmd.AddCommand(cmd)
}

//==============================================================================

var (
	// ErrCollectionExists is return when a collection to be
	// created already exists.
	ErrCollectionExists = errors.New("Collection already exists.")
)

// DBMeta is the container for all db objects.
type DBMeta struct {
	Cols []Collection `json:"collections"`
}

// Collection is the container for a db collection definition.
type Collection struct {
	Name    string  `json:"name"`
	Indexes []Index `json:"indexes"`
}

// Index is the container for an index definition.
type Index struct {
	Name     string  `json:"name"`
	IsUnique bool    `json:"unique"`
	Fields   []Field `json:"fields"`
}

// Field is the container for a field definition.
type Field struct {
	Name      string `json:"name"`
	Type      int    `json:"type"`
	OtherType string `json:"other"`
}

// runCreate is the code that implements the create command.
func runCreate(cmd *cobra.Command, args []string) {
	dbMeta, err := retrieveDatabaseMetadata(create.file)
	if err != nil {
		dbCmd.Printf("Error reading collections : %s : ERROR : %v\n", create.file, err)
		return
	}

	for _, col := range dbMeta.Cols {
		cmd.Println("Creating collection", col.Name)
		if err := createCollection(conn, dbMeta, &col, true); err != nil && err != ErrCollectionExists {
			cmd.Println("ERROR:", err)
			return
		}
	}
}

// retrieveDatabaseMetadata reads the specified file and returns the database
// metadata for creating and updating the database.
func retrieveDatabaseMetadata(fileName string) (*DBMeta, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	var dbMeta DBMeta
	if err := json.NewDecoder(f).Decode(&dbMeta); err != nil {
		return nil, err
	}

	return &dbMeta, nil
}

// createCollection creates a collection in the new database.
func createCollection(db *db.DB, dbMeta *DBMeta, col *Collection, dropIdxs bool) error {
	mCol, err := db.CollectionMGO("", col.Name)
	if err != nil {
		return err
	}

	if err := mCol.Create(new(mgo.CollectionInfo)); err != nil {
		return err
	}

	if err := createIndexes(mCol, col, dropIdxs); err != nil {
		return err
	}

	return nil
}

// createIndexes creates a required indexes in the new database.
func createIndexes(mCol *mgo.Collection, col *Collection, dropIdxs bool) error {
	if dropIdxs == true {
		idxs, err := mCol.Indexes()
		if err != nil {
			return err
		}

		for _, idx := range idxs {
			mCol.DropIndex(idx.Name)
		}
	}

	for _, idx := range col.Indexes {
		newIdx := mgo.Index{
			Key:    parseFields(idx.Fields),
			Unique: idx.IsUnique,
			Name:   idx.Name,
		}

		if err := mCol.EnsureIndex(newIdx); err != nil {
			return err
		}
	}

	return nil
}

// parseFields formats the field array to determine the fields in the index.
func parseFields(idxFields []Field) []string {
	var flds []string

	for _, fld := range idxFields {
		switch fld.Type {
		case -1:
			flds = append(flds, "-"+fld.Name)

		case 0:
			f := fmt.Sprintf("$%s:%s", fld.OtherType, fld.Name)
			flds = append(flds, f)

		default:
			flds = append(flds, fld.Name)
		}
	}

	return flds
}
