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
