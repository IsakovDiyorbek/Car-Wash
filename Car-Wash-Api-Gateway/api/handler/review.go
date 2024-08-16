package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/exam-5/Car-Wash-Api-Gateway/genproto/carwash"
	"github.com/gin-gonic/gin"
)

// @Summary Create a new review
// @Description Create a new review
// @Tags Reviews
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param review body carwash.CreateReviewRequest true "review"
// @Success 200 {object} carwash.CreateReviewResponse
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /reviews [post]
func (h *Handler) CreateReview(c *gin.Context) {
	req := carwash.CreateReviewRequest{}
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

	_, err = h.Client.Provider.GetProvider(c, &carwash.GetProviderRequest{Id: res.Booking.ProviderId})
	if err != nil {
		c.JSON(400, gin.H{"error": "Provider not found"})
		return
	}

	if req.ProviderId != res.Booking.ProviderId {
		c.JSON(400, gin.H{"error": "Provider does not match"})
		return
	}

	if !(req.Rating >= 1 && req.Rating <= 5) {
		c.JSON(400, gin.H{"error": "Rating must be between 1 and 5"})
		return
	}
	_, err = h.Client.Booking.UpdateBooking(c, &carwash.UpdateBookingRequest{
		Id:     req.BookingId,
		Status: "Completed",
	})

	input, err := json.Marshal(req)
	err = h.Client.Kafka.ProduceMessages("review", input)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		log.Println("cannot produce messages via kafka", err)
		return
	}

	go func() {
		providerReviews, err := h.Client.Review.ListReviews(c, &carwash.ListReviewsRequest{
			ProviderId: res.Booking.ProviderId,
		})
		if err != nil {
			log.Fatal("Error while getting provider reviews:", err)
		}

		var sumRating float32
		for _, review := range providerReviews.Reviews {
			sumRating += review.Rating
		}

		avarageRating := sumRating / float32((len(providerReviews.Reviews)))

		_, err = h.Client.Provider.UpdateProvider(c, &carwash.UpdateProviderRequest{
			Id:            res.Booking.ProviderId,
			AverageRating: float32(avarageRating),
		})
		if err != nil {
			fmt.Println("Failed to update provider rating:", err)
		}
	}()

	c.JSON(200, gin.H{"message": "Review created and booking status updated"})
}

// @Summary Get all reviews
// @Description Get all reviews
// @Tags Reviews
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param rating query number false "rating"
// @Param providerId query string false "providerId"
// @Param userId query string false "userId"
// @Param limit query int false "limit"
// @Param offset query int false "offset"
// @Success 200 {object} carwash.ListReviewsResponse
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /reviews [get]
func (h *Handler) GetReviews(c *gin.Context) {
	req := carwash.ListReviewsRequest{}

	ratingStr := c.Query("rating")
	if ratingStr != "" {
		ratingFloat, err := strconv.Atoi(ratingStr)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		req.Rating = float32(ratingFloat)
	}

	limitStr := c.Query("limit")
	if limitStr != "" {
		limitInt, err := strconv.Atoi(limitStr)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		req.Limit = int32(limitInt)
	}

	offsetStr := c.Query("offset")
	if offsetStr != "" {
		offsetInt, err := strconv.Atoi(offsetStr)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		req.Offset = int32(offsetInt)
	}

	req.ProviderId = c.Query("providerId")
	req.BookingId = c.Query("bookingId")
	req.Comment = c.Query("comment")

	res, err := h.Client.Review.ListReviews(c, &req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, res)
}

// @Summary Update a review
// @Description Update a review
// @Tags Reviews
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id query string true "id"
// @Param rating query number false "rating"
// @Param comment query string false "comment"
// @Success 200 {object} carwash.UpdateReviewResponse
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /reviews/{id} [put]
func (h *Handler) UpdateReview(c *gin.Context) {
	req := carwash.UpdateReviewRequest{}
	req.Id = c.Query("id")
	req.BookingId = c.Query("bookingId")
	req.ProviderId = c.Query("providerId")
	req.UserId = c.Query("userId")

	ratingStr := c.Query("rating")
	if ratingStr != "" {
		ratingFloat, err := strconv.Atoi(ratingStr)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid rating value"})
			return
		}
		req.Rating = float32(ratingFloat)
	}

	req.Comment = c.Query("comment")

	req.Id = c.Query("id")
	if req.Id == "" {
		c.JSON(400, gin.H{"error": "Missing review ID"})
		return
	}

	_, err := h.Client.Review.UpdateReview(c, &req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Review updated successfully"})
}

// @Summary Delete a review
// @Description Delete a review
// @Tags Reviews
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id query string true "id"
// @Success 200 {object} carwash.DeleteReviewResponse
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /reviews/{id} [delete]
func (h *Handler) DeleteReview(c *gin.Context) {
	req := carwash.DeleteReviewRequest{}
	req.Id = c.Query("id")

	_, err := h.Client.Review.DeleteReview(c, &req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Review deleted successfully"})
}

// @Summary Get a review
// @Description Get a review
// @Tags Reviews
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id query string true "id"
// @Success 200 {object} carwash.GetReviewResponse
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /reviews/{id} [get]
func (h *Handler) GetReview(c *gin.Context) {
	req := carwash.GetReviewRequest{}
	req.Id = c.Query("id")

	res, err := h.Client.Review.GetReview(c, &req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, res)
}
