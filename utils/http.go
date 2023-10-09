package utils

import (
	"io"
	"net/http"
)

func HttpQuery(route string) ([]byte, error) {
	req, err := http.NewRequest("GET", route, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	return body, err
}
