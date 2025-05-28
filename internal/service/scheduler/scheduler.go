package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/ecoderat/dispatch-go/internal/service/message"
)

const (
	DefaultTickerDuration = 20 * time.Second
)

type Scheduler interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type scheduler struct {
	messageService message.Service

	ctx     context.Context
	cancel  context.CancelFunc
	ticker  *time.Ticker
	running bool
}

func New(messageService message.Service) Scheduler {
	ctx, cancel := context.WithCancel(context.Background())
	return &scheduler{
		messageService: messageService,
		ctx:            ctx,
		cancel:         cancel,
		ticker:         time.NewTicker(DefaultTickerDuration),
		running:        false,
	}
}

func (s *scheduler) Start(ctx context.Context) error {
	err := s.Stop(ctx) // Stop any previous instance before starting a new one
	if err != nil {
		log.Printf("Error stopping previous scheduler: %v", err)
	}

	s.ctx, s.cancel = context.WithCancel(context.Background())
	s.ticker = time.NewTicker(DefaultTickerDuration)
	s.running = true
	errChan := make(chan error)
	go func() {
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
			log.Printf("Error processing messages: %v", err)
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
	messages, err := s.messageService.GetMessages(context.TODO())
	if err != nil {
		return err
	}

	for _, msg := range messages {
		err := s.messageService.SendMessage(context.TODO(), message.MessageRequest{
			Recipient: msg.Recipient,
			Content:   msg.Content,
		})
		if err != nil {
			log.Printf("Failed to send message to %s: %v", msg.Recipient, err)
			continue
		}
	}

	return nil
}
