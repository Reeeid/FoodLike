package usecase

import (
	"context"

	"foodlike-backend/internal/domain/model"
	"foodlike-backend/internal/port/gateway"
	"foodlike-backend/internal/port/repository"
)

// ChatUsecase はグループチャット。メンバー同士のやり取りに加えて、
// AI検索の質問と回答も同じタイムラインにメッセージとして保存する。
type ChatUsecase struct {
	groups   repository.GroupRepository
	messages repository.MessageRepository
	ai       gateway.AIResponder
	suggest  *SuggestionUsecase
}

func NewChatUsecase(
	groups repository.GroupRepository,
	messages repository.MessageRepository,
	ai gateway.AIResponder,
	suggest *SuggestionUsecase,
) *ChatUsecase {
	return &ChatUsecase{groups: groups, messages: messages, ai: ai, suggest: suggest}
}

// List はafterIDより新しいメッセージを返す(ポーリング用)。
func (u *ChatUsecase) List(ctx context.Context, memberID, groupID, afterID uint) ([]model.Message, error) {
	if _, err := u.memberOf(ctx, memberID, groupID); err != nil {
		return nil, err
	}
	return u.messages.ListByGroup(ctx, groupID, afterID)
}

func (u *ChatUsecase) Post(ctx context.Context, memberID, groupID uint, text string) (model.Message, error) {
	name, err := u.memberOf(ctx, memberID, groupID)
	if err != nil {
		return model.Message{}, err
	}
	return u.messages.Create(ctx, model.Message{
		GroupID:    groupID,
		Role:       model.RoleMember,
		MemberID:   memberID,
		MemberName: name,
		Text:       text,
	})
}

// AISearch は質問をメンバーの発言として保存し、候補を添えてAIに回答させる。
// 回答チャンクは生成中にonChunkへ流し、完了後に全文をAIメッセージとして保存する。
// 戻り値は保存済みの(質問, 回答)。
func (u *ChatUsecase) AISearch(
	ctx context.Context,
	memberID, groupID uint,
	query string,
	onChunk func(text string),
) (question, answer model.Message, err error) {
	name, err := u.memberOf(ctx, memberID, groupID)
	if err != nil {
		return model.Message{}, model.Message{}, err
	}

	question, err = u.messages.Create(ctx, model.Message{
		GroupID:    groupID,
		Role:       model.RoleMember,
		MemberID:   memberID,
		MemberName: name,
		Text:       "🔍 " + query,
	})
	if err != nil {
		return model.Message{}, model.Message{}, err
	}

	// AIチャット検索はGemini自身のGoogle検索に任せる。こちらは全員の苦手を
	// 集約した匿名の制約だけを渡す(外部グルメAPIは叩かない)。エリア等の解釈は
	// クエリ文としてGemini側に委ねる。
	constraints, err := u.suggest.Constraints(ctx, memberID, groupID)
	if err != nil {
		return question, model.Message{}, err
	}

	full, err := u.ai.Respond(ctx, query, constraints, onChunk)
	if err != nil {
		return question, model.Message{}, err
	}

	answer, err = u.messages.Create(ctx, model.Message{
		GroupID:  groupID,
		Role:     model.RoleAI,
		MemberID: memberID,
		Text:     full,
	})
	return question, answer, err
}

// memberOf はメンバーシップを確認し、表示名を返す。非メンバーはErrNotMember。
func (u *ChatUsecase) memberOf(ctx context.Context, memberID, groupID uint) (string, error) {
	g, err := u.groups.FindByID(ctx, groupID)
	if err != nil {
		return "", err
	}
	for _, m := range g.Members {
		if m.ID == memberID {
			return m.Name, nil
		}
	}
	return "", ErrNotMember
}
