package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const BaseURL = "https://discord.com/api/v10"

type Client struct {
	token       string
	httpClient  *http.Client
	rateLimiter *RateLimiter
}

func NewClient(token string) *Client {
	return &Client{
		token:       token,
		httpClient:  &http.Client{},
		rateLimiter: NewRateLimiter(),
	}
}

func (c *Client) doRequest(ctx context.Context, method, endpoint string, body any, v any) error {
	url := fmt.Sprintf("%s%s", BaseURL, endpoint)

	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return err
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bot %s", c.token))
	req.Header.Set("Content-Type", "application/json")

	if err := c.rateLimiter.Wait(ctx, endpoint, method); err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("discord API error: %s (%s)", resp.Status, string(respBody))
	}

	if v != nil {
		if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) CreateMessage(ctx context.Context, channelID string, params CreateMessageParams) (*Message, error) {
	var msg Message
	endpoint := fmt.Sprintf("/channels/%s/messages", channelID)
	err := c.doRequest(ctx, http.MethodPost, endpoint, params, &msg)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}
