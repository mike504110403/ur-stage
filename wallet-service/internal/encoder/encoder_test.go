package encoder

import (
	"testing"
)

func TestHmac384_Encoder(t *testing.T) {
	testResult := "939a5c75eec90516af72d2a24bf0e71ef4bd2bbb8b37d73499c6b5409646333239d7821f7fe02d0edb27b6c46b8de106"
	signKey := "gjiowtk49Hw3l"
	testOrigin := "afield=aaa&bfield=bBb&cfield=ccC&tfield=TTT"

	result, err := Hmac384_Encoder(testOrigin, []byte(signKey))
	if err != nil {
		t.Errorf("Hmac384_Encoder() error = %v", err)
		return
	}
	if result != testResult {
		t.Errorf("Hmac384_Encoder() = %v, want %v", result, testResult)
	}
}
