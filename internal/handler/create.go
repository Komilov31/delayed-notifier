package handler

import (
	"net/http"
	"time"

	"github.com/Komilov31/delayed-notifier/internal/dto"
	"github.com/Komilov31/delayed-notifier/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

// CreateNotification godoc
// @Summary Create a new notification
// @Description Create a new delayed notification with text and send time
// @Tags notifications
// @Accept json
// @Produce json
// @Param notification body dto.NotificationDTO true "Notification payload"
// @Success 200 {object} model.Notification
// @Failure 400 {object} ginext.H "Invalid payload or time in the past"
// @Failure 500 {object} ginext.H "Could not create notification"
// @Router /notify [post]
func (h *Handler) CreateNotification(c *gin.Context) {
	var notific dto.NotificationDTO
	if err := c.BindJSON(&notific); err != nil {
		zlog.Logger.Error().Msg("could not parse payload: " + err.Error())
		c.JSON(http.StatusBadRequest, ginext.H{
			"error": "invalid payload",
		})
		return
	}

	if time.Until(notific.SendAt) <= 0 {
		zlog.Logger.Error().Msg("invalid payload: time is in the past")
		c.JSON(http.StatusBadRequest, ginext.H{
			"error": "invalid payload: time should be in future",
		})
		return
	}

	notification := &model.Notification{
		Text:       notific.Text,
		TelegramId: notific.TelegramId,
		SendAt:     int(notific.SendAt.UnixMilli()),
	}

	notification, err := h.service.CreateNotification(*notification)
	if err != nil {
		zlog.Logger.Error().Msg("could not create notification: " + err.Error())
		c.JSON(http.StatusInternalServerError, ginext.H{
			"error": "could not create notification",
		})
		return
	}

	zlog.Logger.Info().Msgf("successfully handled POST request creating new notification")
	c.JSON(http.StatusOK, notification)
}
