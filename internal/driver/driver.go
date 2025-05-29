package driver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
)

//go:generate mockery --name=MessageDriver --output=../../mock/driver --outpkg=mockdriver --case=underscore --with-expecter
type MessageDriver interface {
	Send(ctx context.Context, req MessageRequest) (*MessageResponse, error)
}

type messageDriver struct {
	httpClient *http.Client
	apiURL     string
	logger     *logrus.Logger
}

type MessageRequest struct {
	Recipient string `json:"to"`
	Content   string `json:"content"`
}

type MessageResponse struct {
	Message   string `json:"message"`
	MessageID string `json:"messageId"`
}

var (
	ErrMarshalRequest    = fmt.Errorf("driver: failed to marshal request")
	ErrCreateHTTPRequest = fmt.Errorf("driver: failed to create http request")
	ErrSendHTTPRequest   = fmt.Errorf("driver: failed to send http request")
	ErrUnexpectedStatus  = fmt.Errorf("driver: unexpected status code")
	ErrReadResponseBody  = fmt.Errorf("driver: failed to read response body")
	ErrUnmarshalResponse = fmt.Errorf("driver: failed to unmarshal response")
)

func NewMessageDriver(apiURL string, logger *logrus.Logger) MessageDriver {
	return &messageDriver{
		httpClient: &http.Client{},
		apiURL:     apiURL,
		logger:     logger,
	}
}

func (m *messageDriver) Send(ctx context.Context, req MessageRequest) (*MessageResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		m.logger.WithError(err).Error(ErrMarshalRequest)
		return nil, fmt.Errorf("%w: %v", ErrMarshalRequest, err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, m.apiURL, bytes.NewBuffer(body))
	if err != nil {
		m.logger.WithError(err).Error(ErrCreateHTTPRequest)
		return nil, fmt.Errorf("%w: %v", ErrCreateHTTPRequest, err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := m.httpClient.Do(httpReq)
	if err != nil {
		m.logger.WithError(err).Error(ErrSendHTTPRequest)
		return nil, fmt.Errorf("%w: %v", ErrSendHTTPRequest, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		m.logger.WithField("status_code", resp.StatusCode).Error(ErrUnexpectedStatus)
		return nil, fmt.Errorf("%w: %d", ErrUnexpectedStatus, resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		m.logger.WithError(err).Error(ErrReadResponseBody)
		return nil, fmt.Errorf("%w: %v", ErrReadResponseBody, err)
	}

	var messageResp MessageResponse
	if err := json.Unmarshal(respBody, &messageResp); err != nil {
		m.logger.WithError(err).Error(ErrUnmarshalResponse)
		return nil, fmt.Errorf("%w: %v", ErrUnmarshalResponse, err)
	}

	m.logger.WithFields(logrus.Fields{
		"recipient":  req.Recipient,
		"message_id": messageResp.MessageID,
	}).Info("Message sent successfully via driver")

	return &messageResp, nil
}
