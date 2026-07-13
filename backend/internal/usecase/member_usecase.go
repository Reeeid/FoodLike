package usecase

import (
	"context"

	"foodlike-backend/internal/domain/model"
	"foodlike-backend/internal/port/repository"
)

type MemberUsecase struct {
	members repository.MemberRepository
}

func NewMemberUsecase(members repository.MemberRepository) *MemberUsecase {
	return &MemberUsecase{members: members}
}

func (u *MemberUsecase) Register(ctx context.Context, name string) (model.Member, error) {
	return u.members.Create(ctx, name)
}

func (u *MemberUsecase) Get(ctx context.Context, id uint) (model.Member, error) {
	return u.members.FindByID(ctx, id)
}

// Authenticate は検証済みのFirebase UIDから内部メンバーを解決する。
// 初回ログイン時はその場で作成する(JITプロビジョニング)。
func (u *MemberUsecase) Authenticate(ctx context.Context, firebaseUID, name string) (model.Member, error) {
	if name == "" {
		name = "名無し"
	}
	return u.members.FindOrCreateByFirebaseUID(ctx, firebaseUID, name)
}
