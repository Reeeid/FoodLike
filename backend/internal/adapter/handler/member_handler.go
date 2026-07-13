package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"foodlike-backend/internal/usecase"
)

type MemberHandler struct {
	members *usecase.MemberUsecase
}

func NewMemberHandler(members *usecase.MemberUsecase) *MemberHandler {
	return &MemberHandler{members: members}
}

type registerMemberRequest struct {
	Name string `json:"name" binding:"required,max=64"`
}

// Register POST /api/members (認証不要: ここで作ったIDを以降のX-Member-IDに使う)
func (h *MemberHandler) Register(c *gin.Context) {
	var req registerMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	m, err := h.members.Register(c.Request.Context(), req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register member"})
		return
	}
	c.JSON(http.StatusCreated, toMemberResponse(m))
}

// Me GET /api/me
func (h *MemberHandler) Me(c *gin.Context) {
	m, err := h.members.Get(c.Request.Context(), currentMemberID(c))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "member not found"})
		return
	}
	c.JSON(http.StatusOK, toMemberResponse(m))
}
