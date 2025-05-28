package service

import (
	"context"

	"github.com/ecoderat/dispatch-go/internal/repository"
	"github.com/ecoderat/dispatch-go/model"
)

type MessageService interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	GetMessages(ctx context.Context) ([]model.Message, error)
}

type messageService struct {
	repository repository.MessageRepository
}

func NewMessageService(repo repository.MessageRepository) MessageService {
	return &messageService{
		repository: repo,
	}
}

func (s *messageService) Start(ctx context.Context) error {
	return nil
}

func (s *messageService) Stop(ctx context.Context) error {
	return nil
}

func (s *messageService) GetMessages(ctx context.Context) ([]model.Message, error) {
	messages, err := s.repository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return messages, nil
}
