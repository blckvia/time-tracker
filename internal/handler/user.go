package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"time-tracker/internal/entities"
)

// @Summary Create user
// @Tags Users
// @Description Create user
// @ID create-user
// @Accept  json
// @Produce  json
// @Param input body entities.Users true "users info"
// @Success 200 {object} entities.Users
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/users/ [post]
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

// @Summary Get all users
// @Tags Users
// @Description Get all users
// @ID get-users
// @Accept  json
// @Produce  json
// @Param limit query int false "limit"
// @Param offset query int false "offset"
// @Success 200 {array} entities.GetAllUsers
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/users/ [get]
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

// @Summary Get user
// @Tags Users
// @Description Get single user
// @ID get-user
// @Accept  json
// @Produce  json
// @Param user_id path int true "user_id"
// @Success 200 {object} entities.Users
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/users/{user_id} [get]
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

// @Summary Delete user
// @Tags Users
// @Description Delete user
// @ID delete-user
// @Accept  json
// @Produce  json
// @Param user_id path int true "user_id"
// @Success 200
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/users/{user_id} [delete]
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

// @Summary Update user
// @Tags Users
// @Description Update user
// @ID update-user
// @Accept  json
// @Produce  json
// @Param input body entities.Users true "goods info"
// @Param user_id path int true "user_id"
// @Success 200 {object} entities.Users
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/users/{user_id} [put]
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

// @Summary Get user stats
// @Tags Users
// @Description Get user stats with overall time and tasks
// @ID get-stats-user
// @Accept  json
// @Produce  json
// @Param input body entities.UserStats true "stats info"
// @Param user_id path int true "user_id"
// @Success 200
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/users/{user_id}/stats [get]
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
