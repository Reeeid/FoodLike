package usecase

import (
	"context"

	"foodlike-backend/internal/domain/model"
	"foodlike-backend/internal/domain/service"
	"foodlike-backend/internal/port/gateway"
	"foodlike-backend/internal/port/repository"
)

// SuggestionUsecase は外食先提案のコアフロー。
//
//	好き嫌い取得 → 集約(個人情報を潰す) → 外部API検索 → 絞り込み
//	→ 全員OKが0件なら妥協案ランキング
//
// 個人のPreferenceがこのusecaseより外(handler/DTO)に出ることはない。
type SuggestionUsecase struct {
	groups     repository.GroupRepository
	prefs      repository.PreferenceRepository
	restGw     gateway.RestaurantGateway
	aggregator service.PreferenceAggregator
	filter     service.RestaurantFilter
	ranker     service.CompromiseRanker
}

func NewSuggestionUsecase(
	groups repository.GroupRepository,
	prefs repository.PreferenceRepository,
	restGw gateway.RestaurantGateway,
	aggregator service.PreferenceAggregator,
	filter service.RestaurantFilter,
	ranker service.CompromiseRanker,
) *SuggestionUsecase {
	return &SuggestionUsecase{
		groups:     groups,
		prefs:      prefs,
		restGw:     restGw,
		aggregator: aggregator,
		filter:     filter,
		ranker:     ranker,
	}
}

// Suggest は提案を返す。認可(BOLA対策)として、呼び出し主が対象グループの
// メンバーであることを先に検査し、非メンバーには ErrNotMember を返す。
func (u *SuggestionUsecase) Suggest(ctx context.Context, memberID, groupID uint, area string) ([]model.Candidate, error) {
	g, err := u.groups.FindByID(ctx, groupID)
	if err != nil {
		return nil, err
	}
	if !g.HasMember(memberID) {
		return nil, ErrNotMember
	}

	groupPrefs, err := u.prefs.ListByGroup(ctx, groupID)
	if err != nil {
		return nil, err
	}

	// ここで個人との対応関係が失われる。以降は集合の制約条件だけを扱う。
	constraints := u.aggregator.Aggregate(groupPrefs)

	restaurants, err := u.restGw.Search(ctx, area)
	if err != nil {
		return nil, err
	}

	matched := u.filter.Filter(restaurants, constraints)
	if len(matched) > 0 {
		candidates := make([]model.Candidate, 0, len(matched))
		for _, r := range matched {
			candidates = append(candidates, model.Candidate{Restaurant: r, MatchedAll: true})
		}
		return candidates, nil
	}

	// 全員OKな店が0件 → 妥協案ランキング
	return u.ranker.Rank(restaurants, constraints), nil
}
