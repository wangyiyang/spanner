package rest

import (
	"io"
	"net/http"
	"strings"
)

func HttpDo(url, body, method string, headers map[string][]string) (result http.Response, err error) {
	var payload io.Reader
	if body != "" {
		payload = strings.NewReader(body)
	}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return
	}
	for key, values := range headers {
		req.Header.Add(key, strings.Join(values, ","))
	}
	res, err := http.DefaultClient.Do(req)
	return *res, err
}
