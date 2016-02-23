package auth

import (
	"encoding/base64"
	"errors"
	"strings"
	"time"

	"github.com/ardanlabs/kit/auth/crypto"
	"github.com/ardanlabs/kit/auth/session"
	"github.com/ardanlabs/kit/db"
	"github.com/ardanlabs/kit/db/mongo"
	"github.com/ardanlabs/kit/log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Collection contains the name of the auth_users collection.
const Collection = "auth_users"

//==============================================================================

// CreateUser adds a new user to the database.
func CreateUser(context interface{}, db *db.DB, u *User) error {
	log.Dev(context, "CreateUser", "Started : PublicID[%s]", u.PublicID)

	if err := u.Validate(); err != nil {
		log.Error(context, "CreateUser", err, "Completed")
		return err
	}

	f := func(c *mgo.Collection) error {
		log.Dev(context, "CreateUser", "MGO : db.%s.insert(%s)", c.Name, mongo.Query(&u))
		return c.Insert(u)
	}

	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		log.Error(context, "CreateUser", err, "Completed")
		return err
	}

	log.Dev(context, "CreateUser", "Completed")
	return nil
}

// CreateWebToken return a token and session that can be used to authenticate a user.
func CreateWebToken(context interface{}, db *db.DB, u *User, expires time.Duration) (string, error) {
	log.Dev(context, "CreateWebToken", "Started : PublicID[%s]", u.PublicID)

	// Do we have a valid session right now?
	s, err := session.GetByLatest(context, db, u.PublicID)
	if err != nil && err != mgo.ErrNotFound {
		log.Error(context, "CreateUser", err, "Completed")
		return "", err
	}

	// If we don't have one or it has been expired create
	// a new one.
	if err == mgo.ErrNotFound || s.IsExpired(context) {
		if s, err = session.Create(context, db, u.PublicID, expires); err != nil {
			log.Error(context, "CreateUser", err, "Completed")
			return "", err
		}
	}

	// Set the return arguments though we will explicitly
	// return them. Don't want any confusion.
	token, err := u.WebToken(s.SessionID)
	if err != nil {
		log.Error(context, "CreateUser", err, "Completed")
		return "", err
	}

	log.Dev(context, "CreateWebToken", "Completed : WebToken[%s]", token)
	return token, nil
}

//==============================================================================

// DecodeWebToken breaks a web token into its parts.
func DecodeWebToken(context interface{}, webToken string) (sessionID string, token string, err error) {
	log.Dev(context, "DecodeWebToken", "Started : WebToken[%s]", webToken)

	// Decode the web token to break it into its parts.
	data, err := base64.StdEncoding.DecodeString(webToken)
	if err != nil {
		log.Error(context, "DecodeWebToken", err, "Completed")
		return "", "", err
	}

	// Split the web token.
	str := strings.Split(string(data), ":")
	if len(str) != 2 {
		err := errors.New("Invalid token")
		log.Error(context, "DecodeWebToken", err, "Completed")
		return "", "", err
	}

	// Pull out the session and token.
	sessionID = str[0]
	token = str[1]

	log.Dev(context, "DecodeWebToken", "Completed : SessionID[%s] Token[%s]", sessionID, token)
	return sessionID, token, nil
}

// ValidateWebToken accepts a web token and validates its credibility. Returns
// a User value is the token is valid.
func ValidateWebToken(context interface{}, db *db.DB, webToken string) (*User, error) {
	log.Dev(context, "ValidateWebToken", "Started : WebToken[%s]", webToken)

	// Extract the sessionID and token from the web token.
	sessionID, token, err := DecodeWebToken(context, webToken)
	if err != nil {
		log.Error(context, "ValidateWebToken", err, "Completed")
		return nil, err
	}

	// Find the session in the database.
	s, err := session.GetBySessionID(context, db, sessionID)
	if err != nil {
		log.Error(context, "ValidateWebToken", err, "Completed")
		return nil, err
	}

	// Validate the session has not expired.
	if s.IsExpired(context) {
		err := errors.New("Expired token")
		log.Error(context, "ValidateWebToken", err, "Completed")
		return nil, err
	}

	// Pull the user for this session.
	u, err := GetUserByPublicID(context, db, s.PublicID, true)
	if err != nil {
		log.Error(context, "ValidateWebToken", err, "Completed")
		return nil, err
	}

	// Validate the token against this user.
	if err := crypto.IsTokenValid(u, token); err != nil {
		log.Error(context, "ValidateWebToken", err, "Completed")
		return nil, err
	}

	log.Dev(context, "ValidateWebToken", "Completed : PublicID[%s]", u.PublicID)
	return u, nil
}

//==============================================================================

// GetUserByPublicID retrieves a user record by using the provided PublicID.
func GetUserByPublicID(context interface{}, db *db.DB, publicID string, activeOnly bool) (*User, error) {
	log.Dev(context, "GetUserByPublicID", "Started : PID[%s]", publicID)

	var user User
	f := func(c *mgo.Collection) error {
		var q bson.M
		if activeOnly {
			q = bson.M{"public_id": publicID, "status": StatusActive}
		} else {
			q = bson.M{"public_id": publicID}
		}
		log.Dev(context, "GetUserByPublicID", "MGO : db.%s.findOne(%s)", c.Name, mongo.Query(q))
		return c.Find(q).One(&user)
	}

	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		log.Error(context, "GetUserByPublicID", err, "Completed")
		return nil, err
	}

	log.Dev(context, "GetUserByPublicID", "Completed")
	return &user, nil
}

