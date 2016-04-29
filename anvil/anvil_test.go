package anvil_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ardanlabs/kit/anvil"
)

const succeed = "\u2713"
const failed = "\u2717"

// expPEM contains the PEM document we expect to produce and use.
var expPEM = `-----BEGIN PUBLIC KEY-----
MIICJzANBgkqhkiG9w0BAQEFAAOCAhQAMIICDwKCAgEArlN01z7UMm97vVphJICB
LBFZNv+IMMJq1V/lvprWS96p9s1yiQvlwbxGmSTTqeV4RNeshTfwM6HO/ADEZCP3
PdKLhDMKkqlGP9NLktkdlkalLSJdGyElqRJi9oy7tRmGdYTvI1i7Esup8MadJdFX
RhNwdn/tIHT0uV6SOgX5RtOF3tPybW01gpYqNWW+SWivfbXlC2W/V2BrR/xprDNU
za6BLjmTeUSH7GLoDgt/5OTxYBK0xP3UWlGWQZ8PNUfv/zKvUKvK951doX+WJp92
pK3uS99uIi6lfivHYMX5ncKYY325TXzLgpaBkNH/Uaiw/Lzt3ogaEqxn31eQNuHe
fiOTKazH2e0V61ymdLOn7Gw7ZtzSOBOyFdJzHKprMb94uC1oJYjlUsChIU2vzLr4
X7B48kpo1hhAjkOEc8Jmri4NZBCfo9bUudcRLynNUix6cGD4QnA8fdGi8R6YiTo0
XPOCuJy2K2NtwIdQDRe0CgnhS4EOkMg5Q5YCZytAhRr022sM0JUpNpyZ//IXy/GY
5ZoC16kQ926lVzlHoCbI0UJUpy/425BaDKj7tbVqBYNCHuz2p94v00hFTs6gfYKE
3tOMhPPFd3BhB2Wq2FbT28vmlPcqhr0ZYHVZNQpp33CALQJ1fYDcCg8HJj8R/puT
QLAQFYfJQGBTLJx6x0pYrq0CCCycesdKWK6t
-----END PUBLIC KEY-----
`

//==============================================================================

// mockServer returns the JWKs for the tests.
func mockServer() *httptest.Server {
	jwks := []anvil.JWK{
		{
			KeyType:   "RSA",
			Use:       "sig",
			Algorithm: "RS256",
			Modulus:   "rlN01z7UMm97vVphJICBLBFZNv-IMMJq1V_lvprWS96p9s1yiQvlwbxGmSTTqeV4RNeshTfwM6HO_ADEZCP3PdKLhDMKkqlGP9NLktkdlkalLSJdGyElqRJi9oy7tRmGdYTvI1i7Esup8MadJdFXRhNwdn_tIHT0uV6SOgX5RtOF3tPybW01gpYqNWW-SWivfbXlC2W_V2BrR_xprDNUza6BLjmTeUSH7GLoDgt_5OTxYBK0xP3UWlGWQZ8PNUfv_zKvUKvK951doX-WJp92pK3uS99uIi6lfivHYMX5ncKYY325TXzLgpaBkNH_Uaiw_Lzt3ogaEqxn31eQNuHefiOTKazH2e0V61ymdLOn7Gw7ZtzSOBOyFdJzHKprMb94uC1oJYjlUsChIU2vzLr4X7B48kpo1hhAjkOEc8Jmri4NZBCfo9bUudcRLynNUix6cGD4QnA8fdGi8R6YiTo0XPOCuJy2K2NtwIdQDRe0CgnhS4EOkMg5Q5YCZytAhRr022sM0JUpNpyZ__IXy_GY5ZoC16kQ926lVzlHoCbI0UJUpy_425BaDKj7tbVqBYNCHuz2p94v00hFTs6gfYKE3tOMhPPFd3BhB2Wq2FbT28vmlPcqhr0ZYHVZNQpp33CALQJ1fYDcCg8HJj8R_puTQLAQFYfJQGBTLJx6x0pYrq0",
			Exponent:  "AQAB",
		},
		{
			KeyType:   "RSA",
			Use:       "enc",
			Algorithm: "RS256",
			Modulus:   "2aQ0mVwYhCNr0JijOaq_E47zgWgthZFYZS-zdo9UoKMMyGs_0JTybCZYMc64dQPFAmamBQ8VJcacsF8oAdgWdZAMrXgvxkldLkE9Em_vRdhKjhVkPPBRUSMf6IU78csihuAZ5XsJ4nlUj5fipGaPJuF-PFyBs3Z4rfLCXJCjE7OspgvV13Pgt8R1ucJok204ZyPJ-LonQiqzgWvKm3lj8wVdx6NyozfcTlmMLgWb6HMpsORZ_ZklpDjUwfjlzTYV-wl3pXsXyslGsOVH7ixjLexJyzB5DJXIXRsjiaonvvf1sIOyK0ys3ilbrgv7Is-MfNNxSZ5I7ikCB82fAvt6HAVIY9NjNOjVVQx6AsL_A0YbQvU2kunAvp2GX3knt283O74jWhbgHcE3EjNxLx4EL4CXjTLcQU20-bpaJJmnwozlGsikY2S0Mf7s7VdVjwYeQ6jONfr79QNZKSh0-sbG7a6q7zPQk9rLMq3noKgBnGUkKbaVAARNn6ivCOZw0VGxPRdlGxN_oBXPXMAP54GyzU0wW97Zka_Z9dB22G9TxNg6kBEjLZzY71WWnuQXt7DMunZgLACPrMNpeJnftsk0YPkCxplOjW7ztzSEUg4jYGt5YlgFn5fgcOnAzMbNH5EgJ3K8UDymd_kAQt07PpPIPrQT29s2y7wzhj1iLES4eTc",
			Exponent:  "AQAB",
		},
	}

	f := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(jwks)
	}

	return httptest.NewServer(http.HandlerFunc(f))
}

