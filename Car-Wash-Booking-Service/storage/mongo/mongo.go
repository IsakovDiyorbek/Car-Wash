package mongo

import (
	"context"
	"fmt"

	"github.com/exam-5/Car-Wash-Booking-Service/config"
	"github.com/exam-5/Car-Wash-Booking-Service/storage"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type StorageM struct {
	mongo     *mongo.Database
	Services  storage.Service
	Payments  storage.PaymentI
	Reviews   storage.Review
	Bookings  storage.BookingI
	Providers storage.ProviderI
	Notifications storage.NotificationI
}

func SetupMongoDBConnection(cfg config.Config) (storage.StorageI, error) {
	uri := fmt.Sprintf("mongodb://%s:%d", cfg.MongoHost, cfg.MongoPort)
	clientOptions := options.Client().ApplyURI(uri).SetAuth(options.Credential{Username: "postgres", Password: "20005"})

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	mongoDB := client.Database(cfg.MongoDatabase)

	return &StorageM{
		mongo:     mongoDB,
		Services:  NewServiceRepo(mongoDB),
		Providers: NewProviderRepo(mongoDB),
		Bookings:  NewBookingRepo(mongoDB),
		Reviews:   NewReviewRepo(mongoDB),
		Payments:  NewPaymentRepo(mongoDB),
		Notifications: NewNotificationManager(mongoDB),
	}, nil
}

func (s *StorageM) Service() storage.Service {
	if s.Services == nil {
		s.Services = NewServiceRepo(s.mongo)
	}
	return s.Services
}

func (s *StorageM) Provider() storage.ProviderI {
	if s.Providers == nil {
		s.Providers = NewProviderRepo(s.mongo)
	}
	return s.Providers
}

func (s *StorageM) Booking() storage.BookingI {
	if s.Bookings == nil { 
		s.Bookings = NewBookingRepo(s.mongo)
	}
	return s.Bookings
}

func (s *StorageM) Review() storage.Review {
	if s.Reviews == nil {
		s.Reviews = NewReviewRepo(s.mongo)
	}
	return s.Reviews
}

func (s *StorageM) Payment() storage.PaymentI {
	if s.Payments == nil {
		s.Payments = NewPaymentRepo(s.mongo)
	}
	return s.Payments
}

func (s *StorageM) Notification() storage.NotificationI {
	if s.Notifications == nil {
		s.Notifications = NewNotificationManager(s.mongo)
	}
	return s.Notifications
}