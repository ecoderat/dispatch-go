package message

import (
	"context"

	"github.com/ecoderat/dispatch-go/internal/driver"
	"github.com/ecoderat/dispatch-go/internal/repository"
	"github.com/ecoderat/dispatch-go/model"
)

type Service interface {
	GetMessages(ctx context.Context) ([]model.Message, error)
	SendMessage(ctx context.Context, message MessageRequest) error
}

type service struct {
	repository repository.MessageRepository
	driver     driver.MessageDriver
}

func New(repo repository.MessageRepository, driver driver.MessageDriver) Service {
	return &service{
		repository: repo,
		driver:     driver,
	}
}

func (s *service) GetMessages(ctx context.Context) ([]model.Message, error) {
	messages, err := s.repository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

type MessageRequest struct {
	Recipient string `json:"recipient"`
	Content   string `json:"content"`
}

func (s *service) SendMessage(ctx context.Context, message MessageRequest) error {
	req := driver.MessageRequest{
		Recipient: message.Recipient,
		Content:   message.Content,
	}

	_, err := s.driver.Send(ctx, req)
	if err != nil {
		return err
	}

	return nil
}
