package nymeria

import (
	"fmt"
	"io"
	"net/http"
)

var (
	client = http.Client{}
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
