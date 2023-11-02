package port

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Conversation struct {
	ID             uuid.UUID
	LastMessagedAt time.Time
	Members        []Member
}

type Member struct {
	ID       uuid.UUID
	MemberID uuid.UUID
}

type ConversationPort interface {
	GetItemConversations(ctx context.Context, itemID uuid.UUID) ([]Conversation, error)
}
