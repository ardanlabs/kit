package anvil_test

/*
curl -X POST http://10.0.1.26:3000/signin -d 'response_type=code&client_id=6b6efaae-0ab8-4152-8f92-a87c17921800&redirect_uri=https://anvil.coralproject.net&scope=openid%20profile%20email%20realm&provider=password&email=bill@thekennedyclan.net&password=Qfe^bJ9uD6cgnD-8' -H "referrer: https://anvil.coralproject.net/signin"
Redirecting to https://anvil.coralproject.net?code=c9ce6c03ea6ad8dd3f0a%
curl -X POST http://6b6efaae-0ab8-4152-8f92-a87c17921800:6dafd2b59d6954849a6c@10.0.1.26:3000/token -d 'grant_type=authorization_code&client_id=6b6efaae-0ab8-4152-8f92-a87c17921800&code=b8685e59be97d3cac019&redirect_uri=https://anvil.coralproject.net' -H "referrer: https://anvil.coralproject.net/token"
*/

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ardanlabs/kit/anvil"
)

const succeed = "\u2713"
const failed = "\u2717"

// Response for Anvil.io when requesting keys.
// curl http://10.0.1.26:3000/jwks
var jwks = `{
   "keys":[
      {
         "kty":"RSA",
         "use":"sig",
         "alg":"RS256",
         "n":"wF9BsTKc9DVbWYF359giSHrO1iF2MYeCuUWsOz7xIMoSOoWBBQwWbKC-nTcV12wpycKNGC0AaHjYFgNbuIihSUbDZYNdAba1ecn28LsRzE1X3949kVygeWkG3wyjCSNifuy4g1temjhYrsmpCXIVMUww3-IfqKQZ7aPTf66NgTNpVZzKYSkKHeXPH_n8vHynivEaz_PjuCrc4rv39JgUTMBzpQkJpb3nf-UCZw8AoJCE3qHVcXFIuc2ksEzLjMFu5yECuBFP3Yuq9dRH3aA05Q8E5eV-B_8YGTgWKgfYricE440hX1Tf0_o_JUTlH8FC2Wffa65hon5aojFQ7h9HU6nQRJedHwK1e6g3pqu3PMtmbC3UEEn1vbYzTxs-ARQA18j4nWA2NsHSyE5jIws6KCdqbL8lm4NYoPa01BdJj_Vm1SEFrlcVDtUnxovwxQZ4tVrRRLcoDCEX_y_Maw4HdM63XtP6j1HjNWPU_7kD5v2PJ9l5Ew5YU66WNvT0cMRpovUbIkk0VTzHyptBRk4zzRi9zru_eUNuZ8KCSQyMAG6umtfNOaCiyTHCW3lFYv2AFXqltKLG14cugRx85NtR2LaPC0dlwS2jW3-09pu4HtzjEuN45fboV9V1Cav4LRr4wrY_PMU-YlpOeD80Xjla6gDj2_qTwW-5vQ94EY6-PoM",
         "e":"AQAB"
      },
      {
         "kty":"RSA",
         "use":"enc",
         "alg":"RS256",
         "n":"pim0jTRQqD3_tXUqcMmz2CgYjsFlNbPmAwDFqmBqxAsSH2gYCpImqpZUoo9Cy-E0W63kMFl-SfXbM9uEmYDXgywfNiA1TPV1curfKFOCbUCXTKI2Hlol1pltgm2cZvJhu_zmo8YFH9Zp74XkA2XQDW-4Ri1uBgAV4_e7v4l2WA27OKRP7jPho_Kjg99ILlgNdLrwl0FYct0xSZ_eGn6M931lhtJWLGcnJlM4eIC6dcIuMjL51czgVKjgjnrjYHMHjoQLTs7PTPF_c_ojUhFM2fUtoi2eaNAGhQ8J1aJ2KopoqcFzo5phRvsR5GyGzasZK5fzCccm5QwFcIKwJ26_0jXereHNz7RQTjms1osJMKckDKCQx6u2U_bXXPuM0g0xy97e174YiCb_9QYMcAUeYz8IUtrsAHXDHu9XNB4seGDMTOTrqgJCIrzoJEIFA1sa83goXpkUHPCBYFjNdxJJbnfOzO2KvYOueO94LVYU3-1kEG3PibbFKIMv17PnKnFY6GQt6OhihHLbxyOL7gQ7wWNruVZTfT59MnxBm-yKe5U46lbn4uFhbeAU3iui-N6XyZ9jTzINhKvJnJL1Ukn_k5bhektxd7shQyTpILeK1TmnG8jwHXUpQUqHDPrfWARxhMPrSAkmePISaOop0HXSC5TmalI0EEqrEUI_77dD05E",
         "e":"AQAB"
      }
   ]
}`

//==============================================================================

// mockServer returns the JWKs for the tests.
func mockServer() *httptest.Server {
	f := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, jwks)
	}

	return httptest.NewServer(http.HandlerFunc(f))
}

