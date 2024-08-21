package test

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/Car-Wash/Car-Wash-Booking-Service/config"
	"github.com/Car-Wash/Car-Wash-Booking-Service/genproto/carwash"
	"github.com/Car-Wash/Car-Wash-Booking-Service/storage/mongo"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	m "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestPaymentRepo(t *testing.T) {
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
	paymentRepo := mongo.NewPaymentRepo(mongoDB)

	t.Run("CreatePayment", func(t *testing.T) {
		req := &carwash.CreatePaymentRequest{
			BookingId:     uuid.NewString(),
			Amount:        100.0,
			Status:        "Pending",
			PaymentMethod: "Credit Card",
		}
		resp, err := paymentRepo.CreatePayment(req)
		assert.NoError(t, err, "CreatePayment funksiyasida xato yuz berdi")
		assert.NotNil(t, resp, "CreatePayment javobi nol bo'lmasligi kerak")
	})

	t.Run("GetPayment", func(t *testing.T) {
		// Create a payment to retrieve
		reqCreate := &carwash.CreatePaymentRequest{
			BookingId:     uuid.NewString(),
			Amount:        100.0,
			Status:        "Pending",
			PaymentMethod: "Credit Card",
		}
		createResp, err := paymentRepo.CreatePayment(reqCreate)
		assert.NoError(t, err, "CreatePayment funksiyasida xato yuz berdi")
		assert.NotNil(t, createResp, "CreatePayment javobi nol bo'lmasligi kerak")

		// Retrieve the payment
		filter := bson.M{"transaction_id": "1258885699"}
		var result bson.M
		err = mongoDB.Collection("payments").FindOne(context.TODO(), filter).Decode(&result)
		if err != nil {
			t.Fatalf("Failed to retrieve payment: %v", err)
		}
		paymentID, ok := result["_id"].(primitive.ObjectID)
		if !ok {
			t.Fatalf("Failed to convert _id to ObjectID")
		}

		reqGet := &carwash.GetPaymentRequest{Id: paymentID.Hex()}
		resp, err := paymentRepo.GetPayment(reqGet)
		assert.NoError(t, err, "GetPayment funksiyasida xato yuz berdi")
		assert.NotNil(t, resp, "GetPayment javobi nol bo'lmasligi kerak")
	})

	t.Run("ListPayments", func(t *testing.T) {
		// Create payments to list
		reqCreate1 := &carwash.CreatePaymentRequest{
			BookingId:     uuid.NewString(),
			Amount:        100.0,
			Status:        "Pending",
			PaymentMethod: "Credit Card",
		}
		reqCreate2 := &carwash.CreatePaymentRequest{
			BookingId:     uuid.NewString(),
			Amount:        200.0,
			Status:        "Completed",
			PaymentMethod: "PayPal",
		}
		_, err := paymentRepo.CreatePayment(reqCreate1)
		assert.NoError(t, err, "CreatePayment funksiyasida xato yuz berdi")
		_, err = paymentRepo.CreatePayment(reqCreate2)
		assert.NoError(t, err, "CreatePayment funksiyasida xato yuz berdi")

		req := &carwash.ListPaymentsRequest{}
		resp, err := paymentRepo.ListPayments(req)
		assert.NoError(t, err, "ListPayments funksiyasida xato yuz berdi")
		assert.Greater(t, len(resp.Payments), 1, "ListPayments javobi kerakli miqdordagi to'lovlarni qaytarishi kerak")
	})

}
