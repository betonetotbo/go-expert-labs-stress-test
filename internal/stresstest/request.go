package stresstest

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
	"time"
)

func makeRequest(url, method string, headers []string, data []byte, timeout time.Duration) int {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var body io.Reader
	if data != nil {
		body = bytes.NewReader(data)
	}
	req, _ := http.NewRequestWithContext(ctx, method, url, body)
	if headers != nil {
		for _, h := range headers {
			parts := strings.Split(h, ":")
			req.Header.Set(parts[0], parts[1])
		}
	}

	resp, _ := http.DefaultClient.Do(req)

	if resp != nil {
		_ = resp.Body.Close()
		return resp.StatusCode
	}
	return 0
}
