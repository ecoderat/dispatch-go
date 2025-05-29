package message

import (
	"context"
	"errors"
	"testing"

	"github.com/ecoderat/dispatch-go/internal/driver"
	"github.com/ecoderat/dispatch-go/internal/model"
	mockdriver "github.com/ecoderat/dispatch-go/mock/driver"
	mockrepo "github.com/ecoderat/dispatch-go/mock/repository"
	"github.com/stretchr/testify/assert"
)

func TestService_GetMessages_Success(t *testing.T) {
	repo := mockrepo.NewMessageRepository(t)
	drv := mockdriver.NewMessageDriver(t)
	svc := New(repo, drv)

	ctx := context.Background()
	messages := []model.Message{{ID: 1, Recipient: "+123", Content: "hi", Status: "pending"}}
	repo.On("GetAll", ctx).Return(messages, nil)

	msgs, err := svc.GetMessages(ctx)
	assert.NoError(t, err)
	assert.Equal(t, messages, msgs)
}

func TestService_GetMessages_Fails(t *testing.T) {
	repo := mockrepo.NewMessageRepository(t)
	drv := mockdriver.NewMessageDriver(t)
	svc := New(repo, drv)

	ctx := context.Background()
	repo.On("GetAll", ctx).Return(nil, errors.New("db error"))

	msgs, err := svc.GetMessages(ctx)
	assert.Error(t, err)
	assert.Nil(t, msgs)
}

func TestService_SendMessage_Success(t *testing.T) {
	repo := mockrepo.NewMessageRepository(t)
	drv := mockdriver.NewMessageDriver(t)
	svc := New(repo, drv)

	ctx := context.Background()
	msgReq := MessageRequest{Recipient: "+123", Content: "hi"}
	drv.On("Send", ctx, driver.MessageRequest{Recipient: "+123", Content: "hi"}).Return(&driver.MessageResponse{Message: "ok", MessageID: "1"}, nil)

	err := svc.SendMessage(ctx, msgReq)
	assert.NoError(t, err)
}

func TestService_SendMessage_Fails(t *testing.T) {
	repo := mockrepo.NewMessageRepository(t)
	drv := mockdriver.NewMessageDriver(t)
	svc := New(repo, drv)

	ctx := context.Background()
	msgReq := MessageRequest{Recipient: "+123", Content: "hi"}
	drv.On("Send", ctx, driver.MessageRequest{Recipient: "+123", Content: "hi"}).Return(nil, errors.New("send error"))

	err := svc.SendMessage(ctx, msgReq)
	assert.Error(t, err)
}
