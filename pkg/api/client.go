package api

import (
	"io"
	"io/ioutil"
	"net/http"
	"path"
)

type Client struct {
	baseUrl string
}

func NewClient(baseUrl string) *Client {
	return &Client{baseUrl: baseUrl}
}

func (c Client) AddDocument(name string, data io.Reader) error {
	r, err := http.NewRequest(http.MethodPost, c.baseUrl+"/"+path.Join("document", name), data)
	if err != nil {
		return err
	}

	rsp, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}

	defer rsp.Body.Close()
	return checkResponseError(rsp)
}

func (c Client) GetDocument(name string) ([]byte, error) {
	r, err := http.NewRequest(http.MethodGet, c.baseUrl+"/"+path.Join("document", name), nil)
	if err != nil {
		return nil, err
	}

	rsp, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, err
	}

	if err := checkResponseError(rsp); err != nil {
		return nil, err
	}

	defer rsp.Body.Close()
	return ioutil.ReadAll(rsp.Body)
}

func (c Client) RemoveDocument(name string) error {
	r, err := http.NewRequest(http.MethodDelete, c.baseUrl+"/"+path.Join("document", name), nil)
	if err != nil {
		return err
	}

	rsp, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}

	defer rsp.Body.Close()
	return checkResponseError(rsp)
}
