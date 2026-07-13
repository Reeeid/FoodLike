package repository

import (
	"context"

	"gorm.io/gorm"

	"foodlike-backend/internal/domain/model"
	port "foodlike-backend/internal/port/repository"
)

type preferenceRepository struct {
	db *gorm.DB
}

func NewPreferenceRepository(db *gorm.DB) port.PreferenceRepository {
	return &preferenceRepository{db: db}
}

func (r *preferenceRepository) ListByMember(ctx context.Context, memberID uint) ([]model.Preference, error) {
	var recs []PreferenceRecord
	if err := r.db.WithContext(ctx).Where("member_id = ?", memberID).Find(&recs).Error; err != nil {
		return nil, err
	}
	return toDomainPrefs(recs), nil
}

func (r *preferenceRepository) ListByGroup(ctx context.Context, groupID uint) ([]model.Preference, error) {
	var recs []PreferenceRecord
	err := r.db.WithContext(ctx).
		Joins("JOIN group_members gm ON gm.member_record_id = preferences.member_id").
		Where("gm.group_record_id = ?", groupID).
		Find(&recs).Error
	if err != nil {
		return nil, err
	}
	return toDomainPrefs(recs), nil
}

func (r *preferenceRepository) ReplaceForMember(ctx context.Context, memberID uint, prefs []model.Preference) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("member_id = ?", memberID).Delete(&PreferenceRecord{}).Error; err != nil {
			return err
		}
		if len(prefs) == 0 {
			return nil
		}
		recs := make([]PreferenceRecord, 0, len(prefs))
		for _, p := range prefs {
			recs = append(recs, PreferenceRecord{
				MemberID: memberID,
				Kind:     string(p.Kind),
				Category: string(p.Category),
				Value:    p.Value,
			})
		}
		return tx.Create(&recs).Error
	})
}

func toDomainPrefs(recs []PreferenceRecord) []model.Preference {
	prefs := make([]model.Preference, 0, len(recs))
	for _, rec := range recs {
		prefs = append(prefs, rec.toDomain())
	}
	return prefs
}
