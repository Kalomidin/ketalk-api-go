package google

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"ketalk-api/pkg/provider/model"
	"net"
	"net/http"
	"time"
)

type Config struct {
	ID       string `yaml:"client_id" env:"GOOGLE_CLIENT_ID" env-default:""`
	Secret   string `yaml:"client_secret" env:"CLIENT_SECRET" env-default:""`
	Audience string `yaml:"google_token_audience" env:"GOOGLE_TOKEN_AUDIENCE"`
}

type googleClient struct {
	cfg    Config
	client *http.Client
}

func NewGoogleClient(cfg Config) model.ProviderClient {
	return &googleClient{
		client: &http.Client{
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout: time.Duration(5) * time.Second,
				}).DialContext,
				DisableKeepAlives: true,
			},
		},
		cfg: cfg,
	}
}

func (c *googleClient) Get(ctx context.Context, googleToken *string, endpoint string, result interface{}) error {
	headers := c.headersApi(googleToken)
	return c.handleHTTPRequest(ctx, http.MethodGet, endpoint, nil, result, headers)
}

func (c *googleClient) Post(ctx context.Context, googleToken *string, endpoint string, request []byte, result interface{}) error {
	headers := c.headersApi(googleToken)
	return c.handleHTTPRequest(ctx, http.MethodPost, endpoint, request, result, headers)
}

func (c *googleClient) handleHTTPRequest(ctx context.Context, httpMethod string, endpoint string, httpReq []byte, expectedOutput interface{}, httpHeaders map[string]string) error {
	resp, err := c.sendHTTPRequest(ctx, httpMethod, endpoint, httpReq, httpHeaders)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected status code response: %+v", resp.StatusCode)
	}
	if expectedOutput == nil {
		return nil
	}

	if resp.ContentLength == 0 {
		return fmt.Errorf("returned empty response")
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(expectedOutput); err != nil {
		return err
	}
	return nil
}

func (c *googleClient) sendHTTPRequest(ctx context.Context, httpMethod string, endpoint string, httpReq []byte, httpHeaders map[string]string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, httpMethod, endpoint, bytes.NewBuffer(httpReq))
	if err != nil {
		return nil, err
	}

	for k, v := range httpHeaders {
		req.Header.Add(k, v)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *googleClient) headersApi(token *string) map[string]string {
	if token == nil {
		return map[string]string{
			"Content-Type": "application/json",
		}
	}
	return map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", *token),
		"Content-Type":  "application/json",
	}
}
