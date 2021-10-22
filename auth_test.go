package nymeria

import "testing"

func TestSetAuthErrsOnBlank(t *testing.T) {
	err := SetAuth("\t  \n   	")

	if err != ErrInvalidAuthKey {
		t.Fatal("error: accepted a blank auth key")
	}
}

func TestStripsAuthKey(t *testing.T) {
	err := SetAuth("\t  \nabc-123   	")

	if err == ErrInvalidAuthKey {
		t.Fatal("error: failed to properly set an auth key")
	}

	if apiKey != "abc-123" {
		t.Fatal("error: failed to strip auth key")
	}
}
