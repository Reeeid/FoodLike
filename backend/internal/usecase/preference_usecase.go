package usecase

import (
	"context"

	"foodlike-backend/internal/domain/model"
	"foodlike-backend/internal/port/repository"
)

type PreferenceUsecase struct {
	prefs repository.PreferenceRepository
}

func NewPreferenceUsecase(prefs repository.PreferenceRepository) *PreferenceUsecase {
	return &PreferenceUsecase{prefs: prefs}
}

// ListOwn は本人の好き嫌い一覧。呼び出し元(handler)は認証済みの
// 本人のmemberIDのみを渡すこと(他人のIDを渡せる口を作らない)。
func (u *PreferenceUsecase) ListOwn(ctx context.Context, memberID uint) ([]model.Preference, error) {
	return u.prefs.ListByMember(ctx, memberID)
}

func (u *PreferenceUsecase) Replace(ctx context.Context, memberID uint, prefs []model.Preference) error {
	return u.prefs.ReplaceForMember(ctx, memberID, prefs)
}
