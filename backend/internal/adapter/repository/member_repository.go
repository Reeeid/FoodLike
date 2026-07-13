package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"foodlike-backend/internal/domain/model"
	port "foodlike-backend/internal/port/repository"
)

type memberRepository struct {
	db *gorm.DB
}

func NewMemberRepository(db *gorm.DB) port.MemberRepository {
	return &memberRepository{db: db}
}

func (r *memberRepository) Create(ctx context.Context, name string) (model.Member, error) {
	// firebase_uidはNOT NULL+UNIQUEのため、モック認証経由の登録でも
	// 衝突しないダミーUIDを入れておく。
	rec := MemberRecord{Name: name, FirebaseUID: fmt.Sprintf("mock-%d", time.Now().UnixNano())}
	if err := r.db.WithContext(ctx).Create(&rec).Error; err != nil {
		return model.Member{}, err
	}
	return rec.toDomain(), nil
}

func (r *memberRepository) FindOrCreateByFirebaseUID(ctx context.Context, uid, name string) (model.Member, error) {
	var rec MemberRecord
	err := r.db.WithContext(ctx).Where("firebase_uid = ?", uid).First(&rec).Error
	if err == nil {
		return rec.toDomain(), nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return model.Member{}, err
	}
	rec = MemberRecord{FirebaseUID: uid, Name: name}
	if createErr := r.db.WithContext(ctx).Create(&rec).Error; createErr != nil {
		// 同一UIDの同時リクエストでUNIQUE制約に負けた場合は取り直す。
		if retryErr := r.db.WithContext(ctx).Where("firebase_uid = ?", uid).First(&rec).Error; retryErr != nil {
			return model.Member{}, createErr
		}
	}
	return rec.toDomain(), nil
}

func (r *memberRepository) FindByID(ctx context.Context, id uint) (model.Member, error) {
	var rec MemberRecord
	if err := r.db.WithContext(ctx).First(&rec, id).Error; err != nil {
		return model.Member{}, err
	}
	return rec.toDomain(), nil
}
