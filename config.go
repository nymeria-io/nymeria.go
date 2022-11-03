package nymeria

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	apiVersion = "4"
	baseURL    = "https://www.nymeria.io/api/v" + apiVersion
	userAgent  = "nymeria.go/" + apiVersion
)

var (
	// Timeout is the default timeout used by the HTTP client (default: 5 seconds).
	Timeout = 30 * time.Second

	// ErrInvalidAuthKey is returned any time an invalid or malformed API auth
	// key is detected.
	ErrInvalidAuthKey = errors.New(`error: the supplied auth key is invalid`)

	// The API key that will be used for all authenticated requests.
	apiKey string

	client = http.Client{
		Timeout: Timeout,
	}
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

func request(method, endpoint string, data io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", baseURL, endpoint), data)

	if err != nil {
		return nil, err
	}

	req.Header = http.Header{
		"X-Api-Key":    []string{apiKey},
		"Content-Type": []string{"application/json"},
		"User-Agent":   []string{userAgent},
	}

	return req, nil
}
