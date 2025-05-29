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

func TestService_GetUnsentMessages_Success(t *testing.T) {
	repo := mockrepo.NewMessageRepository(t)
	drv := mockdriver.NewMessageDriver(t)
	svc := New(repo, drv)

	ctx := context.Background()
	messages := []model.Message{{ID: 1, Recipient: "+123", Content: "hi", Status: "pending"}}
	repo.EXPECT().GetAll(ctx, model.StatusPending, model.StatusFailed).Return(messages, nil)

	msgs, err := svc.GetUnsentMessages(ctx)
	assert.NoError(t, err)
	assert.Equal(t, messages, msgs)
}

func TestService_GetUnsentMessages_Fails(t *testing.T) {
	repo := mockrepo.NewMessageRepository(t)
	drv := mockdriver.NewMessageDriver(t)
	svc := New(repo, drv)

	ctx := context.Background()
	repo.EXPECT().GetAll(ctx, model.StatusPending, model.StatusFailed).Return(nil, errors.New("db error"))

	msgs, err := svc.GetUnsentMessages(ctx)
	assert.Error(t, err)
	assert.Nil(t, msgs)
}

func TestService_SendMessage_Success(t *testing.T) {
	repo := mockrepo.NewMessageRepository(t)
	drv := mockdriver.NewMessageDriver(t)
	svc := New(repo, drv)

	ctx := context.Background()
	msgReq := MessageRequest{Recipient: "+123", Content: "hi"}
	drv.EXPECT().Send(ctx, driver.MessageRequest{Recipient: "+123", Content: "hi"}).
		Return(&driver.MessageResponse{Message: "ok", MessageID: "123"}, nil)
	err := svc.SendMessage(ctx, msgReq)
	assert.NoError(t, err)
}

func TestService_SendMessage_Fails(t *testing.T) {
	repo := mockrepo.NewMessageRepository(t)
	drv := mockdriver.NewMessageDriver(t)
	svc := New(repo, drv)

	ctx := context.Background()
	msgReq := MessageRequest{Recipient: "+123", Content: "hi"}
	drv.EXPECT().Send(ctx, driver.MessageRequest{Recipient: "+123", Content: "hi"}).
		Return(nil, errors.New("send error"))

	err := svc.SendMessage(ctx, msgReq)
	assert.Error(t, err)
}

func TestService_GetSentMessages_Success(t *testing.T) {
	repo := mockrepo.NewMessageRepository(t)
	drv := mockdriver.NewMessageDriver(t)
	svc := New(repo, drv)

	ctx := context.Background()
	messages := []model.Message{{ID: 1, Recipient: "+123", Content: "hi", Status: "sent"}}
	repo.EXPECT().GetAll(ctx, model.StatusSent).Return(messages, nil)

	msgs, err := svc.GetSentMessages(ctx)
	assert.NoError(t, err)
	assert.Equal(t, messages, msgs)
}

func TestService_GetSentMessages_Fails(t *testing.T) {
	repo := mockrepo.NewMessageRepository(t)
	drv := mockdriver.NewMessageDriver(t)
	svc := New(repo, drv)

	ctx := context.Background()
	repo.EXPECT().GetAll(ctx, model.StatusSent).Return(nil, errors.New("db error"))

	msgs, err := svc.GetSentMessages(ctx)
	assert.Error(t, err)
	assert.Nil(t, msgs)
}

func TestService_UpdateMessage_Success(t *testing.T) {
	repo := mockrepo.NewMessageRepository(t)
	drv := mockdriver.NewMessageDriver(t)
	svc := New(repo, drv)

	ctx := context.Background()
	repo.EXPECT().Update(ctx, 1, model.StatusSent).Return(nil)

	err := svc.UpdateMessage(ctx, 1, model.StatusSent)
	assert.NoError(t, err)
}

func TestService_UpdateMessage_Fails(t *testing.T) {
	repo := mockrepo.NewMessageRepository(t)
	drv := mockdriver.NewMessageDriver(t)
	svc := New(repo, drv)

	ctx := context.Background()
	repo.EXPECT().Update(ctx, 1, model.StatusSent).Return(errors.New("update error"))

	err := svc.UpdateMessage(ctx, 1, model.StatusSent)
	assert.Error(t, err)
}
