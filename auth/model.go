package auth

import (
	"encoding/base64"
	"errors"
	"strings"
	"time"

	"github.com/ardanlabs/kit/auth/crypto"

	"github.com/pborman/uuid"
	"gopkg.in/bluesuncorp/validator.v8"
	"gopkg.in/mgo.v2/bson"
)

// Set of user status codes.
const (
	StatusUnknown = iota
	StatusActive
	StatusDisabled
)

//==============================================================================

// validate is used to perform model field validation.
var validate *validator.Validate

func init() {
	validate = validator.New(&validator.Config{TagName: "validate"})
}

//==============================================================================

// NUser is provided to create a new user value for use.
type NUser struct {
	Status   int    `bson:"status" json:"status" validate:"required,ne=0"`
	FullName string `bson:"full_name" json:"full_name" validate:"required,min=8"`
	Email    string `bson:"email" json:"email" validate:"required,max=100,email"`
	Password string `bson:"password" json:"-" validate:"required,min=8"`
}

// Validate performs validation on a NUser value before it is processed.
func (nu *NUser) Validate() error {
	if err := validate.Struct(nu); err != nil {
		return err
	}

	return nil
}

//==============================================================================

// User model denotes a user entity for a tenant.
type User struct {
	ID           bson.ObjectId `bson:"_id,omitempty" json:"-"`
	PublicID     string        `bson:"public_id" json:"public_id" validate:"required,uuid"`
	PrivateID    string        `bson:"private_id" json:"-" validate:"required,uuid"`
	Status       int           `bson:"status" json:"status" validate:"required,ne=0"`
	FullName     string        `bson:"full_name" json:"full_name" validate:"required,min=8"`
	Email        string        `bson:"email" json:"email" validate:"required,max=100,email"`
	Password     string        `bson:"password" json:"-" validate:"required,min=55"`
	IsDeleted    bool          `bson:"is_deleted" json:"-"`
	DateModified time.Time     `bson:"date_modified" json:"-"`
	DateCreated  time.Time     `bson:"date_created" json:"-"`
}

// NewUser creates a new user from a NewUser value.
func NewUser(nu NUser) (*User, error) {
	if err := nu.Validate(); err != nil {
		return nil, err
	}

	u := User{
		PublicID:     uuid.New(),
		PrivateID:    uuid.New(),
		Status:       StatusActive,
		FullName:     nu.FullName,
		Email:        strings.ToLower(nu.Email),
		DateModified: time.Now(),
		DateCreated:  time.Now(),
		IsDeleted:    false,
	}

	var err error
	if u.Password, err = crypto.BcryptPassword(u.PrivateID + nu.Password); err != nil {
		return nil, err
	}

	return &u, nil
}

// Validate performs validation on a CrtUser value before it is processed.
func (u *User) Validate() error {
	if err := validate.Struct(u); err != nil {
		return err
	}

	return nil
}

// Pwd implements the secure entity interface.
func (u *User) Pwd() ([]byte, error) {
	if u.Password == "" {
		return nil, errors.New("User password is blank")
	}

	return []byte(u.Password), nil
}

// Salt implements the secure entity interface.
func (u *User) Salt() ([]byte, error) {
	if u.PublicID == "" || u.PrivateID == "" || u.Email == "" {
		return nil, errors.New("Unable to generate user token, missing data")
	}

	s := u.PublicID + ":" + u.PrivateID + ":" + u.Email

	return []byte(s), nil
}

// WebToken returns a token ready for web use.
func (u *User) WebToken(sessionID string) (string, error) {
	t, err := crypto.GenerateToken(u)
	if err != nil {
		return "", err
	}

	token := base64.StdEncoding.EncodeToString([]byte(sessionID + ":" + base64.StdEncoding.EncodeToString(t)))
	return token, nil
}

// IsPasswordValid compares the user provided password with what is in the db.
func (u *User) IsPasswordValid(password string) bool {
	if u.Password == "" {
		return false
	}

	// Hashed Password comes first, then the plain text version.
	if err := crypto.CompareBcryptHashPassword([]byte(u.Password), []byte(u.PrivateID+password)); err != nil {
		return false
	}

	return true
}

//==============================================================================

// UpdUser is provided to update an existing user in the system.
type UpdUser struct {
	PublicID string `bson:"public_id" json:"public_id" validate:"required,uuid"`
	Status   int    `bson:"status" json:"status" validate:"required,ne=0"`
	FullName string `bson:"full_name" json:"full_name" validate:"required,min=8"`
	Email    string `bson:"email" json:"email" validate:"required,max=100,email"`
}

// Validate performs validation on a NewUser value before it is processed.
func (uu *UpdUser) Validate() error {
	if err := validate.Struct(uu); err != nil {
		return err
	}

	return nil
}
