package handler

import (
	"go-test-assesment/internal/cat/domain"
	"go-test-assesment/internal/cat/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CatHandler struct {
	usecase *usecase.CatUsecase
}

func NewCatHandler(r *gin.Engine, uc *usecase.CatUsecase) {
	h := &CatHandler{usecase: uc}

	group := r.Group("/cats")
	{
		group.POST("", h.Create)
		group.GET("/:id", h.GetByID)
		group.PUT("/:id/salary", h.UpdateSalary)
		group.DELETE("/:id", h.Delete)
		group.GET("", h.List)
	}
}

type CatRequest struct {
	Name              string  `json:"name" binding:"required,min=2,max=50" example:"Tom"`
	YearsOfExperience int     `json:"years_of_experience" binding:"required,gte=0,lte=50" example:"3"`
	Breed             string  `json:"breed" binding:"required" example:"Siamese"`
	Salary            float64 `json:"salary" binding:"required,gte=0" example:"1200.50"`
}

// swagger:model UpdateSalaryRequest
type UpdateSalaryRequest struct {
	Salary float64 `json:"salary" binding:"required,gte=0" example:"1300.75"`
}

// swagger:model CatResponse
type CatResponse struct {
	ID                int64   `json:"id" example:"1"`
	Name              string  `json:"name" example:"Tom"`
	YearsOfExperience int     `json:"years_of_experience" example:"3"`
	Breed             string  `json:"breed" example:"Siamese"`
	Salary            float64 `json:"salary" example:"1200.50"`
}

func toCatResponse(c *domain.Cat) *CatResponse {
	return &CatResponse{
		ID:                c.ID,
		Name:              c.Name,
		YearsOfExperience: c.YearsOfExperience,
		Breed:             c.Breed,
		Salary:            c.Salary,
	}
}

// @Summary Create a new cat
// @Description Create a new cat with the provided details.
// @Tags cats
// @Accept json
// @Produce json
// @Param cat body CatRequest true "Cat details"
// @Success 201 {object} CatResponse
// @Failure 400 {object} map[string]string
// @Router /cats [post]
func (h *CatHandler) Create(c *gin.Context) {
	var req CatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}

	cat := &domain.Cat{
		Name:              req.Name,
		YearsOfExperience: req.YearsOfExperience,
		Breed:             req.Breed,
		Salary:            req.Salary,
	}

	if err := h.usecase.Create(c.Request.Context(), cat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, toCatResponse(cat))
}

// GetByID godoc
// @Summary Get a cat by ID
// @Tags cats
// @Produce json
// @Param id path int true "Cat ID"
// @Success 200 {object} CatResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /cats/{id} [get]
func (h *CatHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		c.Error(err)
		return
	}

	cat, err := h.usecase.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "cat not found"})
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, toCatResponse(cat))
}

// UpdateSalary godoc
// @Summary Update a cat's salary
// @Tags cats
// @Accept json
// @Param id path int true "Cat ID"
// @Param salary body UpdateSalaryRequest true "New salary"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Router /cats/{id}/salary [put]
func (h *CatHandler) UpdateSalary(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		c.Error(err)
		return
	}

	var req UpdateSalaryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}

	if err := h.usecase.UpdateSalary(c.Request.Context(), id, req.Salary); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}

// Delete godoc
// @Summary Delete a cat by ID
// @Tags cats
// @Param id path int true "Cat ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Router /cats/{id} [delete]
func (h *CatHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		c.Error(err)
		return
	}

	if err := h.usecase.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}

// List godoc
// @Summary List all cats
// @Tags cats
// @Produce json
// @Success 200 {array} CatResponse
// @Failure 500 {object} map[string]string
// @Router /cats [get]
func (h *CatHandler) List(c *gin.Context) {
	cats, err := h.usecase.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}

	var res []*CatResponse
	for _, cat := range cats {
		res = append(res, toCatResponse(cat))
	}

	c.JSON(http.StatusOK, res)
}
