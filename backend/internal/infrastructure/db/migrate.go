package db

import (
	"gorm.io/gorm"

	"foodlike-backend/internal/adapter/repository"
)

// Migrate はGORMのAutoMigrateでスキーマを作成する。
// MVP用の簡易マイグレーション。本番運用ではgolang-migrate等の
// バージョン管理されたマイグレーションへの移行を検討する。
func Migrate(conn *gorm.DB) error {
	return conn.AutoMigrate(
		&repository.MemberRecord{},
		&repository.GroupRecord{},
		&repository.PreferenceRecord{},
		&repository.MessageRecord{},
	)
}
