package golami

import (
	"bytes"
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

// constants
const (
	APIEndpointCN = "https://cn.olami.ai/cloudservice/api"
	APIEndpointTW = "https://tw.olami.ai/cloudservice/api"

	APIServiceSEG = "seg"
	APIServiceNLI = "nli"
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

func (client *Client) Post(ctx context.Context, service, text string) (*http.Response, error) {
	// get timestamp
	timeStamp := time.Now().Local().UnixNano() / int64(time.Millisecond)
	log.Printf("timestamp: %d\n", timeStamp)

	//  Prepare message to generate an MD5 digest.
	signMsg := fmt.Sprintf("%sapi=%sappkey=%stimestamp=%d%s",
		client.appSecret, service, client.appKey, timeStamp, client.appSecret,
	)

	// Generate MD5 digest.
	sign := fmt.Sprintf("%x", md5.Sum([]byte(signMsg)))

	// Prepare rq JSON data
	var rq string
	var apiName string
	switch service {
	case APIServiceSEG:
		rq = text
		apiName = APIServiceSEG
	case APIServiceNLI:
		rq = fmt.Sprintf(`{"data_type":"stt","data":{"input_type":"1","text":"%s"}}`, text)
		apiName = APIServiceNLI
	}

	// Assemble all the HTTP parameters you want to send
	body := bytes.NewBufferString(fmt.Sprintf("api=%s&appkey=%s&timestamp=%d&sign=%s&rq=%s",
		apiName, client.appKey, timeStamp, sign, rq,
	))

	req, err := http.NewRequest("POST", client.url(), body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if err != nil {
		return nil, err
	}
	return client.do(ctx, req)
}
