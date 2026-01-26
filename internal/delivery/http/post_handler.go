package http

import (
	"net/http"
	"strconv"

	"Test2/internal/domain"

	"github.com/gin-gonic/gin"
)

// PostHandler hứng các request liên quan đến Post
type PostHandler struct {
	PostUseCase domain.PostUseCase
}

// NewPostHandler khởi tạo Handler và đăng ký routes
func NewPostHandler(r *gin.Engine, us domain.PostUseCase) {
	handler := &PostHandler{
		PostUseCase: us,
	}

	// Group routes api/v1
	v1 := r.Group("/api/v1")
	{
		v1.POST("/posts", handler.Store)
		v1.GET("/posts", handler.Fetch)
		v1.GET("/posts/:id", handler.GetByID)
		v1.DELETE("/posts/:id", handler.Delete)
	}
}

// Create Post
func (h *PostHandler) Store(c *gin.Context) {
	var post domain.Post
	// BindJSON giúp parse body và validate struct tag
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	err := h.PostUseCase.Store(ctx, &post)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, post)
}

// Get List Posts
func (h *PostHandler) Fetch(c *gin.Context) {
	// Lấy params page & page_size từ URL
	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	pageSize, _ := strconv.ParseInt(c.DefaultQuery("page_size", "10"), 10, 64)

	posts, err := h.PostUseCase.Fetch(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": posts})
}

// Get One Post
func (h *PostHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	post, err := h.PostUseCase.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, post)
}

// Soft Delete Post
func (h *PostHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	err = h.PostUseCase.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post soft deleted successfully"})
}
