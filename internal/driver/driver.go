package driver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type MessageDriver interface {
	Send(ctx context.Context, req MessageRequest) (*MessageResponse, error)
}

type messageDriver struct {
	httpClient *http.Client
	apiURL     string
}

type MessageRequest struct {
	Recipient string `json:"to"`
	Content   string `json:"content"`
}

type MessageResponse struct {
	Message   string `json:"message"`
	MessageID string `json:"messageId"`
}

func NewMessageDriver(apiURL string) MessageDriver {
	return &messageDriver{
		httpClient: &http.Client{},
		apiURL:     apiURL,
	}
}

func (m *messageDriver) Send(ctx context.Context, req MessageRequest) (*MessageResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, m.apiURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := m.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send http request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var messageResp MessageResponse
	if err := json.Unmarshal(respBody, &messageResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &messageResp, nil
}
