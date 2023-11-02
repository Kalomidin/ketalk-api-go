package conversation_manager

import (
	"context"
	"ketalk-api/pkg/manager/conversation/repository"
	"ketalk-api/pkg/manager/port"

	"github.com/google/uuid"
)

type conversationPort struct {
	conversationRepo repository.ConversationRepository
	messageRepo      repository.MessageRepository
	memberRepo       repository.MemberRepository
}

func NewConversationPort(conversationRepo repository.ConversationRepository, messageRepo repository.MessageRepository, memberRepo repository.MemberRepository) port.ConversationPort {
	return &conversationPort{
		conversationRepo,
		messageRepo,
		memberRepo,
	}
}

func (c *conversationPort) GetItemConversations(ctx context.Context, itemID uuid.UUID) ([]port.Conversation, error) {
	conversations, err := c.conversationRepo.GetConversations(ctx, itemID)
	if err != nil {
		return nil, err
	}
	var conversationPorts []port.Conversation = make([]port.Conversation, len(conversations))
	for i, conversation := range conversations {
		lastMessage, err := c.messageRepo.GetLastMessage(ctx, conversation.ID)
		if err != nil {
			return nil, err
		}
		repoMembers, err := c.memberRepo.GetMembers(ctx, conversation.ID)
		if err != nil {
			return nil, err
		}
		var members []port.Member = make([]port.Member, len(repoMembers))
		for i, member := range repoMembers {
			members[i] = port.Member{
				ID:       member.ID,
				MemberID: member.MemberID,
			}
		}
		conversationPorts[i] = port.Conversation{
			ID:             conversation.ID,
			LastMessagedAt: lastMessage.CreatedAt,
			Members:        members,
		}
	}
	return conversationPorts, nil
}
