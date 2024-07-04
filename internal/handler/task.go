package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"time-tracker/internal/entities"
)

// @Summary Create task
// @Tags Tasks
// @Description Create task
// @ID create-task
// @Accept  json
// @Produce  json
// @Param input body entities.Task true "task info"
// @Param user_id path int true "user_id"
// @Success 200 {object} entities.Task
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/tasks/{user_id} [post]
func (h *Handler) createTask(c *gin.Context) {
	userID, err := GetUserByID(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var input entities.Task
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if input.Task == "" {
		newErrorResponse(c, http.StatusBadRequest, "task is required")
		return
	}

	id, err := h.services.Tasks.Create(&input, userID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"id": id,
	})
}

// @Summary Start task
// @Tags Tasks
// @Description Start task
// @ID start-task
// @Accept  json
// @Produce  json
// @Param task_id path int true "task_id"
// @Success 200 {object} statusResponse
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/tasks/{task_id}/start [post]
func (h *Handler) startTask(c *gin.Context) {
	taskID, err := GetTaskByID(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	err = h.services.Tasks.StartTask(taskID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{"Task started successfully"})
}

// @Summary Stop task
// @Tags Tasks
// @Description Stop task
// @ID stop-task
// @Accept  json
// @Produce  json
// @Param task_id path int true "task_id"
// @Success 200 {object} map[string]string
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/tasks/{task_id}/stop [post]
func (h *Handler) stopTask(c *gin.Context) {
	taskID, err := GetTaskByID(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	overallTime, err := h.services.Tasks.StopTask(taskID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]string{
		"overall_time": overallTime.String(),
		"message":      "Task stopped successfully",
	})
}