// GetUserByEmail retrieves a user record by using the provided email.
func GetUserByEmail(context interface{}, db *db.DB, email string, activeOnly bool) (*User, error) {
	log.Dev(context, "GetUserByEmail", "Started : Email[%s]", email)

	var user User
	f := func(c *mgo.Collection) error {
		var q bson.M
		if activeOnly {
			q = bson.M{"email": strings.ToLower(email), "status": StatusActive}
		} else {
			q = bson.M{"email": strings.ToLower(email)}
		}
		log.Dev(context, "GetUserByEmail", "MGO : db.%s.findOne(%s)", c.Name, mongo.Query(q))
		return c.Find(q).One(&user)
	}

	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		log.Error(context, "GetUserByEmail", err, "Completed")
		return nil, err
	}

	log.Dev(context, "GetUserByEmail", "Completed")
	return &user, nil
}

// GetUserWebToken return a token if one exists and is valid.
func GetUserWebToken(context interface{}, db *db.DB, publicID string) (string, error) {
	log.Dev(context, "GetUserWebToken", "Started : PublicID[%s]", publicID)

	// Do we have a valid session right now?
	s, err := session.GetByLatest(context, db, publicID)
	if err != nil {
		log.Error(context, "GetUserWebToken", err, "Completed")
		return "", err
	}

	// If it is expired return failure.
	if s.IsExpired(context) {
		err := errors.New("Session expired.")
		log.Error(context, "GetUserWebToken", err, "Completed")
		return "", err
	}

	// Pull the user information.
	u, err := GetUserByPublicID(context, db, publicID, true)
	if err != nil {
		log.Error(context, "GetUserWebToken", err, "Completed")
		return "", err
	}

	// Generate a token that works right now.
	token, err := u.WebToken(s.SessionID)
	if err != nil {
		log.Error(context, "GetUserWebToken", err, "Completed")
		return "", err
	}

	log.Dev(context, "GetUserWebToken", "Completed : WebToken[%s]", token)
	return token, nil
}

//==============================================================================

// UpdateUser updates an existing user to the database.
func UpdateUser(context interface{}, db *db.DB, uu UpdUser) error {
	log.Dev(context, "UpdateUser", "Started : PublicID[%s]", uu.PublicID)

	if err := uu.Validate(); err != nil {
		log.Error(context, "UpdateUser", err, "Completed")
		return err
	}

	f := func(c *mgo.Collection) error {
		q := bson.M{"public_id": uu.PublicID}
		upd := bson.M{"$set": bson.M{"full_name": uu.FullName, "email": uu.Email, "status": uu.Status, "modified_at": time.Now().UTC()}}
		log.Dev(context, "UpdateUser", "MGO : db.%s.Update(%s, %s)", c.Name, mongo.Query(q), mongo.Query(upd))
		return c.Update(q, upd)
	}

	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		log.Error(context, "UpdateUser", err, "Completed")
		return err
	}

	log.Dev(context, "UpdateUser", "Completed")
	return nil
}

// UpdateUserPassword updates an existing user's password and token in the database.
func UpdateUserPassword(context interface{}, db *db.DB, u *User, password string) error {
	log.Dev(context, "UpdateUserPassword", "Started : PublicID[%s]", u.PublicID)

	if err := u.Validate(); err != nil {
		log.Error(context, "UpdateUserPassword", err, "Completed")
		return err
	}

	if len(password) < 8 {
		err := errors.New("Invalid password length")
		log.Error(context, "UpdateUserPassword", err, "Completed")
		return err
	}

	newPassHash, err := crypto.BcryptPassword(u.PrivateID + password)
	if err != nil {
		log.Error(context, "UpdateUserPassword", err, "Completed")
		return err
	}

	f := func(c *mgo.Collection) error {
		q := bson.M{"public_id": u.PublicID}
		upd := bson.M{"$set": bson.M{"password": newPassHash, "modified_at": time.Now().UTC()}}
		log.Dev(context, "UpdateUserPassword", "MGO : db.%s.Update(%s, CAN'T SHOW)", c.Name, mongo.Query(q))
		return c.Update(q, upd)
	}

	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		log.Error(context, "UpdateUserPassword", err, "Completed")
		return err
	}

	log.Dev(context, "UpdateUserPassword", "Completed")
	return nil
}

// UpdateUserStatus changes the status of a user to make them active or disabled.
func UpdateUserStatus(context interface{}, db *db.DB, publicID string, status int) error {
	log.Dev(context, "UpdateUserStatus", "Started : PublicID[%s] Status[%d]", publicID, status)

	if status != StatusActive && status != StatusDisabled {
		err := errors.New("Invalid status code")
		log.Error(context, "LoginUser", err, "Completed")
		return err
	}

	f := func(c *mgo.Collection) error {
		q := bson.M{"public_id": publicID}
		upd := bson.M{"$set": bson.M{"status": status, "modified_at": time.Now().UTC()}}
		log.Dev(context, "UpdateUserStatus", "MGO : db.%s.Update(%s, %s)", c.Name, mongo.Query(q), mongo.Query(upd))
		return c.Update(q, upd)
	}

	if err := db.ExecuteMGO(context, Collection, f); err != nil {
		log.Error(context, "UpdateUserStatus", err, "Completed")
		return err
	}

	log.Dev(context, "UpdateUserStatus", "Completed")
	return nil
}

//==============================================================================

// LoginUser authenticates the user and if successful returns the User value.
func LoginUser(context interface{}, db *db.DB, email string, password string) (*User, error) {
	log.Dev(context, "LoginUser", "Started : Email[%s]", email)

	u, err := GetUserByEmail(context, db, email, true)
	if err != nil {
		log.Error(context, "LoginUser", err, "Completed")
		return nil, err
	}

	if ok := u.IsPasswordValid(password); !ok {
		err := errors.New("Invalid password")
		log.Error(context, "LoginUser", err, "Completed")
		return nil, err
	}

	log.Dev(context, "LoginUser", "Completed")
	return u, nil
}
