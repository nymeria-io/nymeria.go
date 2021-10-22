package nymeria

import (
	"errors"
	"strings"
)

var apiKey string

var (
	// ErrInvalidAuthKey is returned any time an invalid or malformed API auth
	// key is detected.
	ErrInvalidAuthKey = errors.New(`error: the supplied auth key is invalid`)
)

// SetAuth will set the libraries authentication key. Only a single API key
// can be used and will be added to all API requests automatically.
func SetAuth(s string) error {
	s = strings.TrimSpace(s)

	if len(s) == 0 {
		return ErrInvalidAuthKey
	}

	apiKey = s

	return nil
}
