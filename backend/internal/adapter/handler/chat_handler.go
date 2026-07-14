package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"foodlike-backend/internal/usecase"
)

type ChatHandler struct {
	chat *usecase.ChatUsecase
}

func NewChatHandler(chat *usecase.ChatUsecase) *ChatHandler {
	return &ChatHandler{chat: chat}
}

type postMessageRequest struct {
	Text string `json:"text" binding:"required,max=500"`
}

// List GET /api/groups/:id/messages?after_id=N (ポーリング用)
func (h *ChatHandler) List(c *gin.Context) {
	groupID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group id"})
		return
	}
	afterID, _ := strconv.ParseUint(c.DefaultQuery("after_id", "0"), 10, 64)
	msgs, err := h.chat.List(c.Request.Context(), currentMemberID(c), uint(groupID), uint(afterID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "group not found"})
		return
	}
	res := make([]messageResponse, 0, len(msgs))
	for _, m := range msgs {
		res = append(res, toMessageResponse(m))
	}
	c.JSON(http.StatusOK, res)
}

// Post POST /api/groups/:id/messages
func (h *ChatHandler) Post(c *gin.Context) {
	groupID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group id"})
		return
	}
	var req postMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	m, err := h.chat.Post(c.Request.Context(), currentMemberID(c), uint(groupID), strings.TrimSpace(req.Text))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "group not found"})
		return
	}
	c.JSON(http.StatusCreated, toMessageResponse(m))
}

// AISearch GET /api/groups/:id/ai-search?q=... (SSE)
// イベント: chunk(回答テキスト断片) → done(保存済みメッセージ) / error。
// EventSourceはAuthorizationヘッダーを積めないため、フロントはfetchのストリームで読む。
func (h *ChatHandler) AISearch(c *gin.Context) {
	groupID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group id"})
		return
	}
	q := strings.TrimSpace(c.Query("q"))
	if q == "" || len([]rune(q)) > 200 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "q must be 1-200 characters"})
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("X-Accel-Buffering", "no")
	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "streaming unsupported"})
		return
	}
	emit := func(event string, v any) {
		data, err := json.Marshal(v)
		if err != nil {
			return
		}
		fmt.Fprintf(c.Writer, "event: %s\ndata: %s\n\n", event, data)
		flusher.Flush()
	}

	question, answer, err := h.chat.AISearch(
		c.Request.Context(), currentMemberID(c), uint(groupID), q,
		func(chunk string) { emit("chunk", chunk) },
	)
	if err != nil {
		emit("error", gin.H{"error": "AI検索に失敗しました"})
		return
	}
	emit("done", gin.H{
		"question": toMessageResponse(question),
		"answer":   toMessageResponse(answer),
	})
}
