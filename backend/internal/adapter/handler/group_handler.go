package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"foodlike-backend/internal/usecase"
)

type GroupHandler struct {
	groups *usecase.GroupUsecase
}

func NewGroupHandler(groups *usecase.GroupUsecase) *GroupHandler {
	return &GroupHandler{groups: groups}
}

type createGroupRequest struct {
	Name string `json:"name" binding:"required,max=64"`
}

type joinGroupRequest struct {
	InviteCode string `json:"invite_code" binding:"required,max=32"`
}

// Create POST /api/groups
func (h *GroupHandler) Create(c *gin.Context) {
	var req createGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	g, err := h.groups.Create(c.Request.Context(), req.Name, currentMemberID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create group"})
		return
	}
	c.JSON(http.StatusCreated, toGroupResponse(g))
}

// Join POST /api/groups/join
func (h *GroupHandler) Join(c *gin.Context) {
	var req joinGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	g, err := h.groups.Join(c.Request.Context(), req.InviteCode, currentMemberID(c))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "invalid invite code"})
		return
	}
	c.JSON(http.StatusOK, toGroupResponse(g))
}

// List GET /api/groups (自分が所属するグループ一覧)
func (h *GroupHandler) List(c *gin.Context) {
	groups, err := h.groups.ListByMember(c.Request.Context(), currentMemberID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list groups"})
		return
	}
	res := make([]groupResponse, 0, len(groups))
	for _, g := range groups {
		res = append(res, toGroupResponse(g))
	}
	c.JSON(http.StatusOK, res)
}

// Get GET /api/groups/:id
func (h *GroupHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group id"})
		return
	}
	g, err := h.groups.Get(c.Request.Context(), currentMemberID(c), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "group not found"})
		return
	}
	c.JSON(http.StatusOK, toGroupResponse(g))
}
