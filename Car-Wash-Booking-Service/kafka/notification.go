package kafka

import (
	"context"
	"encoding/json"
	"log"

	pb "github.com/exam-5/Car-Wash-Booking-Service/genproto/carwash"
	"github.com/exam-5/Car-Wash-Booking-Service/service"
)

func NotifiactionHandler(notificationservice *service.NotificationService) func(message []byte) {
	return func(message []byte) {
		var notification pb.AddNotificationRequest
		if err := json.Unmarshal(message, &notification); err != nil {
			log.Printf("Cannot unmarshal JSON: %v", err)
			return
		}

		respnotification, err := notificationservice.AddNotification(context.Background(), &notification)
		if err != nil {
			log.Printf("Cannot create notification via Kafka: %v", err)
			return
		}
		log.Printf("Created notification: %+v",respnotification)
	}
}
