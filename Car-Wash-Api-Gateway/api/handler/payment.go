package handler

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/exam-5/Car-Wash-Api-Gateway/genproto/carwash"
	"github.com/gin-gonic/gin"
)

// @Summary Create a new payment
// @Description Create a new payment
// @Tags Payments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payment body carwash.CreatePaymentRequest true "payment"
// @Success 200 {object} carwash.CreatePaymentResponse
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /payments [post]
func (h *Handler) CreatePayment(c *gin.Context) {
	req := carwash.CreatePaymentRequest{}
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	res, err := h.Client.Booking.GetBooking(c, &carwash.GetBookingRequest{Id: req.BookingId})
	if err != nil {
		c.JSON(400, gin.H{"error": "Booking not found"})
		return
	}
	if res.Booking.Status != "Confirmed" {
		c.JSON(400, gin.H{"error": "Booking is not confirmed"})
		return
	}

	if res.Booking.TotalPrice != req.Amount {
		c.JSON(400, gin.H{"error": "Amount does not match"})
		return
	}

	req.Status = "Completed"
	input, err := json.Marshal(req)
	err = h.Client.Kafka.ProduceMessages("payment", input)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		log.Println("cannot produce messages via kafka", err)
		return
	}

	notification := carwash.AddNotificationRequest{Message: "Success create payment", IsRead: true, BookingId: req.BookingId}
	request, err := json.Marshal(notification)
	err = h.Client.Kafka.ProduceMessages("notification", request)
	if err != nil {
		log.Println("Error while adding notification", err)
	}

	c.JSON(200, gin.H{"message": "Payment created, status will be updated shortly"})
}

// @Summary Get a payment
// @Description Get a payment
// @Tags Payments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id query string true "id"
// @Success 200 {object} carwash.GetPaymentResponse
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /payments/{id} [get]
func (h *Handler) GetPayment(c *gin.Context) {
	req := carwash.GetPaymentRequest{}
	req.Id = c.Query("id")

	res, err := h.Client.Payment.GetPayment(c, &req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, res)
}

// @Summary List payments
// @Description List payments with optional filters
// @Tags Payments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param booking_id query string false "booking_id"
// @Param amount query number false "amount"
// @Param status query string false "status"
// @Param payment_method query string false "payment_method"
// @Param transaction_id query string false "transaction_id"
// @Param limit query int false "limit"
// @Param offset query int false "offset"
// @Success 200 {object} carwash.ListPaymentsResponse
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /payments [get]
func (h *Handler) ListPayments(c *gin.Context) {
	req := carwash.ListPaymentsRequest{}

	req.BookingId = c.Query("booking_id")

	amountStr := c.Query("amount")
	if amountStr != "" {
		amount, err := strconv.ParseFloat(amountStr, 32)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid amount value"})
			return
		}
		req.Amount = float32(amount)
	}

	req.Status = c.Query("status")
	req.PaymentMethod = c.Query("payment_method")
	req.TransactionId = c.Query("transaction_id")

	limitStr := c.Query("limit")
	if limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid limit value"})
			return
		}
		req.Limit = int32(limit)
	}

	offsetStr := c.Query("offset")
	if offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid offset value"})
			return
		}
		req.Offset = int32(offset)
	}

	res, err := h.Client.Payment.ListPayments(c, &req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, res)
}
