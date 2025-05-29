package driver

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestMessageDriver_Send_Success(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		resp := MessageResponse{Message: "ok", MessageID: "123"}
		w.WriteHeader(http.StatusAccepted)
		_ = json.NewEncoder(w).Encode(resp)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	driver := &messageDriver{
		httpClient: server.Client(),
		apiURL:     server.URL,
		logger:     logrus.New(),
	}

	resp, err := driver.Send(context.Background(), MessageRequest{Recipient: "+123", Content: "hi"})
	assert.NoError(t, err)
	assert.Equal(t, "ok", resp.Message)
	assert.Equal(t, "123", resp.MessageID)
}

func TestMessageDriver_Send_UnmarshallError(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("not-json"))
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	driver := &messageDriver{
		httpClient: server.Client(),
		apiURL:     server.URL,
		logger:     logrus.New(),
	}

	resp, err := driver.Send(context.Background(), MessageRequest{Recipient: "+123", Content: "hi"})
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestMessageDriver_Send_HTTPError(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()
	driver := &messageDriver{
		httpClient: server.Client(),
		apiURL:     server.URL,
		logger:     logrus.New(),
	}

	resp, err := driver.Send(context.Background(), MessageRequest{Recipient: "+123", Content: "hi"})
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestMessageDriver_Send_MultipartMessage(t *testing.T) {
	// Compose a message longer than 160 chars to trigger multipart logic
	longContent := "" // 160 + 10 = 170 chars
	for i := 0; i < 170; i++ {
		longContent += "a"
	}
	var receivedParts []string
	handler := func(w http.ResponseWriter, r *http.Request) {
		var req MessageRequest
		_ = json.NewDecoder(r.Body).Decode(&req)
		receivedParts = append(receivedParts, req.Content)
		resp := MessageResponse{Message: "ok", MessageID: "part"}
		w.WriteHeader(http.StatusAccepted)
		_ = json.NewEncoder(w).Encode(resp)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	driver := &messageDriver{
		httpClient: server.Client(),
		apiURL:     server.URL,
		logger:     logrus.New(),
	}

	resp, err := driver.Send(context.Background(), MessageRequest{Recipient: "+123", Content: longContent})
	assert.NoError(t, err)
	assert.Equal(t, "ok", resp.Message)
	assert.Equal(t, "part", resp.MessageID)
	assert.True(t, len(receivedParts) > 1, "Should send multiple parts")
}
