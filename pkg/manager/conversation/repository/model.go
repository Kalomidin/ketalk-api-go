package repository

import (
	"context"
	"ketalk-api/common"
	"time"

	"github.com/google/uuid"
)

type Conversation struct {
	ID     uuid.UUID `gorm:"primaryKey;default:gen_random_uuid()"`
	ItemID uuid.UUID
	common.CreatedUpdatedDeleted
}

type Member struct {
	ID             uuid.UUID `gorm:"primaryKey;default:gen_random_uuid()"`
	ConvresationID uuid.UUID
	MemberID       uuid.UUID
	LastJoinedAt   time.Time
	common.CreatedUpdatedDeleted
}

type Message struct {
	ID             uuid.UUID `gorm:"primaryKey;default:gen_random_uuid()"`
	ConversationID uuid.UUID
	SenderID       uuid.UUID
	Message        string
	common.CreatedUpdatedDeleted
}

type ConversationMember struct {
	Conversation Conversation `gorm:"embedded"`
	Member       Member       `gorm:"embedded"`
}

type ConversationRepository interface {
	CreateConversation(ctx context.Context, conversation *Conversation) error
	GetConversations(ctx context.Context, itemID uuid.UUID) ([]Conversation, error)
	GetConversation(ctx context.Context, itemID uuid.UUID, userID uuid.UUID) (*Conversation, error)
	Migrate() error
}

type MemberRepository interface {
	AddMembers(ctx context.Context, members []Member) error
	GetConversations(ctx context.Context, userID uuid.UUID) ([]ConversationMember, error)
	SetLastJoinedAt(ctx context.Context, conversationID, userID uuid.UUID) error
	GetMembers(ctx context.Context, conversationID uuid.UUID) ([]Member, error)
	Migrate() error
}

type MessageRepository interface {
	AddMessage(ctx context.Context, message *Message) error
	GetMessages(ctx context.Context, conversationID uuid.UUID) ([]Message, error)
	GetLastMessage(ctx context.Context, conversationID uuid.UUID) (*Message, error)
	Migrate() error
}
