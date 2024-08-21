package handler

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/Car-Wash/Car-Wash-Api-Gateway/genproto/carwash"
	"github.com/gin-gonic/gin"
)

// @Summary Create a new booking
// @Description Create a new booking
// @Tags Bookings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param booking body carwash.CreateBookingRequest true "booking"
// @Success 200 {object} carwash.CreateBookingResponse
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /bookings [post]
func (h *Handler) CreateBooking(c *gin.Context) {
	req := carwash.CreateBookingRequest{}
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// pending - kutilyapti
	// confirmed - tasdiqlangan
	// completed - yakunlangan
	// cancelled- bekorQilingan

	if req.Status != "Pending" {
		c.JSON(400, gin.H{"error": "Invalid status value"})
		return
	}

	res, err := h.Client.Provider.GetProvider(c, &carwash.GetProviderRequest{Id: req.ProviderId})
	if err != nil {
		c.JSON(400, gin.H{"Provider not found": err.Error()})
		return
	}

	var check bool
	for _, id := range res.Provider.ServiceId {
		if id == req.ServiceId {
			check = true
			break
		}
	}

	if !check {
		c.JSON(400, gin.H{"error": "Service id not found"})
		return
	}
	// check if service exist
	_, err = h.Client.Service.GetService(c, &carwash.GetServiceRequest{Id: req.ServiceId})
	if err != nil {
		c.JSON(400, gin.H{"Service not found": err.Error()})
		return
	}

	// // Check if service exists
	// _, err = h.Client.Service.GetService(c, &carwash.GetServiceRequest{Id: req.ServiceId})
	// if err != nil {
	//     c.JSON(400, gin.H{"User not found": err.Error()})
	//     return
	// }

	// Create Book
	input, err := json.Marshal(req)
	err = h.Client.Kafka.ProduceMessages("create-booking", input)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		log.Println("cannot produce messages via kafka", err)
		return
	}

	c.JSON(200, gin.H{"message": "Booking created"})
}

// @Summary Get all bookings
// @Description Get all bookings
// @Tags Bookings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param provider_id query string false "provider_id"
// @Param user_id query string false "user_id"
// @Param scheduled_time query string false "scheduled_time"
// @Param status query string false "status"
// @Param service_id query string false "service_id"
// @Success 200 {object} carwash.ListBookingsResponse
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /bookings [get]
func (h *Handler) GetBookings(c *gin.Context) {
	req := carwash.ListBookingsRequest{}

	req.ProviderId = c.Query("provider_id")
	req.UserId = c.Query("user_id")
	req.ScheduledTime = c.Query("scheduled_time")
	req.Status = c.Query("status")
	req.ServiceId = c.Query("service_id")

	res, err := h.Client.Booking.ListBookings(c, &req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, res)
}

// @Summary Get a booking
// @Description Get a booking
// @Tags Bookings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id query string true "id"
// @Success 200 {object} carwash.GetBookingResponse
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /bookings/{id} [get]
func (h *Handler) GetBooking(c *gin.Context) {
	id := c.Query("id")
	res, err := h.Client.Booking.GetBooking(c, &carwash.GetBookingRequest{Id: id})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, res)
}

// @Summary Update a booking
// @Description Update a booking
// @Tags Bookings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id query string true "id"
// @Param status query string false "status"
// @Param user_id query string false "user_id"
// @Param provider_id query string false "provider_id"
// @Param service_id query string false "service_id"
// @Param scheduled_time query string false "scheduled_time"
// @Param total_price query number false "total_price"
// @Success 200 {object} carwash.UpdateBookingResponse
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /bookings/{id} [put]
func (h *Handler) UpdateBooking(c *gin.Context) {
	req := carwash.UpdateBookingRequest{}
	req.Id = c.Query("id")
	req.Status = c.Query("status")
	req.UserId = c.Query("user_id")
	req.ProviderId = c.Query("provider_id")
	req.ServiceId = c.Query("service_id")
	req.ScheduledTime = c.Query("scheduled_time")
	totalPriceStr := c.Query("total_price")
	if totalPriceStr != "" {
		totalPriceFloat, err := strconv.ParseFloat(totalPriceStr, 32)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		req.TotalPrice = float32(totalPriceFloat)
	}
	if req.ProviderId != "" {
		_, err := h.Client.Provider.GetProvider(c, &carwash.GetProviderRequest{Id: req.ProviderId})
		if err != nil {
			c.JSON(400, gin.H{"Provider not found": err.Error()})
			return
		}
	}

	if req.ServiceId != "" {
		_, err := h.Client.Service.GetService(c, &carwash.GetServiceRequest{Id: req.ServiceId})
		if err != nil {
			c.JSON(400, gin.H{"Service not found": err.Error()})
			return
		}
	}

	input, err := json.Marshal(req)
	err = h.Client.Kafka.ProduceMessages("update-booking", input)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		log.Println("cannot produce messages via kafka", err)
		return
	}

	c.JSON(200, gin.H{"message": "Booking updated"})
}

// @Summary Delete a booking
// @Description Delete a booking
// @Tags Bookings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id query string true "id"
// @Success 200 {object} carwash.DeleteBookingResponse
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /bookings/{id} [delete]
func (h *Handler) DeleteBooking(c *gin.Context) {
	id := c.Query("id")

	input, err := json.Marshal(id)
	err = h.Client.Kafka.ProduceMessages("delete-booking", input)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		log.Println("cannot produce messages via kafka", err)
		return
	}

	c.JSON(200, gin.H{"message": "Booking deleted"})
}

// @Summary Confirm a booking
// @Description Confirm a booking by ID
// @Tags Bookings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id query string true "Booking ID"
// @Success 200 {object} carwash.UpdateBookingResponse
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /bookings/{id}/confirm [put]
func (h *Handler) ConfirmBooking(c *gin.Context) {
	bookingID := c.Query("id")

	res, err := h.Client.Booking.GetBooking(c, &carwash.GetBookingRequest{Id: bookingID})
	if err != nil {
		c.JSON(400, gin.H{"error": "Booking not found"})
		return
	}
	if res.Booking.Status == "Pending" {
		res.Booking.Status = "Confirmed"
		_, err = h.Client.Booking.UpdateBooking(c, &carwash.UpdateBookingRequest{
			Id:     bookingID,
			Status: res.Booking.Status,
		})

	} else {
		c.JSON(400, gin.H{"error": "Booking is not pending"})
		return
	}

	c.JSON(200, gin.H{"message": "Booking confirmed"})
}
