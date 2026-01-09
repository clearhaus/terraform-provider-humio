package api

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Config holds the configuration for the Humio client
type Config struct {
	Address          *url.URL
	Token            string
	CACertificatePEM string
}

// Client is the Humio API client
type Client struct {
	config     Config
	httpClient *http.Client
}

// NewClient creates a new Humio client
func NewClient(config Config) *Client {
	transport := &http.Transport{}

	if config.CACertificatePEM != "" {
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM([]byte(config.CACertificatePEM))
		transport.TLSClientConfig = &tls.Config{
			RootCAs: caCertPool,
		}
	}

	return &Client{
		config: config,
		httpClient: &http.Client{
			Transport: transport,
		},
	}
}

// graphQLRequest represents a GraphQL request
type graphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

// graphQLResponse represents a GraphQL response
type graphQLResponse struct {
	Data   json.RawMessage `json:"data"`
	Errors []graphQLError  `json:"errors,omitempty"`
}

// graphQLError represents a GraphQL error
type graphQLError struct {
	Message    string                 `json:"message"`
	Path       []interface{}          `json:"path,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
	State      map[string]interface{} `json:"state,omitempty"`
	ErrorCode  string                 `json:"errorCode,omitempty"`
}

// Query executes a GraphQL query and unmarshals the result into the provided target
func (c *Client) Query(ctx context.Context, query string, variables map[string]interface{}, target interface{}) error {
	reqBody := graphQLRequest{
		Query:     query,
		Variables: variables,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	graphqlURL := c.config.Address.JoinPath("graphql")
	req, err := http.NewRequestWithContext(ctx, "POST", graphqlURL.String(), bytes.NewReader(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.config.Token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	var gqlResp graphQLResponse
	if err := json.Unmarshal(body, &gqlResp); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(gqlResp.Errors) > 0 {
		messages := make([]string, len(gqlResp.Errors))
		for i, e := range gqlResp.Errors {
			msg := e.Message
			// Include validation details from state (e.g., "name": "Name 'foo' is already taken.")
			if len(e.State) > 0 {
				for field, detail := range e.State {
					msg = fmt.Sprintf("%s: %s=%v", msg, field, detail)
				}
			}
			if e.ErrorCode != "" {
				msg = fmt.Sprintf("%s (errorCode: %s)", msg, e.ErrorCode)
			}
			if len(e.Extensions) > 0 {
				msg = fmt.Sprintf("%s (extensions: %v)", msg, e.Extensions)
			}
			messages[i] = msg
		}
		return fmt.Errorf("GraphQL errors: %v", messages)
	}

	if target != nil {
		if err := json.Unmarshal(gqlResp.Data, target); err != nil {
			return fmt.Errorf("failed to unmarshal data: %w", err)
		}
	}

	return nil
}

// Alerts returns the Alerts API
func (c *Client) Alerts() *Alerts {
	return &Alerts{client: c}
}

// Actions returns the Actions API
func (c *Client) Actions() *Actions {
	return &Actions{client: c}
}

// Parsers returns the Parsers API
func (c *Client) Parsers() *Parsers {
	return &Parsers{client: c}
}

// Repositories returns the Repositories API
func (c *Client) Repositories() *Repositories {
	return &Repositories{client: c}
}

// IngestTokens returns the IngestTokens API
func (c *Client) IngestTokens() *IngestTokens {
	return &IngestTokens{client: c}
}

// Users returns the Users API
func (c *Client) Users() *Users {
	return &Users{client: c}
}
