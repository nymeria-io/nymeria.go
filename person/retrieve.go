package person

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/nymeriaio/nymeria.go"
	"github.com/nymeriaio/nymeria.go/internal/api"
)

type BulkRetrieveParams struct {
	ID       string                 `json:"id"`
	MetaData map[string]interface{} `json:"metadata"`
}

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

	bs, err := io.ReadAll(resp.Body)

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

func BulkRetrieve(params ...BulkRetrieveParams) ([]Person, error) {
	if len(params) == 0 {
		return nil, nymeria.ErrInvalidParameters
	}

	bs, err := json.Marshal(map[string]interface{}{
		"requests": params,
	})

	if err != nil {
		return nil, err
	}

	req, err := api.Request("POST", "/person/retrieve/bulk", bytes.NewBuffer(bs))

	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

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

	bs, err = io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var response struct {
		Status int `json:"status"`
		Data   []struct {
			Status   int                    `json:"status"`
			MetaData map[string]interface{} `json:"metadata"`
			Data     Person                 `json:"data"`
		} `json:"data"`
	}

	if err := json.Unmarshal(bs, &response); err != nil {
		return nil, err
	}

	var records []Person

	for _, v := range response.Data {
		if v.Status == 200 {
			records = append(records, v.Data)
		}
	}

	return records, nil
}
