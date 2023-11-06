package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type conversationRepository struct {
	*gorm.DB
}

func NewConversationRepository(db *gorm.DB) ConversationRepository {
	return &conversationRepository{
		db,
	}
}

func (r *conversationRepository) CreateConversation(ctx context.Context, conversation *Conversation) error {
	res := r.Create(conversation)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (r *conversationRepository) GetConversations(ctx context.Context, itemID uuid.UUID) ([]Conversation, error) {
	var conversations []Conversation = make([]Conversation, 0)
	resp := r.Where("item_id = ?", itemID).Find(&conversations)
	if resp.Error != nil {
		return nil, resp.Error
	}
	return conversations, nil
}

func (r *conversationRepository) GetConversation(ctx context.Context, itemID uuid.UUID, userID uuid.UUID) (*Conversation, error) {
	var conversation Conversation
	resp := r.Model(&conversation).InnerJoins(
		"INNER JOIN ketalk.member ON member.convresation_id = conversation.id",
	).Where("conversation.item_id = ? AND member.member_id = ?", itemID, userID).First(&conversation)
	if resp.Error != nil {
		return nil, resp.Error
	}
	return &conversation, nil
}

func (r *conversationRepository) Migrate() error {
	return r.AutoMigrate(&Conversation{})
}
