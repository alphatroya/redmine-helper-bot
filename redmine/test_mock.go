package redmine

import (
	"io"
	"net/http"
	"strings"
)

type ClientRequestMock struct {
	statusCode int
	mockError  error
}

func (c *ClientRequestMock) Do(req *http.Request) (*http.Response, error) {
	response := &http.Response{}
	if c.statusCode != 0 {
		response.StatusCode = c.statusCode
	} else {
		response.StatusCode = 200
	}
	reader := strings.NewReader(`{ "test": "test" }`)
	response.Body = &bodyMock{reader}
	return response, c.mockError
}

type bodyMock struct{
	reader io.Reader
}

func (b *bodyMock) Read(p []byte) (n int, err error) {
	return b.reader.Read(p)
}

func (b *bodyMock) Close() error {
	return nil
}
