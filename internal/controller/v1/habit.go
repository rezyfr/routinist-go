package v1

import (
	"net/http"
	"routinist/internal/dto/request"
	"routinist/internal/dto/response"
	"routinist/internal/middleware"
	"routinist/internal/usecase"
	"routinist/pkg/logger"
	"strconv"

	"github.com/gin-gonic/gin"
)

type HabitHandler struct {
	usecase usecase.HabitUsecase
	logger  logger.Interface
}

func NewHabitRoutes(handler *gin.RouterGroup, t usecase.HabitUsecase, l logger.Interface) {
	r := &HabitHandler{t, l}

	h1 := handler.Group("/habit")
	{
		h1.GET("/random", r.getRandomHabits)
	}

	auth := handler.Group("/protected/habit", middleware.JWTAuthMiddleware())
	{
		auth.GET("/today", r.getTodayHabits)
		auth.POST("/:user_habit_id/progress", r.postCreateProgress)
	}
}

func (h *HabitHandler) getRandomHabits(c *gin.Context) {
	r := response.Response{}

	habits, err := h.usecase.GetRandomHabits()
	if err != nil {
		h.logger.Error(err)
		r.SetMessage("Failed to get random habits")
		c.JSON(http.StatusInternalServerError, r)
		return
	}

	r.Data = habits
	c.JSON(http.StatusOK, r)
}

func (h *HabitHandler) getTodayHabits(c *gin.Context) {
	r := response.Response{}

	userIDVal, _ := c.Get("user_id")

	userId := userIDVal.(uint)

	habits, err := h.usecase.GetTodayHabits(userId)
	if err != nil {
		h.logger.Error(err)
		r.SetMessage("Failed to get today habits")
		c.JSON(http.StatusInternalServerError, r)
		return
	}

	r.Data = habits
	c.JSON(http.StatusOK, r)
}

func (h *HabitHandler) postCreateProgress(c *gin.Context) {
	r := response.Response{}

	habitIdVal := c.Param("user_habit_id")
	if habitIdVal == "" {
		r.SetMessage("Invalid habit ID")
		c.JSON(http.StatusBadRequest, r)
		return
	}

	habitId, e := strconv.Atoi(habitIdVal)

	if e != nil {
		r.SetMessage("Invalid habit ID")
		c.JSON(http.StatusBadRequest, r)
	}

	userIDVal, _ := c.Get("user_id")
	userId := userIDVal.(uint)

	var req request.CreateHabitProgressRequestDTO

	if err := c.Bind(&req); err != nil {
		r.SetMessage("Invalid request")
		c.JSON(http.StatusBadRequest, r)
		return
	}

	if req.Value < 0 || req.Value > 100 {
		r.SetMessage("Value must be between 0 and 100")
		c.JSON(http.StatusBadRequest, r)
		return
	}

	updatedHabit, err := h.usecase.PostCreateHabitProgress(userId, uint(habitId), req.Value)

	if err != nil {
		h.logger.Error(err)
		r.SetMessage("Failed to create habit progress")
		c.JSON(http.StatusInternalServerError, r)
		return
	}

	r.Data = updatedHabit
	c.JSON(http.StatusOK, r)
}