// mockHandler executes the full path of the package for processing
// Anvil based JWTs.
func mockHandler(rw http.ResponseWriter, r *http.Request) {
	server := mockServer()
	defer server.Close()

	pk, err := anvil.RetrievePublicKey(server.URL)
	if err != nil {
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(501)
		json.NewEncoder(rw).Encode(struct{ Error string }{err.Error()})
		return
	}

	claims, err := anvil.ValidateFromRequest(r, pk)
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

// TestRetrievePublicKey validates we can retrieve the JWKs properly.
func TestRetrievePublicKey(t *testing.T) {
	server := mockServer()
	defer server.Close()

	t.Log("Given the need to retrieve the JWKs.")
	{
		t.Logf("\tTest 0:\tWhen reuqesting %q", server.URL)
		{
			pk, err := anvil.RetrievePublicKey(server.URL)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve the PK : %v", failed, err)
			}
			t.Logf("\t%s\tShould be able to retrieve the PK.", succeed)

			ret, err := json.MarshalIndent(pk, "", "")
			if err != nil {
				t.Fatalf("\t%s\tShould be able to marshal the PK : %v", failed, err)
			}
			t.Logf("\t%s\tShould be able to marshal the PK.", succeed)

			// This is the expected result after Marshaling.
			exp := `{
"N": 784809680841980536561529838924636149103449928316484596406116945129881369840077837517596139572912190395854257993301789733592665445295898859892166274020932183697143377939747435366813124394058378213343980741184145376336476231836861359530306236756445340904488118055275444552511730853241683175503268766877900659616390486338132043563565552190809445762864794696718346026795808446765496336615161884911532204322302445505437739851411975361730313913000854404600624939943612536789426331981633743258595734028526125833761419421318414177820312647507263252939737663321336263154244298457031507152579300006431471179894141636035549223078005480851679579386015857523948100201636288760151106160096865536375782411635750471179896321131039674803909899472377216822782752586581532763322460715959009431778261210331291249969189919890579549046585440212724813640517622216207319083857607999373960557900077464606985837960547916645629834135531265952114715064375896618901220309888450784118555943640521327502847292719084166217091145498212058167824103241846701097960044890602248440757107759469336921832091964708017413835057853122900526586776846752370917199664338287767228433527211335710460098655068102431991879967268769339131977905835968526713868975639404449718568107651,
"E": 65537
}`
			if string(ret) != exp {
				t.Logf("\tRCV\n%+v", string(ret))
				t.Logf("\tEXP\n%+v", exp)
				t.Errorf("\t%s\tShould have the correct PK document.", failed)
			} else {
				t.Logf("\t%s\tShould have the correct PK document.", succeed)
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
	r.Header.Add("Authorization", "Bearer eyJhbGciOiJSUzI1NiJ9.eyJqdGkiOiI1ZGY0MzE4NmEyYjM4ZWI1Nzc4MSIsImlzcyI6Imh0dHBzOi8vYW52aWwuY29yYWxwcm9qZWN0Lm5ldCIsInN1YiI6IjI3MzdmMjllLWE5NWItNGVhOC1iNGQxLTMxZDE2YjIzZGVlZSIsImF1ZCI6IjZiNmVmYWFlLTBhYjgtNDE1Mi04ZjkyLWE4N2MxNzkyMTgwMCIsImV4cCI6MTQ2MjAyNjU1NCwiaWF0IjoxNDYyMDIyOTU0LCJzY29wZSI6Im9wZW5pZCBwcm9maWxlIGVtYWlsIHJlYWxtIn0.uEEIYHroZFnYCF3P5tXNJvmbMydBDSlWDYzO2Cx9HC6O8QUr-9r7hmX2bt3t_GUb3VYHMFZLwedB4ZQ2KzDnld3VKrxLzzMf7aMUeCDkOoyiinBiCiy2SIWBHG4cicI2FN7ANWl3Ok21Mp1umTDJwYNbzSi-SnL0n6Mk3FH4_hJM6zKoKQ8qijRXSo8WD4OCjez5iwBZyBrCZPgc1tngB7R1hT57i5j0gDHjoRvGba0YQe-SFp19BvjGZ1h_DmDQi59DCtPvH6BnrVx79nobd6vIFNVuraC8UM72Rx3N2GhNHV_msJWp_wJZfkc8_1gYtrLjVFN8cOkXWJWbZ66LAFqFUNhCyCFwGr4b2FqYePcAYgdsM9gEIcPLtPfO-HthL2AOhGG_oF0VxmyHyGSloaqN3svk7coKdceSK-fAx0vVS3LyW62wSt1UkzzdxXZoGFOdgaB8k65feD6553SFuE2mIA5KwAVd9cHgeX5j6bX6TE_VDwPhTkLH199kzE66rOtoPIUpO3rQYHWdMFIONNHuJtSJmU-gNCF4GVBYCGqx8WiHNE4tKplVRuhpF28RYGEemjsdeStovvBJGKki7B3NcDY-eIopwY3qyBZ72lPONeLTP3zeldu7_O4MQsHOddXFoO_HuWRUgrMy38zIrcgJdoyRW9IHummZc2WS--8")

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

			exp := "openid profile email realm"
			if claims.Scope != exp {
				t.Logf("\tRCV\n%+v", claims.Scope)
				t.Logf("\tEXP\n%+v", exp)
				t.Errorf("\t%s\tShould have the correct Scope : %s", failed, claims.Scope)
			} else {
				t.Logf("\t%s\tShould have the correct Scope.", succeed)
			}
		}
	}
}
