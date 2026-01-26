package http

import (
	"net/http"
	"strconv"

	"Test2/internal/domain"

	"github.com/gin-gonic/gin"
)

type CateHandler struct {
	CateUseCase domain.CategoryUseCase
}

func NewCateHandler(r *gin.Engine, us domain.CategoryUseCase) {
	handler := &CateHandler{
		CateUseCase: us,
	}

	v1 := r.Group("/api/v1")
	{
		v1.POST("/categories/add", handler.Store)
		v1.GET("/categories/list", handler.Fetch)
		v1.GET("/categories/find/:id", handler.GetByID)
		v1.PUT("/categories/update/:id", handler.Update)
		v1.DELETE("/categories/delete/:id", handler.Delete)
	}
}

func (h *CateHandler) Store(c *gin.Context) {
	var cate domain.Category

	if err := c.ShouldBindJSON(&cate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	err := h.CateUseCase.Store(ctx, &cate)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, cate)
}

func (h *CateHandler) Fetch(c *gin.Context) {
	// Lấy params page & page_size từ URL
	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	pageSize, _ := strconv.ParseInt(c.DefaultQuery("page", "10"), 10, 64)

	categories, err := h.CateUseCase.Fetch(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": categories})
}

func (h *CateHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	categories, err := h.CateUseCase.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, categories)
}

func (h *CateHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID to Update"})
		return
	}

	var cate domain.Category
	if err := c.ShouldBindJSON(&cate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cate.ID = id

	err = h.CateUseCase.Update(c.Request.Context(), &cate)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": cate})
}

func (h *CateHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	err = h.CateUseCase.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Category soft deleted successfully"})
}
