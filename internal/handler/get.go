package handler

import (
	"errors"
	"net/http"
	"strconv"

	_ "github.com/Komilov31/delayed-notifier/internal/dto"
	_ "github.com/Komilov31/delayed-notifier/internal/model"

	"github.com/Komilov31/delayed-notifier/internal/repository"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

// GetNotificationStatus godoc
// @Summary Get notification status by ID
// @Description Retrieve the status of a notification by its ID
// @Tags notifications
// @Produce json
// @Param id path int true "Notification ID"
// @Success 200 {object} dto.NotificationStatus
// @Failure 400 {object} ginext.H "Invalid ID or notification not found"
// @Failure 500 {object} ginext.H "Could not get notification status"
// @Router /notify/{id} [get]
func (h *Handler) GetNotificationStatus(c *ginext.Context) {
	id := c.Param("id")

	notifID, err := strconv.Atoi(id)
	if err != nil {
		zlog.Logger.Error().Msg("invalid id was provided: " + err.Error())
		c.JSON(http.StatusBadRequest, ginext.H{
			"error": "invalid id was provided",
		})
		return
	}

	status, err := h.service.GetNotificationStatus(notifID)
	if err != nil {
		zlog.Logger.Error().Msg("could not get notificatino status: " + err.Error())
		if errors.Is(err, repository.ErrNoSuchNotification) {
			c.JSON(http.StatusBadRequest, ginext.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, ginext.H{
			"error": "could not get notification status: " + err.Error(),
		})
		return
	}

	zlog.Logger.Info().Msgf("successfully handled GET request for getting status with id: %d", notifID)
	c.JSON(http.StatusOK, status)
}

// GetAllNotifications godoc
// @Summary Get all notifications
// @Description Retrieve a list of all notifications
// @Tags notifications
// @Produce json
// @Success 200 {array} model.Notification
// @Failure 500 {object} ginext.H "Could not get notifications"
// @Router /notify [get]
func (h *Handler) GetAllNotifications(c *ginext.Context) {
	notifications, err := h.service.GetAllNotifications()
	if err != nil {
		zlog.Logger.Error().Msg(err.Error())
		c.JSON(http.StatusInternalServerError, ginext.H{
			"error": "could not get notifications: " + err.Error(),
		})
		return
	}

	zlog.Logger.Info().Msgf("successfully handled GET request for getting all notifications")
	c.JSON(http.StatusOK, notifications)
}

// GetMainPage godoc
// @Summary Get main page
// @Description Serve the main index.html page
// @Tags main
// @Produce html
// @Success 200 {string} string "HTML page"
// @Router / [get]
func (h *Handler) GetMainPage(c *ginext.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}
