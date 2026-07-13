package repository

import "foodlike-backend/internal/domain/model"

// GORM用の永続化モデル。ドメインモデルとは分離し、DBの都合(タグ・リレーション)を
// この層に閉じ込める。層=責務の単位、GORMモデル=データの形。

type MemberRecord struct {
	ID          uint   `gorm:"primaryKey"`
	FirebaseUID string `gorm:"size:64;uniqueIndex;not null"`
	Name        string `gorm:"size:64;not null"`
}

func (MemberRecord) TableName() string { return "members" }

type GroupRecord struct {
	ID         uint           `gorm:"primaryKey"`
	Name       string         `gorm:"size:64;not null"`
	InviteCode string         `gorm:"size:32;uniqueIndex;not null"`
	Members    []MemberRecord `gorm:"many2many:group_members"`
}

func (GroupRecord) TableName() string { return "groups" }

type PreferenceRecord struct {
	ID       uint   `gorm:"primaryKey"`
	MemberID uint   `gorm:"index;not null"`
	Kind     string `gorm:"size:16;not null"`
	Category string `gorm:"size:16;not null"`
	Value    string `gorm:"size:64;not null"`
}

func (PreferenceRecord) TableName() string { return "preferences" }

func (r MemberRecord) toDomain() model.Member {
	return model.Member{ID: r.ID, Name: r.Name}
}

func (r GroupRecord) toDomain() model.Group {
	members := make([]model.Member, 0, len(r.Members))
	for _, m := range r.Members {
		members = append(members, m.toDomain())
	}
	return model.Group{ID: r.ID, Name: r.Name, InviteCode: r.InviteCode, Members: members}
}

func (r PreferenceRecord) toDomain() model.Preference {
	return model.Preference{
		ID:       r.ID,
		MemberID: r.MemberID,
		Kind:     model.PreferenceKind(r.Kind),
		Category: model.PreferenceCategory(r.Category),
		Value:    r.Value,
	}
}
