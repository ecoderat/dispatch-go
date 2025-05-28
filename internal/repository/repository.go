package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/ecoderat/dispatch-go/model"
)

type MessageRepository interface {
	Create(ctx context.Context, message string) error
	Update(ctx context.Context, id int, status model.MessageStatus) error
	Delete(ctx context.Context, id int) error
	GetAll(ctx context.Context) ([]model.Message, error)
}

type messageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) MessageRepository {
	return &messageRepository{
		db: db,
	}
}

func (r *messageRepository) Create(ctx context.Context, message string) error {
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

func (r *messageRepository) GetAll(ctx context.Context) ([]model.Message, error) {
	var messages []model.Message

	err := r.db.WithContext(ctx).
		Find(&messages).
		Error
	if err != nil {
		return nil, err
	}

	return messages, nil
}
