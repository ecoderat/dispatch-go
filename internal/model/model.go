package model

import (
	"time"

	"gorm.io/gorm"
)

type MessageStatus string

const (
	StatusSent    MessageStatus = "sent"
	StatusFailed  MessageStatus = "failed"
	StatusPending MessageStatus = "pending"
)

type Message struct {
	ID        int           `json:"id"`
	Recipient string        `json:"recipient"`
	Content   string        `json:"content"`
	Status    MessageStatus `json:"status"`

	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`
}

// GORM uses plural table names, so we need to override the table name
// https://gorm.io/docs/conventions.html#TableName
func (Message) TableName() string {
	return "message"
}
