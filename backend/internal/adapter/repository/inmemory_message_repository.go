package repository

import (
	"context"
	"sync"
	"time"

	"foodlike-backend/internal/domain/model"
)

// InMemoryMessageRepository はチャットメッセージの仮実装。
// サーバー再起動で消える。本実装(GORM)に差し替えるまでのつなぎ。
// TODO: messagesテーブル + GORM実装に差し替える。
type InMemoryMessageRepository struct {
	mu      sync.RWMutex
	seq     uint
	byGroup map[uint][]model.Message
}

func NewInMemoryMessageRepository() *InMemoryMessageRepository {
	return &InMemoryMessageRepository{byGroup: map[uint][]model.Message{}}
}

func (r *InMemoryMessageRepository) Create(_ context.Context, m model.Message) (model.Message, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.seq++
	m.ID = r.seq
	m.CreatedAt = time.Now()
	r.byGroup[m.GroupID] = append(r.byGroup[m.GroupID], m)
	return m, nil
}

func (r *InMemoryMessageRepository) ListByGroup(_ context.Context, groupID, afterID uint) ([]model.Message, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	all := r.byGroup[groupID]
	res := make([]model.Message, 0, len(all))
	for _, m := range all {
		if m.ID > afterID {
			res = append(res, m)
		}
	}
	return res, nil
}
