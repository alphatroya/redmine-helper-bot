package redmine

import "net/http"

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
	response.Body = &bodyMock{}
	return response, c.mockError
}

type bodyMock struct{}

func (b *bodyMock) Read(p []byte) (n int, err error) {
	return 0, nil
}

func (b *bodyMock) Close() error {
	return nil
}
