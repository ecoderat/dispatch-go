package scheduler

import (
	"context"
	"errors"
	"time"

	"github.com/ecoderat/dispatch-go/internal/model"
	"github.com/ecoderat/dispatch-go/internal/service/message"
	"github.com/sirupsen/logrus"
)

const (
	defaultTickerDuration = 20 * time.Second
)

var (
	ErrSchedulerStop       = errors.New("scheduler: failed to stop previous instance")
	ErrProcessMessages     = errors.New("scheduler: failed to process messages")
	ErrSendMessage         = errors.New("scheduler: failed to send message")
	ErrUpdateMessageStatus = errors.New("scheduler: failed to update message status")
)

type Scheduler interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type scheduler struct {
	ctx     context.Context
	cancel  context.CancelFunc
	ticker  *time.Ticker
	running bool

	messageService message.Service
	logger         *logrus.Logger
}

func New(messageService message.Service, logger *logrus.Logger) Scheduler {
	ctx, cancel := context.WithCancel(context.Background())
	return &scheduler{
		messageService: messageService,
		ctx:            ctx,
		cancel:         cancel,
		ticker:         time.NewTicker(defaultTickerDuration),
		running:        false,
		logger:         logger,
	}
}

func (s *scheduler) Start(ctx context.Context) error {
	err := s.Stop(ctx) // Stop any previous instance before starting a new one
	if err != nil {
		s.logger.WithError(err).Error(ErrSchedulerStop)
	}

	s.ctx, s.cancel = context.WithCancel(context.Background())
	s.ticker = time.NewTicker(defaultTickerDuration)
	s.running = true
	errChan := make(chan error)
	// Immediately process messages on start
	go func() {
		if err := s.processMessages(); err != nil {
			errChan <- err
		}
		for {
			select {
			case <-s.ticker.C:
				if err := s.processMessages(); err != nil {
					errChan <- err
				}
			case <-ctx.Done():
				s.running = false
				close(errChan)
				return
			case <-s.ctx.Done():
				s.running = false
				close(errChan)
				return
			}
		}
	}()

	// Error logging goroutine
	go func() {
		for err := range errChan {
			s.logger.WithError(err).Error(ErrProcessMessages)
		}
	}()

	return nil
}

func (s *scheduler) Stop(ctx context.Context) error {
	if s.running {
		s.cancel()
		s.ticker.Stop()
		s.running = false
	}

	return nil
}

func (s *scheduler) processMessages() error {
	messages, err := s.messageService.GetUnsentMessages(context.TODO())
	if err != nil {
		s.logger.WithError(err).Error(ErrProcessMessages)
		return ErrProcessMessages
	}

	for _, msg := range messages {
		err := s.messageService.SendMessage(context.TODO(), message.MessageRequest{
			Recipient: msg.Recipient,
			Content:   msg.Content,
		})
		if err != nil {
			s.logger.WithFields(logrus.Fields{"recipient": msg.Recipient, "id": msg.ID}).WithError(err).Error(ErrSendMessage)
			err = s.messageService.UpdateMessage(context.TODO(), msg.ID, model.StatusFailed)
			if err != nil {
				s.logger.WithFields(logrus.Fields{"id": msg.ID}).WithError(err).Error(ErrUpdateMessageStatus)
				continue
			}
			continue
		}

		s.logger.WithFields(logrus.Fields{"recipient": msg.Recipient, "id": msg.ID}).Info("Message sent successfully")

		err = s.messageService.UpdateMessage(context.TODO(), msg.ID, model.StatusSent)
		if err != nil {
			s.logger.WithFields(logrus.Fields{"id": msg.ID}).WithError(err).Error(ErrUpdateMessageStatus)
			continue
		}

		s.logger.WithFields(logrus.Fields{"id": msg.ID}).Info("Message status updated to sent")
	}

	return nil
}
