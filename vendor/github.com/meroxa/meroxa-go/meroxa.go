package meroxa

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

const (
	apiURL          = "https://api.meroxa.io/v1"
	jsonContentType = "application/json"
	textContentType = "text/plain"
	clientTimeOut   = 5 * time.Second
)

// encodeFunc encodes v into w
type encodeFunc func(w io.Writer, v interface{}) error

func jsonEncodeFunc(w io.Writer, v interface{}) error {
	return json.NewEncoder(w).Encode(v)
}

func stringEncodeFunc(w io.Writer, v interface{}) error {
	if s, ok := v.(string); ok {
		_, err := w.Write([]byte(s))
		return err
	}

	return fmt.Errorf("Body is not a string")
}

func noopEncodeFunc(w io.Writer, v interface{}) error {
	return nil
}

// Client represents the Meroxa API Client
type Client struct {
	BaseURL   *url.URL
	userAgent string
	token     string

	httpClient *http.Client
}

// New returns a configured Meroxa API Client
func New(token, ua string, debugMode bool) (*Client, error) {
	u, err := url.Parse(getAPIURL())
	if err != nil {
		return nil, err
	}

	c := &Client{
		BaseURL:    u,
		userAgent:  userAgent(ua),
		token:      token,
		httpClient: httpClient(),
	}

	if debugMode {
		c.httpClient.Transport = &DumpTransport{http.DefaultTransport}
	}

	return c, nil
}

func (c *Client) MakeRequestString(ctx context.Context, method, path string, body string) (*http.Response, error) {
	return c.makeRequestRaw(ctx, method, path, body, nil, stringEncodeFunc)
}

func (c *Client) makeRequest(ctx context.Context, method, path string, body interface{}, params url.Values) (*http.Response, error) {
	return c.makeRequestRaw(ctx, method, path, body, params, jsonEncodeFunc)
}

func (c *Client) makeRequestRaw(ctx context.Context, method, path string, body interface{}, params url.Values, encode encodeFunc) (*http.Response, error) {
	req, err := c.newRequest(ctx, method, path, body, encode)
	if err != nil {
		return nil, err
	}

	// Merge params
	if params != nil {
		q := req.URL.Query()
		for k, v := range params { // v is a []string
			for _, vv := range v {
				q.Add(k, vv)
			}
			req.URL.RawQuery = q.Encode()
		}
	}
	resp, err := c.do(ctx, req)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) newRequest(ctx context.Context, method, path string, body interface{}, encode encodeFunc) (*http.Request, error) {
	u, err := c.BaseURL.Parse(path)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	if body != nil {
		if err := encode(buf, body); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	// Set Auth
	bearer := fmt.Sprintf("Bearer %s", c.token)

	// add authorization header to the req
	req.Header.Add("Authorization", bearer)
	req.Header.Add("Content-Type", jsonContentType)
	req.Header.Add("Accept", jsonContentType)
	req.Header.Add("User-Agent", c.userAgent)
	return req, nil
}
func (c *Client) do(ctx context.Context, req *http.Request) (*http.Response, error) {
	req = req.WithContext(ctx)
	return c.httpClient.Do(req)
}

func httpClient() *http.Client {
	return &http.Client{
		Timeout: clientTimeOut,
	}
}

func getAPIURL() string {
	if u := os.Getenv("API_URL"); u != "" {
		return u
	}

	return apiURL
}

func userAgent(ua string) string {
	if ua != "" {
		return ua
	}
	return "meroxa-go"
}
