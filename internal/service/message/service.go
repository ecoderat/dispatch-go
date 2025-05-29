package message

import (
	"context"

	"github.com/ecoderat/dispatch-go/internal/driver"
	"github.com/ecoderat/dispatch-go/internal/model"
	"github.com/ecoderat/dispatch-go/internal/repository"
)

//go:generate mockery --name=Service --output=../../../mock/service/message --outpkg=mock_service_message --case=underscore --with-expecter
type Service interface {
	GetUnsentMessages(ctx context.Context) ([]model.Message, error)
	GetSentMessages(ctx context.Context) ([]model.Message, error)
	UpdateMessage(ctx context.Context, id int, status model.MessageStatus) error
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

func (s *service) GetSentMessages(ctx context.Context) ([]model.Message, error) {
	messages, err := s.repository.GetAll(ctx, model.StatusSent)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (s *service) GetUnsentMessages(ctx context.Context) ([]model.Message, error) {
	messages, err := s.repository.GetAll(ctx, model.StatusPending, model.StatusFailed)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

type MessageRequest struct {
	Recipient string `json:"recipient"`
	Content   string `json:"content"`
}

func (s *service) UpdateMessage(ctx context.Context, id int, status model.MessageStatus) error {
	err := s.repository.Update(ctx, id, status)
	if err != nil {
		return err
	}

	return nil
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
