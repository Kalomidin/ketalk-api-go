package conversation_manager

import (
	"context"
	"encoding/json"
	"errors"
	conn_redis "ketalk-api/pkg/manager/conversation/redis"
	"ketalk-api/pkg/manager/conversation/repository"
	"ketalk-api/pkg/manager/conversation/ws"
	"ketalk-api/pkg/manager/port"
	"ketalk-api/storage"

	"gorm.io/gorm"
)

type conversationManager struct {
	conversationRepo repository.ConversationRepository
	memberRepo       repository.MemberRepository
	messageRepo      repository.MessageRepository
	itemPort         port.ItemPort
	blobStorage      storage.AzureBlobStorage
	userPort         port.UserPort
	redis            conn_redis.RedisClient
}

func NewConversationManager(
	ctx context.Context,
	conversationRepo repository.ConversationRepository,
	memberRepo repository.MemberRepository,
	messageRepo repository.MessageRepository,
	itemPort port.ItemPort,
	blobStorage storage.AzureBlobStorage,
	userPort port.UserPort,
	redisClient conn_redis.RedisClient,
) ConversationManager {
	connManager := conversationManager{
		conversationRepo,
		memberRepo,
		messageRepo,
		itemPort,
		blobStorage,
		userPort,
		redisClient,
	}

	// start a task to read from redis and write into db
	go redisClient.Handle(ctx, connManager.AddMessage, "db")

	return &connManager
}

func (c *conversationManager) CreateConversation(ctx context.Context, request CreateConversationRequest) (*CreateConversationResponse, error) {
	// 1. check if conversation already exists
	if conversation, err := c.conversationRepo.GetConversation(ctx, request.ItemID, request.UserID); err == nil {
		item, err := c.itemPort.GetItem(ctx, request.ItemID)
		if err != nil {
			return nil, err
		}
		return &CreateConversationResponse{
			ID:              conversation.ID,
			SecondaryUserID: item.OwnerID,
			ItemID:          request.ItemID,
		}, nil
	} else if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// 2. Create conversation
	conversation := repository.Conversation{
		ItemID: request.ItemID,
	}
	if err := c.conversationRepo.CreateConversation(ctx, &conversation); err != nil {
		return nil, err
	}

	item, err := c.itemPort.GetItem(ctx, request.ItemID)
	if err != nil {
		return nil, err
	}

	if item.OwnerID == request.UserID {
		return nil, errors.New("user owns the item")
	}

	// 3. Create members
	var members []repository.Member = make([]repository.Member, 2)
	members[0] = repository.Member{
		MemberID:       request.UserID,
		ConvresationID: conversation.ID,
		LastJoinedAt:   conversation.CreatedAt,
	}
	members[1] = repository.Member{
		MemberID:       item.OwnerID,
		ConvresationID: conversation.ID,
		LastJoinedAt:   conversation.CreatedAt,
	}
	if err := c.memberRepo.AddMembers(ctx, members); err != nil {
		return nil, err
	}

	return &CreateConversationResponse{
		ID:              conversation.ID,
		SecondaryUserID: item.OwnerID,
		ItemID:          request.ItemID,
	}, nil
}

func (c *conversationManager) GetConversations(ctx context.Context, request GetConversationsRequest) ([]Conversation, error) {
	conversations, err := c.memberRepo.GetConversations(ctx, request.UserID)
	if err != nil {
		return nil, err
	}

	var resp []Conversation = make([]Conversation, 0)
	for _, conversation := range conversations {
		item, err := c.itemPort.GetItem(ctx, conversation.Conversation.ItemID)
		if err != nil {
			continue
		}

		message, err := c.messageRepo.GetLastMessage(ctx, conversation.Conversation.ID)
		if err != nil || message == nil {
			continue
		}

		members, err := c.memberRepo.GetMembers(ctx, conversation.Conversation.ID)
		if err != nil {
			continue
		}
		firstMember, err := c.userPort.GetUser(ctx, members[0].MemberID)
		if err != nil {
			continue
		}
		secondMember, err := c.userPort.GetUser(ctx, members[1].MemberID)
		if err != nil {
			continue
		}

		coverImage, err := c.itemPort.GetCovertImage(ctx, item.ID)
		if err != nil {
			continue
		}

		thumbnail := c.blobStorage.GetFrontDoorUrl(coverImage)

		// get the secondary user image url
		var secondaryUserImageUrl string
		image := firstMember.Image
		if request.UserID == firstMember.ID {
			image = secondMember.Image
		}
		if image != nil {
			url := c.blobStorage.GetFrontDoorUrl(*image)
			secondaryUserImageUrl = url
		}

		resp = append(resp, Conversation{
			ID:                    conversation.Conversation.ID,
			Title:                 item.Title,
			LastMessage:           message.Message,
			LastMessageAt:         message.CreatedAt,
			LastMessageSenderID:   message.SenderID,
			SecondaryUserImageUrl: secondaryUserImageUrl,
			ItemID:                item.ID,
			ItemThumbnail:         thumbnail,
			IsMessageRead:         conversation.Member.LastJoinedAt.After(message.CreatedAt),
		})
	}
	return resp, nil
}

func (c *conversationManager) AddMessage(ctx context.Context, payload string) error {
	var mes ws.Message
	if err := json.Unmarshal([]byte(payload), &mes); err != nil {
		return err
	}

	switch mes.Type {
	case ws.MessageTypeLeave:
		return nil
	case ws.MessageTypeRead:
		if err := c.memberRepo.SetLastJoinedAt(ctx, mes.ConversationID, mes.UserID); err != nil {
			return err
		}
		return nil
	case ws.MessageTypeMessage:
		message := repository.Message{
			ConversationID: mes.ConversationID,
			SenderID:       mes.UserID,
			Message:        mes.Message,
		}
		if err := c.messageRepo.AddMessage(ctx, &message); err != nil {
			return err
		}
		return nil
	default:
		return errors.New("invalid message type")
	}
}

func (c *conversationManager) GetMembers(ctx context.Context, request GetMembersRequest) ([]Member, error) {
	members, err := c.memberRepo.GetMembers(ctx, request.ConversationID)
	if err != nil {
		return nil, err
	}

	var resp []Member = make([]Member, len(members))
	for i, member := range members {
		user, err := c.userPort.GetUser(ctx, member.MemberID)
		if err != nil {
			continue
		}
		var url *string
		if user.Image != nil {
			imageUrl := c.blobStorage.GetFrontDoorUrl(*user.Image)
			url = &imageUrl
		}
		resp[i] = Member{
			ID:         user.ID,
			Name:       user.Username,
			Avatar:     url,
			LastSeenAt: member.LastJoinedAt,
		}
	}
	return resp, nil
}
