package main

import (
	"log"
	"log/slog"
	"net"

	"github.com/exam-5/Car-Wash-Booking-Service/config"
	"github.com/exam-5/Car-Wash-Booking-Service/genproto/carwash"
	"github.com/exam-5/Car-Wash-Booking-Service/kafka"
	"github.com/exam-5/Car-Wash-Booking-Service/service"
	"github.com/exam-5/Car-Wash-Booking-Service/storage/mongo"
	"google.golang.org/grpc"
)

func main() {

	cfg := config.Load()
	db, err := mongo.SetupMongoDBConnection(cfg)
	if err != nil {
		slog.Info("Error connecting mongo db", err)
	}

	liss, err := net.Listen("tcp", cfg.HTTPPort)

	if err != nil {
		slog.Info("Error listening port", err)
	}

	bookingService := service.NewBookingService(db)
	paymentService := service.NewPaymentService(db)
	reviewService := service.NewReviewService(db)

	brokers := []string{"kafka:9092"}

	kcm := kafka.NewKafkaConsumerManager()

	if err := kcm.RegisterConsumer(brokers, "create-booking", "product", kafka.BookingHandler(bookingService)); err != nil {
		if err == kafka.ErrConsumerAlreadyExists {
			log.Printf("Consumer for topic 'cr-booking' already exists")
		} else {
			log.Printf("Failed to register consumer for topic 'cr-booking': %v", err)

		}
	}
	if err := kcm.RegisterConsumer(brokers, "update-booking", "product", kafka.UpdateHandler(bookingService)); err != nil {
		if err == kafka.ErrConsumerAlreadyExists {
			log.Printf("Consumer for topic 'cr-booking' already exists")
		} else {
			log.Printf("Failed to register consumer for topic 'cr-booking': %v", err)

		}
	}
	if err := kcm.RegisterConsumer(brokers, "delete-booking", "product", kafka.DeleteBookingHandler(bookingService)); err != nil {
		if err == kafka.ErrConsumerAlreadyExists {
			log.Printf("Consumer for topic 'cr-booking' already exists")
		} else {
			log.Printf("Failed to register consumer for topic 'cr-booking': %v", err)

		}
	}

	if err := kcm.RegisterConsumer(brokers, "payment", "product", kafka.PaymentHandler(paymentService)); err != nil {
		if err == kafka.ErrConsumerAlreadyExists {
			log.Printf("Consumer for topic 'payment' already exists")
		} else {
			log.Printf("Failed to register consumer for topic 'payment': %v", err)

		}
	}
	if err := kcm.RegisterConsumer(brokers, "review", "product", kafka.ReviewHandler(reviewService)); err != nil {
		if err == kafka.ErrConsumerAlreadyExists {
			log.Printf("Consumer for topic 'review' already exists")
		} else {
			log.Printf("Failed to register consumer for topic 'review': %v", err)

		}
	}

	s := grpc.NewServer()
	carwash.RegisterBookingServiceServer(s, service.NewBookingService(db))
	carwash.RegisterServicesServiceServer(s, service.NewServicesService(db))
	carwash.RegisterProviderServiceServer(s, service.NewProviderService(db))
	carwash.RegisterReviewServiceServer(s, service.NewReviewService(db))
	carwash.RegisterPaymentServiceServer(s, service.NewPaymentService(db))

	log.Printf("Server started on port: %v", cfg.HTTPPort)
	if err := s.Serve(liss); err != nil {
		log.Fatalf("error while serving: %v", err)
	}

}
