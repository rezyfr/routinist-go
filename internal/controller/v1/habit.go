package v1

import (
	"net/http"
	"routinist/internal/dto"
	"routinist/internal/usecase"
	"routinist/pkg/logger"

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
}

func (r *HabitHandler) getRandomHabits(c *gin.Context) {
	response := dto.Response{}

	habits, err := r.usecase.GetRandomHabits()
	if err != nil {
		r.logger.Error(err)
		response.SetMessage("Failed to get random habits")
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	c.JSON(http.StatusOK, gin.H{"habits": habits})
}
