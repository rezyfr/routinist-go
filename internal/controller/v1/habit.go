package v1

import (
	"net/http"
	"routinist/internal/dto/request"
	"routinist/internal/dto/response"
	"routinist/internal/middleware"
	"routinist/internal/usecase"
	"routinist/pkg/logger"
	"strconv"
	"time"

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
		auth.POST("/create", r.createUserHabit)
		auth.GET("/today", r.getTodayHabits)
		auth.POST("/:user_habit_id/progress", r.postCreateProgress)
		auth.GET("/summary", r.GetSummaryProgress)
	}
}

func (h *HabitHandler) createUserHabit(c *gin.Context) {
	r := response.Response{}

	userIDVal, _ := c.Get("user_id")
	userId := userIDVal.(uint)

	var req request.CreateUserHabitRequestDTO

	if err := c.Bind(&req); err != nil {
		r.SetMessage("Invalid request")
		c.JSON(http.StatusBadRequest, r)
		return
	}

	_, err := h.usecase.CreateUserHabit(userId, req.HabitId, &req.UnitId, &req.Goal)

	if err != nil {
		h.logger.Error(err)
		r.SetMessage("Failed to create user habit")
		c.JSON(http.StatusInternalServerError, r)
	}

	r.Data = "User habit created successfully"

	c.JSON(http.StatusOK, r)
	return
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

	if req.Value < 0 {
		r.SetMessage("Value must be more than 0")
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

func (h *HabitHandler) GetSummaryProgress(c *gin.Context) {
	r := response.Response{}

	userIDVal, _ := c.Get("user_id")

	userId := userIDVal.(uint)

	mode := c.DefaultQuery("mode", "today")

	var from, to time.Time
	now := time.Now()
	switch mode {
	case "today":
		from = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		to = from.Add(24 * time.Hour)
	case "this_week":
		weekday := int(now.Weekday())
		from = time.Date(now.Year(), now.Month(), now.Day()-weekday, 0, 0, 0, 0, now.Location())
		to = from.AddDate(0, 0, 7)
	case "this_month":
		from = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		to = time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location())
	default:
		from = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		to = from.Add(24 * time.Hour)
	}

	d, e := h.usecase.GetProgressSummary(userId, from, to)

	if e != nil {
		h.logger.Error(e)
		r.SetMessage("Failed to get aggregate progress")
		c.JSON(http.StatusInternalServerError, r)
		return
	}

	r.Data = d
	c.JSON(http.StatusOK, r)
}
