package handler

import (
	"net/http"
	"strconv"

	"go-test-assesment/internal/mission/domain"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	usecase domain.Usecase
}

type MissionDTO struct {
	CatID     *int64 `json:"cat_id,omitempty" example:"123"`
	Completed bool   `json:"completed" example:"false"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type TargetDTO struct {
	Name      string `json:"name" example:"Target name"`
	Country   string `json:"country" example:"Country name"`
	Notes     string `json:"notes,omitempty" example:"Additional notes"`
	Completed bool   `json:"completed" example:"false"`
}

func NewHandler(u domain.Usecase) *Handler {
	return &Handler{usecase: u}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	missions := r.Group("/missions")
	{
		missions.POST("", h.createMission)
		missions.GET("", h.listMissions)
		missions.GET("/:id", h.getMissionByID)
		missions.PUT("/:id", h.updateMission)
		missions.DELETE("/:id", h.deleteMission)
		missions.POST("/:id/cat/:catID", h.assignCatToMission)
		missions.POST("/:id/targets", h.addTargets)
	}

	targets := r.Group("/targets")
	{
		targets.PUT("/:id", h.updateTarget)
		targets.DELETE("/:id", h.deleteTarget)
	}
}

// createMission godoc
// @Summary Create a new mission
// @Description Create a new mission with the provided details.
// @Tags Missions
// @Accept json
// @Produce json
// @Param mission body MissionDTO true "Mission details"
// @Success 201 {object} domain.Mission
// @Failure 400 {object} ErrorResponse "Wrong request format"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /missions [post]
func (h *Handler) createMission(c *gin.Context) {
	var missionDTO MissionDTO
	if err := c.ShouldBindJSON(&missionDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}

	mission := domain.Mission{
		CatID:     missionDTO.CatID,
		Completed: missionDTO.Completed,
	}

	if err := h.usecase.CreateMission(c.Request.Context(), &mission); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, mission)
}

// getMissionByID godoc
// @Summary Get a mission by ID
// @Description Retrieve a mission by its ID.
// @Tags Missions
// @Produce json
// @Param id path int true "Mission ID"
// @Success 200 {object} domain.Mission
// @Failure 400 {object} ErrorResponse "Invalid mission ID"
// @Failure 404 {object} ErrorResponse "Mission not found"
// @Router /missions/{id} [get]
func (h *Handler) getMissionByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid mission id"})
		c.Error(err)
		return
	}
	mission, err := h.usecase.GetMissionByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, mission)
}

// listMissions godoc
// @Summary List all missions
// @Description Retrieve a list of all missions.
// @Tags Missions
// @Produce json
// @Success 200 {array} domain.Mission
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /missions [get]
func (h *Handler) listMissions(c *gin.Context) {
	missions, err := h.usecase.ListMissions(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, missions)
}

// updateMission godoc
// @Summary Update a mission
// @Description Update an existing mission with the provided details.
// @Tags Missions
// @Accept json
// @Produce json
// @Param id path int true "Mission ID"
// @Param mission body MissionDTO true "Mission details"
// @Success 200 {object} domain.Mission
// @Failure 400 {object} ErrorResponse "Invalid mission ID or request format"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /missions/{id} [put]
func (h *Handler) updateMission(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid mission id"})
		c.Error(err)
		return
	}

	var dto MissionDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}

	mission := domain.Mission{
		ID:        id,
		CatID:     dto.CatID,
		Completed: dto.Completed,
	}

	if err := c.ShouldBindJSON(&mission); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}
	mission.ID = id
	if err := h.usecase.UpdateMission(c.Request.Context(), &mission); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, mission)
}

// deleteMission godoc
// @Summary Delete a mission
// @Description Delete a mission by its ID.
// @Tags Missions
// @Param id path int true "Mission ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse "Invalid mission ID"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /missions/{id} [delete]
func (h *Handler) deleteMission(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid mission id"})
		c.Error(err)
		return
	}
	if err := h.usecase.DeleteMission(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}

// @Summary Assign Cat to Mission
// @Description Assign a cat to a mission by their IDs.
// @Tags Missions
// @Param id path int true "Mission ID"
// @Param catID path int true "Cat ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse "Wrong request format"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /missions/{id}/cat/{catID} [post]
func (h *Handler) assignCatToMission(c *gin.Context) {
	missionID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid mission id"})
		c.Error(err)
		return
	}
	catID, err := strconv.ParseInt(c.Param("catID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid cat id"})
		c.Error(err)
		return
	}
	if err := h.usecase.AssignCatToMission(c.Request.Context(), missionID, catID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}

// addTargets godoc
// @Summary Add Targets to Mission
// @Description Add multiple targets to a mission by its ID.
// @Tags Missions
// @Accept json
// @Param id path int true "Mission ID"
// @Param targets body []TargetDTO true "Targets to add"
// @Success 201 "Created"
// @Failure 400 {object} ErrorResponse "Invalid request format"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /missions/{id}/targets [post]
func (h *Handler) addTargets(c *gin.Context) {
	missionID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid mission id"})
		c.Error(err)
		return
	}

	var targetsDTO []TargetDTO
	if err := c.ShouldBindJSON(&targetsDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}

	domainTargets := make([]domain.Target, 0, len(targetsDTO))
	for _, t := range targetsDTO {
		domainTargets = append(domainTargets, domain.Target{
			MissionID: missionID,
			Name:      t.Name,
			Country:   t.Country,
			Notes:     t.Notes,
			Completed: t.Completed,
		})
	}

	if err := h.usecase.AddTargets(c.Request.Context(), missionID, domainTargets); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}

	c.Status(http.StatusCreated)
}

// updateTarget godoc
// @Summary Update Target
// @Description Update an existing target by its ID.
// @Tags Targets
// @Accept json
// @Produce json
// @Param id path int true "Target ID"
// @Param target body domain.Target true "Target details"
// @Success 200 {object} domain.Target
// @Failure 400 {object} ErrorResponse "Wrong request format or invalid target ID"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /targets/{id} [put]
func (h *Handler) updateTarget(c *gin.Context) {
	var target domain.Target
	if err := c.ShouldBindJSON(&target); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid target id"})
		c.Error(err)
		return
	}
	target.ID = id
	if err := h.usecase.UpdateTarget(c.Request.Context(), &target); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, target)
}

// deleteTarget godoc
// @Summary Delete Target
// @Description Delete a target by its ID.
// @Tags Targets
// @Param id path int true "Target ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse "Invalid target ID"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /targets/{id} [delete]
func (h *Handler) deleteTarget(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid target id"})
		c.Error(err)
		return
	}
	if err := h.usecase.DeleteTarget(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}
