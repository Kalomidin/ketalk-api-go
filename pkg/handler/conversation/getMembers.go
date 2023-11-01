package conversation_handler

import (
	conversation_manager "ketalk-api/pkg/manager/conversation"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GetMembersResponse struct {
	Members []Member `json:"members"`
}

type Member struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Avatar     *string   `json:"avatar"`
	LastSeenAt int64     `json:"lastSeenAt"`
}

func (h *HttpHandler) GetMembers(ctx *gin.Context, req *http.Request) (interface{}, error) {
	resp, err := h.handler.GetMembers(ctx)
	return resp, err
}

func (h *handler) GetMembers(ctx *gin.Context) (*GetMembersResponse, error) {
	conversationId, err := uuid.Parse(ctx.Param("conversationId"))
	if err != nil {
		return nil, err
	}
	var request = conversation_manager.GetMembersRequest{
		ConversationID: conversationId,
	}
	resp, err := h.manager.GetMembers(ctx, request)
	if err != nil {
		return nil, err
	}
	var members []Member = make([]Member, len(resp))
	for i, member := range resp {
		members[i] = Member{
			ID:         member.ID,
			Name:       member.Name,
			Avatar:     member.Avatar,
			LastSeenAt: member.LastSeenAt.UTC().Unix(),
		}
	}
	return &GetMembersResponse{
		Members: members,
	}, nil
}
