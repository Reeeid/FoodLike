package model

type Group struct {
	ID   uint
	Name string
	// InviteCode はQRコード/ID招待に使う参加用コード。
	// URLに埋め込んでQR化するのはフロントエンドの責務。
	InviteCode string
	Members    []Member
}

func (g Group) HasMember(memberID uint) bool {
	for _, m := range g.Members {
		if m.ID == memberID {
			return true
		}
	}
	return false
}
