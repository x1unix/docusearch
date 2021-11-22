package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type ErrorResponse struct {
	StatusCode int    `json:"-"`
	Status     string `json:"-"`
	Message    string `json:"message"`
}

func (rsp ErrorResponse) Error() string {
	if rsp.Message == "" {
		return rsp.Status
	}

	return fmt.Sprintf("%s (%s)", rsp.Message, rsp.Status)
}

func checkResponseError(r *http.Response) error {
	if r.StatusCode < 400 {
		return nil
	}

	errRsp := &ErrorResponse{StatusCode: r.StatusCode, Status: r.Status}
	contentType := r.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "application/json") {
		return errRsp
	}

	if err := json.NewDecoder(r.Body).Decode(&errRsp); err != nil {
		return err
	}

	return errRsp
}
