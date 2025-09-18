package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Komilov31/delayed-notifier/internal/repository"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

// UpdateNotificationStatus godoc
// @Summary Cancel a notification by updating its status
// @Description Update the status of a notification to "canceled" by its ID
// @Tags notifications
// @Produce json
// @Param id path int true "Notification ID"
// @Success 200 {object} ginext.H "Notification cancellation status"
// @Failure 400 {object} ginext.H "Invalid ID or notification not found"
// @Failure 500 {object} ginext.H "Could not update notification status"
// @Router /notify/{id} [delete]
func (h *Handler) UpdateNotificationStatus(c *ginext.Context) {
	id := c.Param("id")

	notifID, err := strconv.Atoi(id)
	if err != nil {
		zlog.Logger.Error().Msg("invalid id was provided: " + err.Error())
		c.JSON(http.StatusBadRequest, ginext.H{
			"error": "invalid id was provided",
		})
		return
	}

	err = h.service.UpdateNotificationStatus(notifID, "canceled")
	if err != nil {
		zlog.Logger.Error().Msg("could not update notification status: " + err.Error())
		if errors.Is(err, repository.ErrNoSuchNotification) {
			c.JSON(http.StatusBadRequest, ginext.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, ginext.H{
			"error": "could not update notification status: " + err.Error(),
		})
		return
	}

	zlog.Logger.Info().Msgf("successfully handled DELETE request for updating notif status with id: %d", notifID)
	c.JSON(http.StatusOK, ginext.H{
		"status": "notification was cancelled succesfully",
	})
}
