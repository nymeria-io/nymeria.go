package nymeria

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"strings"
)

var (
	// ErrInvalidAuthKey is returned any time an invalid or malformed API auth
	// key is detected.
	ErrInvalidAuthKey = errors.New(`error: the supplied auth key is invalid`)

	// The API key that will be used for all authenticated requests.
	apiKey string
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

// CheckAuthentication will send the apiKey to the Nymeria server, and check
// if the API key is valid. If it's invalid, an error is returned.
func CheckAuthentication() error {
	req, err := request("GET", "/check-authentication", nil)

	if err != nil {
		return err
	}

	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	bs, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	type response struct {
		Status string `json:"status"`
	}

	var jsonResp response

	if err := json.Unmarshal(bs, &jsonResp); err != nil {
		return err
	}

	if jsonResp.Status != "success" {
		return ErrInvalidAuthKey
	}

	return nil
}
