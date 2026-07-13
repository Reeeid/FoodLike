package repository

import (
	"context"

	"foodlike-backend/internal/domain/model"
)

type GroupRepository interface {
	Create(ctx context.Context, name, inviteCode string, ownerID uint) (model.Group, error)
	FindByID(ctx context.Context, id uint) (model.Group, error)
	FindByInviteCode(ctx context.Context, code string) (model.Group, error)
	ListByMember(ctx context.Context, memberID uint) ([]model.Group, error)
	AddMember(ctx context.Context, groupID, memberID uint) error
}
