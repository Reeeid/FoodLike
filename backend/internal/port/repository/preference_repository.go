package repository

import (
	"context"

	"foodlike-backend/internal/domain/model"
)

type PreferenceRepository interface {
	// ListByMember は本人の好き嫌い一覧。本人向けAPIでのみ使うこと。
	ListByMember(ctx context.Context, memberID uint) ([]model.Preference, error)
	// ListByGroup はグループ全員分の好き嫌い。
	// PreferenceAggregatorへの入力専用で、APIレスポンスに直接載せてはならない。
	ListByGroup(ctx context.Context, groupID uint) ([]model.Preference, error)
	// ReplaceForMember は本人の好き嫌いを全件入れ替える(MVPは差分更新なしの単純方式)。
	ReplaceForMember(ctx context.Context, memberID uint, prefs []model.Preference) error
}
