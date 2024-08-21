package test

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/Car-Wash/Car-Wash-Booking-Service/config"
	"github.com/Car-Wash/Car-Wash-Booking-Service/genproto/carwash"
	"github.com/Car-Wash/Car-Wash-Booking-Service/storage/mongo"
	"github.com/stretchr/testify/assert"
	m "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestNotificationManager(t *testing.T) {
	cfg := config.Load()
	uri := fmt.Sprintf("mongodb://%s:%d", cfg.MongoHost, cfg.MongoPort)
	clientOptions := options.Client().ApplyURI(uri).SetAuth(options.Credential{Username: "postgres", Password: "20005"})

	client, err := m.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal("Error connecting to MongoDB")
	}

	mongoDB := client.Database(cfg.MongoDatabase)
	notificationManager := mongo.NewNotificationManager(mongoDB)

	t.Run("AddNotification", func(t *testing.T) {
		req := &carwash.AddNotificationRequest{
			BookingId: "test-booking-id",
			Message:   "This is a test notification",
			IsRead:    false,
		}
		resp, err := notificationManager.AddNotification(req)
		assert.NoError(t, err, "AddNotification funksiyasida xato yuz berdi")
		assert.NotNil(t, resp, "AddNotification javobi nol bo'lmasligi kerak")
	})

	t.Run("GetNotifications", func(t *testing.T) {
		req := &carwash.GetNotificationsRequest{
			BookingId: "test-booking-id",
		}
		resp, err := notificationManager.GetNotifications(req)
		assert.NoError(t, err, "GetNotifications funksiyasida xato yuz berdi")
		assert.NotNil(t, resp, "GetNotifications javobi nol bo'lmasligi kerak")
		assert.Greater(t, len(resp.Notifications), 0, "GetNotifications javobi bo'sh bo'lmasligi kerak")
	})

	t.Run("MarkNotificationAsRead", func(t *testing.T) {
		// Birinchi testni qo'shib olish
		reqAdd := &carwash.AddNotificationRequest{
			BookingId: "test-booking-id",
			Message:   "Test notification for marking as read",
			IsRead:    false,
		}
		respAdd, err := notificationManager.AddNotification(reqAdd)
		assert.NoError(t, err, "AddNotification funksiyasida xato yuz berdi")
		assert.NotNil(t, respAdd, "AddNotification javobi nol bo'lmasligi kerak")

		// Notificationni o'qilgan sifatida belgilash
		notificationsResp, err := notificationManager.GetNotifications(&carwash.GetNotificationsRequest{BookingId: "test-booking-id"})
		assert.NoError(t, err, "GetNotifications funksiyasida xato yuz berdi")
		assert.Greater(t, len(notificationsResp.Notifications), 0, "GetNotifications javobi bo'sh bo'lmasligi kerak")
		notificationID := notificationsResp.Notifications[0].Id

		reqMark := &carwash.MarkNotificationAsReadRequest{
			Id: notificationID,
		}
		respMark, err := notificationManager.MarkNotificationAsRead(reqMark)
		assert.NoError(t, err, "MarkNotificationAsRead funksiyasida xato yuz berdi")
		assert.True(t, respMark.Success, "Notificationni o'qilgan sifatida belgilash muvaffaqiyatsiz bo'ldi")

		// O'qilgan notificationni qayta olish
		notificationsResp, err = notificationManager.GetNotifications(&carwash.GetNotificationsRequest{BookingId: "test-booking-id"})
		assert.NoError(t, err, "GetNotifications funksiyasida xato yuz berdi")
		assert.NotNil(t, notificationsResp, "GetNotifications javobi nol bo'lmasligi kerak")
		assert.Equal(t, true, notificationsResp.Notifications[0].IsRead, "Notification hali o'qilgan deb belgilanmagan")
	})
}
