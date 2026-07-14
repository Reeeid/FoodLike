package repository

import (
	"context"

	"foodlike-backend/internal/domain/model"
)

// MessageRepository はグループチャットのメッセージ永続化ポート。
type MessageRepository interface {
	Create(ctx context.Context, m model.Message) (model.Message, error)
	// ListByGroup はafterID(0なら先頭)より新しいメッセージをID昇順で返す。
	ListByGroup(ctx context.Context, groupID, afterID uint) ([]model.Message, error)
}
