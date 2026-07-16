package repository

import (
	"time"

	"foodlike-backend/internal/domain/model"
)

// GORM用の永続化モデル。ドメインモデルとは分離し、DBの都合(タグ・リレーション)を
// この層に閉じ込める。層=責務の単位、GORMモデル=データの形。

type MemberRecord struct {
	ID          uint   `gorm:"primaryKey"`
	FirebaseUID string `gorm:"size:64;uniqueIndex;not null"`
	Name        string `gorm:"size:64;not null"`
}

func (MemberRecord) TableName() string { return "members" }

func (r MemberRecord) toDomain() model.Member {
	return model.Member{ID: r.ID, Name: r.Name}
}

type GroupRecord struct {
	ID         uint           `gorm:"primaryKey"`
	Name       string         `gorm:"size:64;not null"`
	InviteCode string         `gorm:"size:32;uniqueIndex;not null"`
	Members    []MemberRecord `gorm:"many2many:group_members"`
}

func (GroupRecord) TableName() string { return "groups" }

func (r GroupRecord) toDomain() model.Group {
	members := make([]model.Member, 0, len(r.Members))
	for _, m := range r.Members {
		members = append(members, m.toDomain())
	}
	return model.Group{ID: r.ID, Name: r.Name, InviteCode: r.InviteCode, Members: members}
}

type PreferenceRecord struct {
	ID       uint   `gorm:"primaryKey"`
	MemberID uint   `gorm:"index;not null"`
	Kind     string `gorm:"size:16;not null"`
	Category string `gorm:"size:16;not null"`
	Value    string `gorm:"size:64;not null"`
}

func (PreferenceRecord) TableName() string { return "preferences" }

func (r PreferenceRecord) toDomain() model.Preference {
	return model.Preference{
		ID:       r.ID,
		MemberID: r.MemberID,
		Kind:     model.PreferenceKind(r.Kind),
		Category: model.PreferenceCategory(r.Category),
		Value:    r.Value,
	}
}

type MessageRecord struct {
	ID       uint   `gorm:"primaryKey"`
	GroupID  uint   `gorm:"index;not null"`
	Role     string `gorm:"size:16;not null"`
	MemberID uint   `gorm:"index;not null"`
	// MemberName は発言時点の表示名スナップショット(退会・改名に依存しない)。
	// 将来read model/CQRSへリファクタする際は、この列を落として
	// ListByGroupでmembersをJOINする読み取り専用型に置き換える。
	MemberName string    `gorm:"size:64;not null"`
	Text       string    `gorm:"size:512;not null"`
	CreatedAt  time.Time // GORMが作成時刻を自動セット
}

func (MessageRecord) TableName() string { return "messages" }

func (r MessageRecord) toDomain() model.Message {
	return model.Message{
		ID:         r.ID,
		GroupID:    r.GroupID,
		Role:       model.MessageRole(r.Role),
		MemberID:   r.MemberID,
		MemberName: r.MemberName,
		Text:       r.Text,
		CreatedAt:  r.CreatedAt,
	}
}
