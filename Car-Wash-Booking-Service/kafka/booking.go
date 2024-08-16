package kafka

import (
	"context"
	"encoding/json"
	"log"

	pb "github.com/exam-5/Car-Wash-Booking-Service/genproto/carwash"
	"github.com/exam-5/Car-Wash-Booking-Service/service"
)

func BookingHandler(bookingservice *service.BookingService) func(message []byte) {
	return func(message []byte) {
		var booking pb.CreateBookingRequest
		if err := json.Unmarshal(message, &booking); err != nil {
			log.Printf("Cannot unmarshal JSON: %v", err)
			return
		}

		respbooking, err := bookingservice.CreateBooking(context.Background(), &booking)
		if err != nil {
			log.Printf("Cannot create booking via Kafka: %v", err)
			return
		}
		log.Printf("Created booking: %+v",respbooking)
	}
}

func UpdateHandler(bookservice *service.BookingService) func(message []byte) {
	return func(message []byte) {
		var booking pb.UpdateBookingRequest
		if err := json.Unmarshal(message, &booking); err != nil {
			log.Printf("Cannot unmarshal JSON: %v", err)
			return
		}

		respbooking, err := bookservice.UpdateBooking(context.Background(), &booking)
		if err != nil {
			log.Printf("Cannot create booking via Kafka: %v", err)
			return
		}
		log.Printf("Created booking: %+v",respbooking)
	}
}
func DeleteBookingHandler(bookingservice *service.BookingService) func(message []byte) {
	return func(message []byte) {
		var booking pb.DeleteBookingRequest
		if err := json.Unmarshal(message, &booking); err != nil {
			log.Printf("Cannot unmarshal JSON: %v", err)
			return
		}

		respbooking, err := bookingservice.DeleteBooking(context.Background(), &booking)
		if err != nil {
			log.Printf("Cannot create booking via Kafka: %v", err)
			return
		}
		log.Printf("Created booking: %+v",respbooking)
	}
}