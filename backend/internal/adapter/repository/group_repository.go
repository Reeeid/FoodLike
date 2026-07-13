package repository

import (
	"context"

	"gorm.io/gorm"

	"foodlike-backend/internal/domain/model"
	port "foodlike-backend/internal/port/repository"
)

type groupRepository struct {
	db *gorm.DB
}

func NewGroupRepository(db *gorm.DB) port.GroupRepository {
	return &groupRepository{db: db}
}

func (r *groupRepository) Create(ctx context.Context, name, inviteCode string, ownerID uint) (model.Group, error) {
	rec := GroupRecord{Name: name, InviteCode: inviteCode}
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&rec).Error; err != nil {
			return err
		}
		return tx.Model(&rec).Association("Members").Append(&MemberRecord{ID: ownerID})
	})
	if err != nil {
		return model.Group{}, err
	}
	return r.FindByID(ctx, rec.ID)
}

func (r *groupRepository) FindByID(ctx context.Context, id uint) (model.Group, error) {
	var rec GroupRecord
	if err := r.db.WithContext(ctx).Preload("Members").First(&rec, id).Error; err != nil {
		return model.Group{}, err
	}
	return rec.toDomain(), nil
}

func (r *groupRepository) FindByInviteCode(ctx context.Context, code string) (model.Group, error) {
	var rec GroupRecord
	if err := r.db.WithContext(ctx).Where("invite_code = ?", code).First(&rec).Error; err != nil {
		return model.Group{}, err
	}
	return rec.toDomain(), nil
}

func (r *groupRepository) ListByMember(ctx context.Context, memberID uint) ([]model.Group, error) {
	var recs []GroupRecord
	err := r.db.WithContext(ctx).
		Joins("JOIN group_members gm ON gm.group_record_id = groups.id").
		Where("gm.member_record_id = ?", memberID).
		Preload("Members").
		Find(&recs).Error
	if err != nil {
		return nil, err
	}
	groups := make([]model.Group, 0, len(recs))
	for _, rec := range recs {
		groups = append(groups, rec.toDomain())
	}
	return groups, nil
}

func (r *groupRepository) AddMember(ctx context.Context, groupID, memberID uint) error {
	rec := GroupRecord{ID: groupID}
	return r.db.WithContext(ctx).Model(&rec).Association("Members").Append(&MemberRecord{ID: memberID})
}
