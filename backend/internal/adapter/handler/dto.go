package handler

import (
	"time"

	"foodlike-backend/internal/domain/model"
)

// DTO層。ドメインモデルをそのままJSONにせず、外に出す形をここで明示的に決める。
// 特にPreferenceは本人向けレスポンス以外に絶対に含めないこと。

type memberResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type groupResponse struct {
	ID         uint             `json:"id"`
	Name       string           `json:"name"`
	InviteCode string           `json:"invite_code"`
	Members    []memberResponse `json:"members"`
}

type preferenceItem struct {
	Kind     string `json:"kind" binding:"required,oneof=like dislike"`
	Category string `json:"category" binding:"required,oneof=genre ingredient"`
	Value    string `json:"value" binding:"required,max=64"`
}

type candidateResponse struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Area           string   `json:"area"`
	Genres         []string `json:"genres"`
	Budget         int      `json:"budget"`
	MatchedAll     bool     `json:"matched_all"`
	ViolationCount int      `json:"violation_count"`
}

type messageResponse struct {
	ID         uint   `json:"id"`
	Role       string `json:"role"`
	MemberID   uint   `json:"member_id"`
	MemberName string `json:"member_name"`
	Text       string `json:"text"`
	CreatedAt  string `json:"created_at"`
}

func toMessageResponse(m model.Message) messageResponse {
	return messageResponse{
		ID:         m.ID,
		Role:       string(m.Role),
		MemberID:   m.MemberID,
		MemberName: m.MemberName,
		Text:       m.Text,
		CreatedAt:  m.CreatedAt.Format(time.RFC3339),
	}
}

func toMemberResponse(m model.Member) memberResponse {
	return memberResponse{ID: m.ID, Name: m.Name}
}

func toGroupResponse(g model.Group) groupResponse {
	members := make([]memberResponse, 0, len(g.Members))
	for _, m := range g.Members {
		members = append(members, toMemberResponse(m))
	}
	return groupResponse{ID: g.ID, Name: g.Name, InviteCode: g.InviteCode, Members: members}
}

func toCandidateResponses(cs []model.Candidate) []candidateResponse {
	res := make([]candidateResponse, 0, len(cs))
	for _, c := range cs {
		res = append(res, candidateResponse{
			ID:             c.Restaurant.ID,
			Name:           c.Restaurant.Name,
			Area:           c.Restaurant.Area,
			Genres:         c.Restaurant.Genres,
			Budget:         c.Restaurant.Budget,
			MatchedAll:     c.MatchedAll,
			ViolationCount: c.ViolationCount,
		})
	}
	return res
}
