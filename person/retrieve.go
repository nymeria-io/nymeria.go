package person

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"git.nymeria.io/nymeria.go"
	"git.nymeria.io/nymeria.go/internal/api"
)

func Retrieve(id string) (*Person, error) {
	if len(id) == 0 {
		return nil, nymeria.ErrInvalidParameters
	}

	req, err := api.Request("GET", fmt.Sprintf("/person/retrieve/%s", id), nil)

	if err != nil {
		return nil, err
	}

	resp, err := api.Client.Do(req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		if e, ok := nymeria.ErrMap[resp.StatusCode]; ok {
			return nil, e
		}

		return nil, nymeria.ErrServerError
	}

	defer resp.Body.Close()

	bs, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var response struct {
		Status int    `json:"status"`
		Data   Person `json:"data"`
	}

	if err := json.Unmarshal(bs, &response); err != nil {
		return nil, err
	}

	return &response.Data, nil
}
