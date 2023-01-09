package http

import (
	"context"
	"errors"
	"github.com/avast/retry-go"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func (client *Client) request(ctx context.Context, apiInfo *APIInfo) ([]byte, http.Header, error) {

	var requestBody io.Reader
	if apiInfo.Body != "" {
		requestBody = strings.NewReader(apiInfo.Body)
	}
	req, err := http.NewRequest(apiInfo.Method, client.opts.Url, requestBody)
	if err != nil {
		return []byte(""), nil, errors.New("构建request失败")
	}

	for k, v := range client.opts.Header {
		req.Header.Set(k, strings.Join(v, ";"))
	}

	var resp []byte
	var headers http.Header

	err = retry.Do(func() error {
		var needRetry bool
		resp, _, _, headers, err, needRetry = client.makeRequest(ctx, req, client.opts.Timeout)
		if needRetry {
			return err
		}
		return nil
	}, retry.Attempts(3))

	return resp, headers, err
}

func (client *Client) makeRequest(inputContext context.Context, req *http.Request, timeout time.Duration) ([]byte, int, string, http.Header, error, bool) {
	ctx := inputContext
	if ctx == nil {
		ctx = context.Background()
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	req = req.WithContext(ctx)

	resp, err := client.c.Do(req)
	if err != nil || resp == nil {
		// should retry when client sends request error.
		return []byte(""), http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), nil, err, true
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte(""), resp.StatusCode, resp.Status, resp.Header, err, false
	}

	return body, resp.StatusCode, resp.Status, resp.Header, nil, false
}
