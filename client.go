package nymeria

import (
	"fmt"
	"io"
	"net/http"
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

	client = http.Client{
		Timeout: Timeout,
	}
)

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
