package driver

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

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
	}

	resp, err := driver.Send(context.Background(), MessageRequest{Recipient: "+123", Content: "hi"})
	assert.Error(t, err)
	assert.Nil(t, resp)
}
