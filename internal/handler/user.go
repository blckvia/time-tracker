package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"time-tracker/internal/entities"
)

func (h *Handler) createUser(c *gin.Context) {
	var input entities.Users
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if input.Name == "" || input.Surname == "" || input.PassportSeries == "" || input.PassportSeries == "" || input.Address == "" {
		newErrorResponse(c, http.StatusBadRequest, "invalid input. Name, Surname, PassportSeries, PassportNumber, Address should not be empty")
		return
	}
	if _, err := strconv.Atoi(input.PassportSeries); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "series should consist of 4 numbers")
		return
	}
	if _, err := strconv.Atoi(input.PassportNumber); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "series should consist of 6 numbers")
		return
	}

	id, err := h.services.Users.Create(&input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]int{
		"id": id,
	})
}

func (h *Handler) getUsers(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	filters := map[string]string{}
	if passportSeries := c.Query("passport_series"); passportSeries != "" {
		filters["passport_series"] = passportSeries
	}
	if passportNumber := c.Query("passport_number"); passportNumber != "" {
		filters["passport_number"] = passportNumber
	}
	if name := c.Query("name"); name != "" {
		filters["name"] = name
	}
	if surname := c.Query("surname"); surname != "" {
		filters["surname"] = surname
	}
	if patronymic := c.Query("patronymic"); patronymic != "" {
		filters["patronymic"] = patronymic
	}
	if address := c.Query("address"); address != "" {
		filters["address"] = address
	}

	users, err := h.services.Users.GetAll(filters, limit, offset)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h *Handler) getUser(c *gin.Context) {
	userID, err := GetUserByID(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	user, err := h.services.Users.GetByID(userID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) deleteUser(c *gin.Context) {
	UserID, err := GetUserByID(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.Users.Delete(UserID); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})
}

func (h *Handler) updateUser(c *gin.Context) {
	userID, err := GetUserByID(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var input entities.Users
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.Users.Update(userID, &input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})
}

func (h *Handler) getUsersStats(c *gin.Context) {
	userID, err := GetUserByID(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	userStats, err := h.services.Users.Stats(userID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, userStats)
}
