package handler

import (
	"github.com/exam-5/Car-Wash-Api-Gateway/genproto/carwash"
	"github.com/gin-gonic/gin"
)



// @Summary Get all notifications
// @Description Get all notifications
// @Tags Notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param booking_id query string true "booking_id"
// @Success 200 {object} carwash.GetNotificationsResponse
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /notifications/{id} [get]
func (h *Handler) GetNotigication(c *gin.Context){
	req := carwash.GetNotificationsRequest{}
	req.BookingId = c.Query("booking_id")

	res, err := h.Client.Notification.GetNotifications(c, &req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, res)
}



















































































