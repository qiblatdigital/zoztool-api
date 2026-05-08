package threads

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const baseURL = "https://graph.threads.net/v1.0"

type Client struct {
	httpClient *http.Client
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *Client) Publish(threadsUserID, accessToken, text string) error {
	if threadsUserID == "" {
		return errors.New("threads_user_id not configured for this account")
	}

	containerID, err := c.createContainer(threadsUserID, accessToken, text)
	if err != nil {
		return fmt.Errorf("create container: %w", err)
	}

	return c.publishContainer(threadsUserID, accessToken, containerID)
}

func (c *Client) createContainer(threadsUserID, accessToken, text string) (string, error) {
	endpoint := fmt.Sprintf("%s/%s/threads", baseURL, url.PathEscape(threadsUserID))

	params := url.Values{}
	params.Set("text", text)
	params.Set("media_type", "TEXT")
	params.Set("access_token", accessToken)

	resp, err := c.httpClient.Post(endpoint, "application/x-www-form-urlencoded",
		strings.NewReader(params.Encode()))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response body: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		slog.Error("threads create container error", "status", resp.StatusCode, "body", string(body))
		return "", fmt.Errorf("threads API error: status %d", resp.StatusCode)
	}

	var result struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("parse response: %w", err)
	}
	if result.ID == "" {
		return "", errors.New("empty container id from threads API")
	}
	return result.ID, nil
}

func (c *Client) publishContainer(threadsUserID, accessToken, containerID string) error {
	endpoint := fmt.Sprintf("%s/%s/threads_publish", baseURL, url.PathEscape(threadsUserID))

	params := url.Values{}
	params.Set("creation_id", containerID)
	params.Set("access_token", accessToken)

	resp, err := c.httpClient.Post(endpoint, "application/x-www-form-urlencoded",
		strings.NewReader(params.Encode()))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read publish response body: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		slog.Error("threads publish error", "status", resp.StatusCode, "body", string(body))
		return fmt.Errorf("threads publish error: status %d", resp.StatusCode)
	}

	var result struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("parse publish response: %w", err)
	}
	if result.ID == "" {
		return errors.New("empty post id from threads publish API")
	}
	return nil
}

var _ = context.Background
