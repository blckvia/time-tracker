package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetUserByID(c *gin.Context) (int, error) {
	userIDStr := c.Param("id")
	if userIDStr == "" {
		newErrorResponse(c, http.StatusBadRequest, "user id is required")
		return 0, errors.New("empty user id field")
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid user id format")
		return 0, errors.New("invalid user id")
	}
	return userID, nil
}

func GetTaskByID(c *gin.Context) (int, error) {
	taskIDStr := c.Param("id")
	if taskIDStr == "" {
		newErrorResponse(c, http.StatusBadRequest, "task id is required")
		return 0, errors.New("empty task id field")
	}

	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid task id format")
		return 0, errors.New("invalid task id")
	}
	return taskID, nil
}
