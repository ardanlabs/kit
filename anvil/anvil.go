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

// key contains the JSON WEB KEY required to encrypt the JWT received on
// every request.
type key struct {
	Type      string `json:"kty"` // Key Type: RSA
	Use       string `json:"use"` // Verifiying signatures for `sig` or `enc`
	Algorithm string `json:"alg"` // Algorithm to use: RS256
	Modulus   string `json:"n"`   // The modulus section of the key
	Exponent  string `json:"e"`   // The exponent of the key
}

type keys struct {
	Set []key `json:"keys"`
}

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

// ValidateFromRequest takes a request and extracts the JWT. Then it performs
// validation and returns the claims if everything is valid.
func ValidateFromRequest(r *http.Request, pk *rsa.PublicKey) (Claims, error) {

	// Function is required to return the Public Key that is needed to
	// validate the JWT and extract the claims.
	f := func(token *jwt.Token) (interface{}, error) {
		return pk, nil
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

// RetrievePublicKey calls into Anvil to get the Public Key information.
func RetrievePublicKey(host string) (*rsa.PublicKey, error) {

	// Ask Anvil for the public keys.
	r, err := http.Get(host + "/jwks")
	if err != nil {
		return nil, err
	}

	// Validate we successful requested them.
	if r.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Invalid Status: %d %s", r.StatusCode, r.Status)
	}

	defer r.Body.Close()

	// Decode the two keys we will receive, `sig` and `enc`.
	var ks keys
	if err := json.NewDecoder(r.Body).Decode(&ks); err != nil {
		return nil, err
	}

	// Find the `sig` key since this is what we need.
	for _, jwk := range ks.Set {
		if jwk.Use == "sig" {
			return jwkToPK(jwk)
		}
	}

	return nil, errors.New("Sig keys not found")
}

// jwkToPK converts the Anvil JWK into a RSA Public Key.
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

//==============================================================================

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
