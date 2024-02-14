package httpcall

import (
	"bytes"
	"crypto/tls"
	"github.com/skhatri/go-http-cache/pkg/target/model"
	"io"
	"net/http"
	"os"
	"strings"
)

type _httpCall struct {
	client   *http.Client
	notifier model.Notifier
}

func (hc *_httpCall) Invoke(request model.Request) (*model.Response, error) {
	req, err := http.NewRequest(request.Method, request.Url, bytes.NewReader(request.Body))
	for k, values := range request.Headers {
		for _, value := range values {
			req.Header.Add(strings.ToUpper(k), value)
		}
	}
	if err != nil {
		return nil, err
	}
	res, er := hc.client.Do(req)
	if er != nil {
		return nil, er
	}
	defer res.Body.Close()
	data, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		return nil, readErr
	}
	resp := &model.Response{
		StatusCode: res.StatusCode,
		Headers:    res.Header,
		Data:       io.NopCloser(bytes.NewReader(data)),
	}
	if hc.notifier != nil {
		hc.notifier.OnNotify(request, &model.Response{
			StatusCode: res.StatusCode,
			Headers:    res.Header,
			Data:       io.NopCloser(bytes.NewReader(data)),
		})
	}
	return resp, nil
}
func httpClient() *http.Client {
	skipVerify := false
	if skipFlag := os.Getenv("SKIP_VERIFY_TLS"); skipFlag != "" {
		skipVerify = strings.ToLower(skipFlag) == "true"
	}
	var transport *http.Transport
	transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: skipVerify},
	}
	client := http.Client{
		Transport: transport,
	}
	return &client
}

func New(notifier model.Notifier) model.ResourceClient {
	return &_httpCall{
		client:   httpClient(),
		notifier: notifier,
	}
}
