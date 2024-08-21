package test

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/Car-Wash/Car-Wash-Booking-Service/config"
	"github.com/Car-Wash/Car-Wash-Booking-Service/genproto/carwash"
	"github.com/Car-Wash/Car-Wash-Booking-Service/storage/mongo"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	m "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/stretchr/testify/assert"
)

func TestBookingRepo(t *testing.T) {
	cfg := config.Load()
	uri := fmt.Sprintf("mongodb://%s:%d", cfg.MongoHost, cfg.MongoPort)
	clientOptions := options.Client().ApplyURI(uri).SetAuth(options.Credential{Username: "postgres", Password: "20005"})

	client, err := m.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Println(err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Println("Error connecting to MongoDB")
	}

	mongoDB := client.Database(cfg.MongoDatabase)
	redis := redis.Client{}
	bookingRepo := mongo.NewBookingRepo(mongoDB, &redis)

	t.Run("CreateBooking", func(t *testing.T) {

		testBooking := &carwash.CreateBookingRequest{
			UserId:        uuid.NewString(),
			ServiceId:     uuid.NewString(),
			ProviderId:    uuid.NewString(),
			Status:        "Pending",
			ScheduledTime: "90",
			Location: &carwash.GeoPoint{
				Latitude:  10.2555,
				Longitude: 10.41222,
			},
			TotalPrice: 100,
		}
		createdID, err := bookingRepo.CreateBooking(testBooking)

		assert.NoError(t, err, "CreateBooking funksiyasida xato yuz berdi")
		assert.NotEmpty(t, createdID, "Yaratilgan booking haqiqiy IDga ega bo'lishi kerak")
	})

	t.Run("GetBooking", func(t *testing.T) {

		testBooking := &carwash.CreateBookingRequest{
			UserId:        uuid.NewString(),
			ServiceId:     uuid.NewString(),
			ProviderId:    uuid.NewString(),
			Status:        "Pending",
			ScheduledTime: "90",
			Location: &carwash.GeoPoint{
				Latitude:  10.2555,
				Longitude: 10.41222,
			},
			TotalPrice: 100,
		}
		createdID, err := bookingRepo.CreateBooking(testBooking)
		assert.NoError(t, err, "Yaratish jarayonida xato yuz berdi")
		assert.NotEmpty(t, createdID, "Yaratilgan booking haqiqiy IDga ega bo'lishi kerak")

		retrievedBooking, err := bookingRepo.GetBooking(&carwash.GetBookingRequest{Id: "66bd9994b671b4f10b3ff5cd"})

		assert.NoError(t, err, "GetBooking funksiyasida xato yuz berdi")
		assert.NotNil(t, retrievedBooking, "GetBooking javobi nol bo'lmasligi kerak")

	})

	t.Run("UpdateBooking", func(t *testing.T) {
		updatedBooking := &carwash.UpdateBookingRequest{
			Id:            "66bd9994b671b4f10b3ff5cd",
			UserId:        uuid.NewString(),
			ServiceId:     uuid.NewString(),
			ProviderId:    uuid.NewString(),
			Status:        "Completed",
			ScheduledTime: "120",
			Location: &carwash.GeoPoint{
				Latitude:  20.2555,
				Longitude: 20.41222,
			},
			TotalPrice: 200,
		}
		_, err := bookingRepo.UpdateBooking(updatedBooking)
		assert.NoError(t, err, "UpdateBooking funksiyasida xato yuz berdi")

		_, err = bookingRepo.GetBooking(&carwash.GetBookingRequest{Id: "66bd9994b671b4f10b3ff5cd"})
		assert.NoError(t, err, "GetBooking funksiyasida xato yuz berdi")

	})

	t.Run("DeleteBooking", func(t *testing.T) {
		_, err = bookingRepo.DeleteBooking(&carwash.DeleteBookingRequest{Id: "66bd9994b671b4f10b3ff5cd"})
		assert.NoError(t, err, "DeleteBooking funksiyasida xato yuz berdi")

		retrievedBooking, err := bookingRepo.GetBooking(&carwash.GetBookingRequest{Id: "66bd9994b671b4f10b3ff5cd"})
		assert.Error(t, err, "GetBooking funksiyasi o'chirilgan bookingni olishda xato berishi kerak")
		assert.Nil(t, retrievedBooking, "GetBooking javobi o'chirilgan booking uchun nol bo'lishi kerak")
	})

	t.Run("ListBookings", func(t *testing.T) {

		userID := uuid.NewString()
		testBookings := []*carwash.CreateBookingRequest{
			{
				UserId:        userID,
				ServiceId:     uuid.NewString(),
				ProviderId:    uuid.NewString(),
				Status:        "Pending",
				ScheduledTime: "90",
				Location: &carwash.GeoPoint{
					Latitude:  10.2555,
					Longitude: 10.41222,
				},
				TotalPrice: 100,
			},
			{
				UserId:        userID,
				ServiceId:     uuid.NewString(),
				ProviderId:    uuid.NewString(),
				Status:        "Completed",
				ScheduledTime: "120",
				Location: &carwash.GeoPoint{
					Latitude:  20.2555,
					Longitude: 20.41222,
				},
				TotalPrice: 200,
			},
		}

		for _, booking := range testBookings {
			_, err := bookingRepo.CreateBooking(booking)
			assert.NoError(t, err, "CreateBooking funksiyasida xato yuz berdi")
		}

		req := &carwash.ListBookingsRequest{
			UserId: userID,
		}
		bookings, err := bookingRepo.ListBookings(req)
		assert.NoError(t, err, "ListBookings funksiyasida xato yuz berdi")
		assert.Len(t, bookings, len(testBookings), "ListBookings javobi kerakli miqdordagi bookinglarni qaytarishi kerak")
	})
}
