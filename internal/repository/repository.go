package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/ecoderat/dispatch-go/internal/model"
	"github.com/sirupsen/logrus"
)

//go:generate mockery --name=MessageRepository --output=../../mock/repository --outpkg=mockrepository --case=underscore --with-expecter
type MessageRepository interface {
	Create(ctx context.Context, message model.Message) error
	Update(ctx context.Context, id int, status model.MessageStatus) error
	Delete(ctx context.Context, id int) error
	GetAll(ctx context.Context, status ...model.MessageStatus) ([]model.Message, error)
}

type messageRepository struct {
	db     *gorm.DB
	logger *logrus.Logger
}

func NewMessageRepository(db *gorm.DB, logger *logrus.Logger) MessageRepository {
	return &messageRepository{
		db:     db,
		logger: logger,
	}
}

func (r *messageRepository) Create(ctx context.Context, message model.Message) error {
	return r.db.Create(&message).Error
}

func (r *messageRepository) Update(ctx context.Context, id int, status model.MessageStatus) error {
	return r.db.WithContext(ctx).
		Model(&model.Message{}).
		Where("id = ?", id).
		Update("status", status).
		Error
}

func (r *messageRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&model.Message{}).
		Error
}

func (r *messageRepository) GetAll(ctx context.Context, status ...model.MessageStatus) ([]model.Message, error) {
	var messages []model.Message
	query := r.db.WithContext(ctx)

	if len(status) > 0 {
		query = query.Where("status IN ?", status)
	}

	err := query.Find(&messages).Error
	if err != nil {
		return nil, err
	}

	return messages, nil
}
