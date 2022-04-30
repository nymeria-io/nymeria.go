package nymeria

import "errors"

var (
	// ErrInvalidRequest is returned any time a request looks malformed or inavlid.
	// This can be due to a bad parameter, or bad encoding.
	ErrInvalidRequest = errors.New(`error: the request looks invalid or malformed`)
)
