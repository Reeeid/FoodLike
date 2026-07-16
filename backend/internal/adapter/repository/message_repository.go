package repository

import (
	"context"
	"foodlike-backend/internal/domain/model"
	port "foodlike-backend/internal/port/repository"

	"gorm.io/gorm"
)

type messageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) port.MessageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) Create(ctx context.Context, m model.Message) (model.Message, error) {
	rec := MessageRecord{
		GroupID:    m.GroupID,
		Role:       string(m.Role),
		MemberID:   m.MemberID,
		MemberName: m.MemberName,
		Text:       m.Text,
	}
	err := r.db.WithContext(ctx).Create(&rec).Error
	if err != nil {
		return model.Message{}, err
	}
	return rec.toDomain(), nil
}

func (r *messageRepository) ListByGroup(ctx context.Context, groupID, afterID uint) ([]model.Message, error) {
	var recs []MessageRecord
	err := r.db.WithContext(ctx).
		Where("group_id = ? AND id > ?", groupID, afterID).
		Order("id ASC").
		Find(&recs).Error
	if err != nil {
		return nil, err
	}
	messages := make([]model.Message, 0, len(recs))
	for _, rec := range recs {
		messages = append(messages, rec.toDomain())
	}
	return messages, nil
}
