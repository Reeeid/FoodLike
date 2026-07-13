package repository

import (
	"context"

	"foodlike-backend/internal/domain/model"
)

type MemberRepository interface {
	Create(ctx context.Context, name string) (model.Member, error)
	FindByID(ctx context.Context, id uint) (model.Member, error)
	// FindOrCreateByFirebaseUID はFirebase UIDに対応するメンバーを返し、
	// 存在しなければ作成する(初回ログイン時のJITプロビジョニング)。
	FindOrCreateByFirebaseUID(ctx context.Context, uid, name string) (model.Member, error)
}
