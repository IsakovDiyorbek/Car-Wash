package kafka

import (
	"context"
	"encoding/json"
	"log"

	pb "github.com/Car-Wash/Car-Wash-Booking-Service/genproto/carwash"
	"github.com/Car-Wash/Car-Wash-Booking-Service/service"
)

func ReviewHandler(reviewservice *service.ReviewService) func(message []byte) {
	return func(message []byte) {
		var review pb.CreateReviewRequest
		if err := json.Unmarshal(message, &review); err != nil {
			log.Printf("Cannot unmarshal JSON: %v", err)
			return
		}

		respreview, err := reviewservice.CreateReview(context.Background(), &review)
		if err != nil {
			log.Printf("Cannot create review via Kafka: %v", err)
			return
		}
		log.Printf("Created review: %+v", respreview)
	}
}