// mockHandler executes the full path of the package for processing
// Anvil based JWTs.
func mockHandler(rw http.ResponseWriter, r *http.Request) {
	server := mockServer()
	defer server.Close()

	pem, err := anvil.RetrievePEM(server.URL)
	if err != nil {
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(501)
		json.NewEncoder(rw).Encode(struct{ Error string }{err.Error()})
		return
	}

	claims, err := anvil.ValidateFromRequest(r, pem)
	if err != nil {
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(502)
		json.NewEncoder(rw).Encode(struct{ Error string }{err.Error()})
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(200)
	json.NewEncoder(rw).Encode(claims)
}

//==============================================================================

// TestRetrievePEM validates we can retrieve the JWKs and convert them to a
// PEM document.
func TestRetrievePEM(t *testing.T) {
	server := mockServer()
	defer server.Close()

	t.Log("Given the need to retrieve the JWKs.")
	{
		t.Logf("\tTest 0:\tWhen reuqesting %q", server.URL)
		{
			pem, err := anvil.RetrievePEM(server.URL)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve the PEM : %v", failed, err)
			}
			t.Logf("\t%s\tShould be able to retrieve the PEM.", succeed)

			rcvPEM := string(pem)

			if rcvPEM != expPEM {
				t.Logf("\tRCV\n%+v", rcvPEM)
				t.Logf("\tEXP\n%+v", expPEM)
				t.Errorf("\t%s\tShould have the correct PEM document.", failed)
			} else {
				t.Logf("\t%s\tShould have the correct PEM document.", succeed)
			}
		}
	}
}

// TestValidateFromRequest validates a JWT can be processed from a request and
// validated with proper claims extraction.
func TestValidateFromRequest(t *testing.T) {
	http.HandleFunc("/api", mockHandler)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/api", nil)
	r.Header.Add("Authorization", "Bearer eyJhbGciOiJSUzI1NiJ9.eyJqdGkiOiJhNTdjNTI4ZTFiNDkzNWNiODYxYSIsImlzcyI6Imh0dHBzOi8vZm9yZ2UuYW52aWwuaW8iLCJzdWIiOiI0ZTVhYjg3NS1mZWViLTRiYWItODlhMS1jY2MwNjBhNDZiZjciLCJhdWQiOiI2YzVmNWEwZC05Mzk1LTQ0ZGMtYTQ1OC0yOGYyNmM0OTNmOTMiLCJleHAiOjE0NjA5OTcxNTksImlhdCI6MTQ2MDk5MzU1OSwic2NvcGUiOiJvcGVuaWQgcHJvZmlsZSByZWFsbSJ9.PAnBQRgZ-EtetYrR1pb7AHY3A43QB_quTC24DC6SGnACehuOHhLgZ7e6XxS1aZsaHI5Jt0MW4bJMA1mMfGYOOl7naSlFFnEbWWo4YOkZjoG--cDDWh4N5Aim1Y8cUQszyTZahxizx9NglAmzwUr31EEyTF8Yj0xp9FlWUqX0Dr9YHqj7iwU01fKruZ-9iG2ckhnvQ3mfLwKq-lJ0cBj8iGQYLktZ33GDrwvNcRdBnbHl-q_5pgKENbsXPpN6_FuaJrZa-GgnSbcrVM3tHNngubLxtCgueKzKBsppxTWsBRKSWa1lx6ZnBDbkSDOk1wsaQQeIZqwefxSB8q6lrVOoVqeZeRgAYJgoqW8zdX7kX_LCoEZ8i0-f0AQpemPHfa-K0XtSE3ENrz-9gWNbbxtgcnr3RPXmY4ke5-m2G4cb7F9iilgX9U8BOCBvdwaEtiA4O3YUnHlOGCGWV7dPN034iC69-TRUNlXt83GD-algAEnL-q_FrlvTFuNa0WyNgH4BJ8Si0BWUphqkrj4WDAgl4O8gfevEb-rMGtPlbcoAJnSpPUBDjFyZzZdeLJzM-B8uenKCBCL3LqtcFpp8k7_I3Q9bUC9fbjawG1IBVtZEucxWGvytOXVuo8S_VS7z3iGVvuDUEWYTPITsO2EE5VG3kFy8KTIjwOUBULVt7buXvus")

	http.DefaultServeMux.ServeHTTP(w, r)

	t.Log("Given the need to validate JWTs from Anvil.")
	{
		t.Log("\tTest 0:\tWhen making a stadard web api call.")
		{
			if w.Code != 200 {
				t.Errorf("\t%s\tShould receive a status code of 200 for the response. Received[%d].", failed, w.Code)

				if w.Code >= 501 {
					var rspErr struct{ Error string }
					if err := json.NewDecoder(w.Body).Decode(&rspErr); err != nil {
						t.Fatalf("\t%s\tShould be able to decode the error response : %v", failed, err)
					}

					t.Fatalf("\t%s\t%+v", failed, rspErr)
				}
			}
			t.Logf("\t%s\tShould receive a status code of 200 for the response.", succeed)

			var claims anvil.Claims
			if err := json.NewDecoder(w.Body).Decode(&claims); err != nil {
				t.Fatalf("\t%s\tShould be able to decode the claims response : %v", failed, err)
			}
			t.Logf("\t%s\tShould be able to decode the claims response.", succeed)

			t.Logf("%+v", claims)
		}
	}
}
