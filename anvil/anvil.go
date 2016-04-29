// Package anvil provides support for validating an Anvil JWT and extracting
// the claims for authorization.
package anvil

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
)

// JWK contains the JSON WEB KEY required to encrypt the JWT received on
// every request.
type JWK struct {
	KeyType   string `json:"kty"` // Key Type: RSA
	Use       string `json:"use"` // Verifiying signatures for `sig` or `enc`
	Algorithm string `json:"alg"` // Algorithm to use: RS256
	Modulus   string `json:"n"`   // The modulus section of the key
	Exponent  string `json:"e"`   // The exponent of the key
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
func ValidateFromRequest(r *http.Request, pem []byte) (Claims, error) {

	// Function is required to return the PEM that is needed to
	// validate the JWT and extract the claims.
	f := func(token *jwt.Token) (interface{}, error) {
		return pem, nil
	}

	// Parse the request, looking for the JWT and peforming transformations.
	token, err := jwt.ParseFromRequest(r, f)
	if err != nil {
		return Claims{}, err
	}

	// Was the token valid.
	if !token.Valid {
		return Claims{}, errors.New("Token is invalid")
	}

	// Trying to reduce the cost of unmarshaling the map into our struct so
	// doing it manually. Calling Marshal/Unmarshal will cost more.

	var claims Claims

	if v, exists := token.Claims["Jti"]; exists {
		claims.Jti = v.(string)
	}

	if v, exists := token.Claims["Iss"]; exists {
		claims.Iss = v.(string)
	}

	if v, exists := token.Claims["Sub"]; exists {
		claims.Sub = v.(string)
	}

	if v, exists := token.Claims["Aud"]; exists {
		claims.Aud = v.(string)
	}

	if v, exists := token.Claims["Exp"]; exists {
		claims.Exp = int64(v.(int))
	}

	if v, exists := token.Claims["Iat"]; exists {
		claims.Iat = int64(v.(int))
	}

	if v, exists := token.Claims["Scope"]; exists {
		claims.Scope = v.(string)
	}

	fmt.Println("***********", claims)

	return claims, nil
}

//==============================================================================

// RetrievePEM makes a call to Anvil and retrieves the public key for validating
// and extracting the JWT payload. It returns a PEM document for validating and
// extracting the JWT's claims.
func RetrievePEM(host string) ([]byte, error) {

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
	var jwks []JWK
	if err := json.NewDecoder(r.Body).Decode(&jwks); err != nil {
		return nil, err
	}

	// Find the `sig` key since this is what we need.
	for _, jwk := range jwks {
		if jwk.Use == "sig" {
			return JWKToPEM(jwk)
		}
	}

	return nil, errors.New("Sig keys not found")
}

// JWKToPEM takes an Anvil JWK and converts it to the PEM format.
func JWKToPEM(jwk JWK) ([]byte, error) {

	// Convert the Modulus into a Big Int.
	n, err := base64RawURLEncToBigInt(jwk.Modulus)
	if err != nil {
		return nil, err
	}

	// Convert the Exponent into a Big Int.
	e, err := base64RawURLEncToBigInt(jwk.Modulus)
	if err != nil {
		return nil, err
	}

	// Create an rsa public key value based on the JWK received from Anvil.
	key := rsa.PublicKey{
		N: n,
		E: int(e.Int64()),
	}

	// Serialize a public key to DER-encoded PKIX format.
	pubBytes, err := x509.MarshalPKIXPublicKey(&key)
	if err != nil {
		return nil, err
	}

	// Create a PEM block value using the serialized DER-encoded PKIX format.
	block := pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubBytes,
	}

	// Generate the PEM document.
	var buf bytes.Buffer
	if err := pem.Encode(&buf, &block); err != nil {
		return nil, err
	}

	// Return the PEM document.
	return buf.Bytes(), nil
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
