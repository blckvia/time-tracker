package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"time-tracker/internal/entities"
)

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

func (h *Handler) stopTask(c *gin.Context) {

}
