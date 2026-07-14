package model

import "time"

// MessageRole はチャットメッセージの発言者種別。
type MessageRole string

const (
	// RoleMember はグループメンバーの発言。
	RoleMember MessageRole = "member"
	// RoleAI はAI検索の回答。
	RoleAI MessageRole = "ai"
)

// Message はグループチャットの1メッセージ。
// AI検索の回答も1メッセージとして保存し、履歴として全員が見られる。
type Message struct {
	ID         uint
	GroupID    uint
	Role       MessageRole
	MemberID   uint   // RoleAIのときは質問したメンバーのID
	MemberName string // 表示用スナップショット(退会・改名に依存しない)
	Text       string
	CreatedAt  time.Time
}
