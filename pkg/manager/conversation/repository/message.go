package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type messageRepository struct {
	*gorm.DB
}

func NewMessageRepository(db *gorm.DB) MessageRepository {
	return &messageRepository{
		db,
	}
}

func (r *messageRepository) AddMessage(ctx context.Context, message *Message) error {
	res := r.Create(message)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (r *messageRepository) GetMessages(ctx context.Context, conversationID uuid.UUID) ([]Message, error) {
	var messages []Message = make([]Message, 0)
	resp := r.Where("conversation_id = ?", conversationID).Find(&messages)
	if resp.Error != nil {
		return nil, resp.Error
	}
	return messages, nil
}

func (r *messageRepository) GetLastMessage(ctx context.Context, conversationID uuid.UUID) (*Message, error) {
	var message Message
	resp := r.Where("conversation_id = ?", conversationID).Order("created_at DESC").First(&message)
	if resp.Error != nil {
		return nil, resp.Error
	}
	return &message, nil
}

func (r *messageRepository) Migrate() error {
	return r.AutoMigrate(&Message{})
}
