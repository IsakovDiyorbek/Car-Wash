package kafka

import (
	"context"
	"encoding/json"
	"log"

	pb "github.com/exam-5/Car-Wash-Booking-Service/genproto/carwash"
	"github.com/exam-5/Car-Wash-Booking-Service/service"
)

func PaymentHandler(paymentservice *service.PaymentService) func(message []byte) {
	return func(message []byte) {
		var payment pb.CreatePaymentRequest
		if err := json.Unmarshal(message, &payment); err != nil {
			log.Printf("Cannot unmarshal JSON: %v", err)
			return
		}

		resppayment, err := paymentservice.CreatePayment(context.Background(), &payment)
		if err != nil {
			log.Printf("Cannot create payment via Kafka: %v", err)
			return
		}
		log.Printf("Created payment: %+v",resppayment)
	}
}