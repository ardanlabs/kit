// Package anvil provides support for validating an Anvil JWT and extracting
// the claims for authorization.
package anvil

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
)

// Claims is the payload we are looking to get on each request.
type Claims struct {
	Jti   string
	Iss   string
	Sub   string
	Aud   string
	Exp   int64
	Iat   int64
	Scope string // This is where we find the authorization markers.
}

// key contains the JSON WEB KEY required to encrypt the JWT received on
// every request.
type key struct {
	Type      string `json:"kty"` // Key Type: RSA
	Use       string `json:"use"` // Verifiying signatures for `sig` or `enc`
	Algorithm string `json:"alg"` // Algorithm to use: RS256
	Modulus   string `json:"n"`   // The modulus section of the key
	Exponent  string `json:"e"`   // The exponent of the key
}

// keys match the document returned by Anvil.io on the request.
// curl http://HOST/jwks
type keys struct {
	Set []key `json:"keys"`
}

//==============================================================================

// Anvil provides support for validating Anvil.io based JWTs and extracting
// claims for authorization.
type Anvil struct {
	PublicKey *rsa.PublicKey
}

// New create a new Anvil value for use with handling JWTs from Anvil.io
func New(host string) (*Anvil, error) {

	// Ask Anvil.io for the public keys.
	r, err := http.Get(host + "/jwks")
	if err != nil {
		return nil, err
	}

	// Validate we successfully received the keys.
	if r.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Invalid Status: %d %s", r.StatusCode, r.Status)
	}

	defer r.Body.Close()

	// Decode the document we received.
	var ks keys
	if err := json.NewDecoder(r.Body).Decode(&ks); err != nil {
		return nil, err
	}

	// Find the `sig` key since this is what we need.
	var jwk *key
	for _, key := range ks.Set {
		if key.Use == "sig" {
			jwk = &key
			break
		}
	}

	// Did we find the `sig` key.
	if jwk == nil {
		return nil, errors.New("`Sig` key not found")
	}

	// Convert the `sig` key into an rsa public key.
	pk, err := jwkToPK(*jwk)
	if err != nil {
		return nil, err
	}

	return &Anvil{PublicKey: pk}, nil
}

// ValidateFromRequest takes a request and extracts the JWT. Then it performs
// validation and returns the claims if everything is valid.
func (a *Anvil) ValidateFromRequest(r *http.Request) (Claims, error) {

	// Function is required to return the public key that is needed to
	// validate the JWT and extract the claims.
	f := func(token *jwt.Token) (interface{}, error) {
		return a.PublicKey, nil
	}

	// Parse the request, looking for the JWT and peforming transformations.
	token, err := jwt.ParseFromRequest(r, f)
	if err != nil {
		return Claims{}, fmt.Errorf("Parse Error: %v", err)
	}

	// Was the token valid.
	if !token.Valid {
		return Claims{}, errors.New("Token is invalid")
	}

	// Trying to reduce the cost of unmarshaling the map into our struct so
	// doing it manually. Calling Marshal/Unmarshal will cost more.

	var claims Claims

	if v, exists := token.Claims["jti"]; exists {
		claims.Jti = v.(string)
	}

	if v, exists := token.Claims["iss"]; exists {
		claims.Iss = v.(string)
	}

	if v, exists := token.Claims["sub"]; exists {
		claims.Sub = v.(string)
	}

	if v, exists := token.Claims["aud"]; exists {
		claims.Aud = v.(string)
	}

	if v, exists := token.Claims["exp"]; exists {
		claims.Exp = int64(v.(float64))
	}

	if v, exists := token.Claims["iat"]; exists {
		claims.Iat = int64(v.(float64))
	}

	if v, exists := token.Claims["scope"]; exists {
		claims.Scope = v.(string)
	}

	return claims, nil
}

//==============================================================================

// jwkToPK converts the Anvil JWK into a RSA public key.
func jwkToPK(k key) (*rsa.PublicKey, error) {

	// Convert the Modulus into a Big Int.
	n, err := base64RawURLEncToBigInt(k.Modulus)
	if err != nil {
		return nil, err
	}

	// Convert the Exponent into a Big Int.
	e, err := base64RawURLEncToBigInt(k.Exponent)
	if err != nil {
		return nil, err
	}

	// Create an rsa public key value based on the JWK received from Anvil.
	key := rsa.PublicKey{
		N: n,
		E: int(e.Int64()),
	}

	return &key, nil
}

// base64StdEncToBigInt takes a base64 standard encoded string and converts
// it to a Big Int.
func base64StdEncToBigInt(str string) (*big.Int, error) {
	decoded, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return nil, err
	}

	var bint big.Int
	bint.SetBytes([]byte(decoded))

	return &bint, nil
}

// base64RawURLEncToBigInt takes a base64 raw url encoded string and converts
// it to a Big Int.
func base64RawURLEncToBigInt(str string) (*big.Int, error) {
	decoded, err := base64.RawURLEncoding.DecodeString(str)
	if err != nil {
		return nil, err
	}

	var bint big.Int
	bint.SetBytes([]byte(decoded))

	return &bint, nil
}
