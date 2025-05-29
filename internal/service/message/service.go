package message

import (
	"context"
	"errors"

	"github.com/ecoderat/dispatch-go/internal/driver"
	"github.com/ecoderat/dispatch-go/internal/model"
	"github.com/ecoderat/dispatch-go/internal/repository"
	"github.com/sirupsen/logrus"
)

var (
	ErrGetSentMessages   = errors.New("service: failed to get sent messages")
	ErrGetUnsentMessages = errors.New("service: failed to get unsent messages")
	ErrUpdateMessage     = errors.New("service: failed to update message status")
	ErrSendMessage       = errors.New("service: failed to send message")
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
	logger     *logrus.Logger
}

func New(repo repository.MessageRepository, driver driver.MessageDriver, logger *logrus.Logger) Service {
	return &service{
		repository: repo,
		driver:     driver,
		logger:     logger,
	}
}

func (s *service) GetSentMessages(ctx context.Context) ([]model.Message, error) {
	messages, err := s.repository.GetAll(ctx, model.StatusSent)
	if err != nil {
		s.logger.WithError(err).Error(ErrGetSentMessages)
		return nil, ErrGetSentMessages
	}

	s.logger.WithField("count", len(messages)).Info("Fetched sent messages")
	return messages, nil
}

func (s *service) GetUnsentMessages(ctx context.Context) ([]model.Message, error) {
	messages, err := s.repository.GetAll(ctx, model.StatusPending, model.StatusFailed)
	if err != nil {
		s.logger.WithError(err).Error(ErrGetUnsentMessages)
		return nil, ErrGetUnsentMessages
	}

	s.logger.WithField("count", len(messages)).Info("Fetched unsent messages")
	return messages, nil
}

type MessageRequest struct {
	Recipient string `json:"recipient"`
	Content   string `json:"content"`
}

func (s *service) UpdateMessage(ctx context.Context, id int, status model.MessageStatus) error {
	err := s.repository.Update(ctx, id, status)
	if err != nil {
		s.logger.WithFields(logrus.Fields{"id": id, "status": status}).WithError(err).Error(ErrUpdateMessage)
		return ErrUpdateMessage
	}

	s.logger.WithFields(logrus.Fields{"id": id, "status": status}).Info("Message status updated")
	return nil
}

func (s *service) SendMessage(ctx context.Context, message MessageRequest) error {
	req := driver.MessageRequest{
		Recipient: message.Recipient,
		Content:   message.Content,
	}

	_, err := s.driver.Send(ctx, req)
	if err != nil {
		s.logger.WithFields(logrus.Fields{"recipient": message.Recipient}).WithError(err).Error(ErrSendMessage)
		return ErrSendMessage
	}

	s.logger.WithFields(logrus.Fields{"recipient": message.Recipient}).Info("Message sent successfully")
	return nil
}
