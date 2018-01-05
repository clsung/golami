package golami

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// constants
const (
	APIEndpointCN = "https://cn.olami.ai/cloudservice/api"
	APIEndpointTW = "https://tw.olami.ai/cloudservice/api"

	APIServiceSEG = "seg"
	APIServiceNLI = "nli"
	APIServiceASR = "asr"
)

// Client type
type Client struct {
	appKey     string
	appSecret  string
	endpoint   *url.URL     // default APIEndpointTW
	httpClient *http.Client // default http.DefaultClient
}

// ClientOption type
type ClientOption func(*Client) error

// New returns a new olami client instance.
func New(appKey, appSecret string, options ...ClientOption) (*Client, error) {
	if appKey == "" {
		return nil, errors.New("missing application key")
	}
	if appSecret == "" {
		return nil, errors.New("missing application secret")
	}
	c := &Client{
		appKey:     appKey,
		appSecret:  appSecret,
		httpClient: http.DefaultClient,
	}
	for _, option := range options {
		err := option(c)
		if err != nil {
			return nil, err
		}
	}
	if c.endpoint == nil {
		c.endpoint, _ = url.ParseRequestURI(APIEndpointTW)
	}
	return c, nil
}

// WithHTTPClient function
func WithHTTPClient(c *http.Client) ClientOption {
	return func(client *Client) error {
		client.httpClient = c
		return nil
	}
}

// WithLocalization function
func WithLocalization(location string) ClientOption {
	return func(client *Client) error {
		switch location {
		case "tw":
			client.endpoint, _ = url.ParseRequestURI(APIEndpointTW)
		case "cn":
			client.endpoint, _ = url.ParseRequestURI(APIEndpointCN)
		default:
			return errors.New("missing location, specify tw or cn")
		}
		return nil
	}
}

func (client *Client) url() string {
	return client.endpoint.String()
}

func (client *Client) do(ctx context.Context, req *http.Request) (*http.Response, error) {
	if ctx != nil {
		req = req.WithContext(ctx)
	}
	return client.httpClient.Do(req)

}

func (client *Client) get(ctx context.Context, query url.Values) (*http.Response, error) {
	req, err := http.NewRequest("GET", client.url(), nil)
	if err != nil {
		return nil, err
	}
	if query != nil {
		req.URL.RawQuery = query.Encode()
	}
	return client.do(ctx, req)
}
