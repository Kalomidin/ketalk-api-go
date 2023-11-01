package conversation_manager

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type CreateConversationRequest struct {
	UserID uuid.UUID
	ItemID uuid.UUID
}

type CreateConversationResponse struct {
	ID              uuid.UUID
	SecondaryUserID uuid.UUID
	ItemID          uuid.UUID
}

type Conversation struct {
	ID                    uuid.UUID
	Title                 string
	LastMessage           string
	LastMessageAt         time.Time
	LastMessageSenderID   uuid.UUID
	SecondaryUserImageUrl string
	ItemID                uuid.UUID
	ItemThumbnail         string
	IsMessageRead         bool
}

type GetConversationsRequest struct {
	UserID uuid.UUID
}

type GetMembersRequest struct {
	ConversationID uuid.UUID
}

type Member struct {
	ID         uuid.UUID
	Name       string
	Avatar     *string
	LastSeenAt time.Time
}

type ConversationManager interface {
	CreateConversation(ctx context.Context, req CreateConversationRequest) (*CreateConversationResponse, error)
	GetConversations(ctx context.Context, req GetConversationsRequest) ([]Conversation, error)
	AddMessage(ctx context.Context, mes string) error
	GetMembers(ctx context.Context, req GetMembersRequest) ([]Member, error)
}
