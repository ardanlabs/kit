package anvil_test

import (
	"testing"

	"github.com/ardanlabs/kit/anvil"
)

func TestConversion(t *testing.T) {
	n := "rlN01z7UMm97vVphJICBLBFZNv-IMMJq1V_lvprWS96p9s1yiQvlwbxGmSTTqeV4RNeshTfwM6HO_ADEZCP3PdKLhDMKkqlGP9NLktkdlkalLSJdGyElqRJi9oy7tRmGdYTvI1i7Esup8MadJdFXRhNwdn_tIHT0uV6SOgX5RtOF3tPybW01gpYqNWW-SWivfbXlC2W_V2BrR_xprDNUza6BLjmTeUSH7GLoDgt_5OTxYBK0xP3UWlGWQZ8PNUfv_zKvUKvK951doX-WJp92pK3uS99uIi6lfivHYMX5ncKYY325TXzLgpaBkNH_Uaiw_Lzt3ogaEqxn31eQNuHefiOTKazH2e0V61ymdLOn7Gw7ZtzSOBOyFdJzHKprMb94uC1oJYjlUsChIU2vzLr4X7B48kpo1hhAjkOEc8Jmri4NZBCfo9bUudcRLynNUix6cGD4QnA8fdGi8R6YiTo0XPOCuJy2K2NtwIdQDRe0CgnhS4EOkMg5Q5YCZytAhRr022sM0JUpNpyZ__IXy_GY5ZoC16kQ926lVzlHoCbI0UJUpy_425BaDKj7tbVqBYNCHuz2p94v00hFTs6gfYKE3tOMhPPFd3BhB2Wq2FbT28vmlPcqhr0ZYHVZNQpp33CALQJ1fYDcCg8HJj8R_puTQLAQFYfJQGBTLJx6x0pYrq0"
	e := "AQAB"

	jwk := anvil.JWK{
		Modulus:  n,
		Exponent: e,
	}

	pem, err := anvil.JWKToPEM(jwk)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(pem))
}
