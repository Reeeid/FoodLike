package usecase

import (
	"context"

	"foodlike-backend/internal/domain/model"
	"foodlike-backend/internal/port/repository"
	"foodlike-backend/pkg/invitecode"
)

type GroupUsecase struct {
	groups repository.GroupRepository
}

func NewGroupUsecase(groups repository.GroupRepository) *GroupUsecase {
	return &GroupUsecase{groups: groups}
}

func (u *GroupUsecase) Create(ctx context.Context, name string, ownerID uint) (model.Group, error) {
	code, err := invitecode.New()
	if err != nil {
		return model.Group{}, err
	}
	return u.groups.Create(ctx, name, code, ownerID)
}

// Join は招待コード(QR/ID共通)でグループに参加する。
func (u *GroupUsecase) Join(ctx context.Context, inviteCode string, memberID uint) (model.Group, error) {
	g, err := u.groups.FindByInviteCode(ctx, inviteCode)
	if err != nil {
		return model.Group{}, err
	}
	if err := u.groups.AddMember(ctx, g.ID, memberID); err != nil {
		return model.Group{}, err
	}
	return u.groups.FindByID(ctx, g.ID)
}

func (u *GroupUsecase) Get(ctx context.Context, memberID, groupID uint) (model.Group, error) {
	groups, err := u.groups.FindByID(ctx, groupID)
	if err != nil {
		return model.Group{}, err
	}
	if !groups.HasMember(memberID) {
		return model.Group{}, ErrNotMember
	}
	return groups, nil

}

func (u *GroupUsecase) ListByMember(ctx context.Context, memberID uint) ([]model.Group, error) {
	return u.groups.ListByMember(ctx, memberID)
}
