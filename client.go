// Package flowroute provides a Go client library for the Flowroute API.
//
// The Flowroute API allows you to manage phone numbers, routes, and messaging
// services. This client provides a type-safe, idiomatic Go interface to the API.
//
// Example usage:
//
//	client := flowroute.NewClient("your-api-key", "your-api-secret")
//	ctx := context.Background()
//
//	// List available area codes
//	resp, err := client.Numbers().Purchaseable().
//	    ListAvailableAreaCodes(ctx).
//	    Limit(10).
//	    Send()
//	if err != nil {
//	    log.Fatal(err)
//	}
package flowroute

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	// DefaultBaseURL is the base URL for the Flowroute API
	DefaultBaseURL = "https://api.flowroute.com"

	// DefaultTimeout is the default HTTP client timeout
	DefaultTimeout = 30 * time.Second
)

// Client is the main Flowroute API client
type Client struct {
	apiKey     string
	apiSecret  string
	baseURL    string
	httpClient *http.Client
}

// ClientOption is a functional option for configuring the Client
type ClientOption func(*Client)

// WithBaseURL sets a custom base URL for the API
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithTimeout sets the HTTP client timeout
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.httpClient.Timeout = timeout
	}
}

// NewClient creates a new Flowroute API client
//
// Example:
//
//	client := flowroute.NewClient("api-key", "api-secret")
//
//	// With custom options:
//	client := flowroute.NewClient("api-key", "api-secret",
//	    flowroute.WithTimeout(60 * time.Second),
//	)
func NewClient(apiKey, apiSecret string, opts ...ClientOption) *Client {
	c := &Client{
		apiKey:    apiKey,
		apiSecret: apiSecret,
		baseURL:   DefaultBaseURL,
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// Numbers returns a client for the Numbers API
func (c *Client) Numbers() *NumbersClient {
	return &NumbersClient{client: c}
}

// Messages returns a client for the Messaging API
func (c *Client) Messages() *MessagesClient {
	return &MessagesClient{client: c}
}

// doRequest performs an HTTP request with authentication
func (c *Client) doRequest(ctx context.Context, method, path string, query url.Values, body interface{}) (*http.Response, error) {
	u := c.baseURL + path
	if len(query) > 0 {
		u += "?" + query.Encode()
	}

	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, u, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(c.apiKey, c.apiSecret)
	req.Header.Set("Accept", "application/vnd.api+json")
	if body != nil {
		req.Header.Set("Content-Type", "application/vnd.api+json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	return resp, nil
}

// ErrorResponse represents the top-level error response from Flowroute
type ErrorResponse struct {
	Errors []ErrorDetail `json:"errors"`
}

// ErrorDetail contains detailed error information from the Flowroute API
type ErrorDetail struct {
	Status int             `json:"status"`
	ID     string          `json:"id"`
	Title  *string         `json:"title,omitempty"`
	Detail json.RawMessage `json:"detail,omitempty"`
}

// Error implements the error interface
func (e *ErrorDetail) Error() string {
	title := ""
	if e.Title != nil {
		title = *e.Title
	}

	detail := ""
	if len(e.Detail) > 0 {
		detail = string(e.Detail)
	}

	return fmt.Sprintf("Flowroute API error: status=%d, id=%s, title=%s, detail=%s",
		e.Status, e.ID, title, detail)
}

// APIError represents an error from the Flowroute API
type APIError struct {
	StatusCode int
	Detail     *ErrorDetail
}

// Error implements the error interface
func (e *APIError) Error() string {
	if e.Detail != nil {
		return e.Detail.Error()
	}
	return fmt.Sprintf("API error: status code %d", e.StatusCode)
}

// parseErrorResponse attempts to parse an error response from the API
func parseErrorResponse(resp *http.Response) error {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &APIError{
			StatusCode: resp.StatusCode,
		}
	}

	var errResp ErrorResponse
	if err := json.Unmarshal(body, &errResp); err != nil {
		return &APIError{
			StatusCode: resp.StatusCode,
		}
	}

	if len(errResp.Errors) > 0 {
		return &APIError{
			StatusCode: resp.StatusCode,
			Detail:     &errResp.Errors[0],
		}
	}

	return &APIError{
		StatusCode: resp.StatusCode,
	}
}

// decodeResponse decodes a successful response into the target struct
func decodeResponse(resp *http.Response, target interface{}) error {
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return parseErrorResponse(resp)
	}

	if target == nil {
		return nil
	}

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}
