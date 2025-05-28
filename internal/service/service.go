package service

type MessageService interface {
	Start() error
	Stop() error
	GetMessages() ([]string, error)
}

type messageService struct{}

func NewMessageService() MessageService {
	return &messageService{}
}

func (s *messageService) Start() error {
	return nil
}

func (s *messageService) Stop() error {
	return nil
}

func (s *messageService) GetMessages() ([]string, error) {
	messages := []string{"Message 1", "Message 2", "Message 3"}
	return messages, nil
}
