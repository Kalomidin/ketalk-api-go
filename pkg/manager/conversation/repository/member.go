package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type memberRepository struct {
	*gorm.DB
}

func NewMemberRepository(db *gorm.DB) MemberRepository {
	return &memberRepository{
		db,
	}
}

func (r *memberRepository) AddMembers(ctx context.Context, members []Member) error {
	res := r.CreateInBatches(&members, len(members))
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected != int64(len(members)) {
		return fmt.Errorf("unexpected number of rows affected")
	}
	return nil
}

func (r *memberRepository) GetConversations(ctx context.Context, userID uuid.UUID) ([]ConversationMember, error) {
	var conversations []ConversationMember = make([]ConversationMember, 0)
	resp := r.Model(&Conversation{}).Select(
		"conversation.id, conversation.item_id, conversation.created_at, conversation.updated_at, conversation.deleted_at, member.id, member.convresation_id, member.member_id, member.last_joined_at, member.created_at, member.updated_at, member.deleted_at",
	).
		InnerJoins(
			"INNER JOIN member ON member.convresation_id = conversation.id",
		).Where("member.member_id = ?", userID).Find(&conversations)
	if resp.Error != nil {
		return nil, resp.Error
	}
	return conversations, nil
}

func (r *memberRepository) SetLastJoinedAt(ctx context.Context, conversationID, userID uuid.UUID) error {
	res := r.Model(&Member{}).Where("convresation_id = ? AND member_id = ?", conversationID, userID).Update("last_joined_at", "NOW()")
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected != 1 {
		return fmt.Errorf("unexpected number of rows affected")
	}
	return nil
}

func (r *memberRepository) GetMembers(ctx context.Context, conversationID uuid.UUID) ([]Member, error) {
	var members []Member = make([]Member, 0)
	resp := r.Where("convresation_id = ?", conversationID).Find(&members)
	if resp.Error != nil {
		return nil, resp.Error
	}
	return members, nil
}

func (r *memberRepository) Migrate() error {
	return r.AutoMigrate(&Member{})
}
