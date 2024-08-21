package mongo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Car-Wash/Car-Wash-Booking-Service/genproto/carwash"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/exp/rand"
)

type PaymentRepo struct {
	collection *mongo.Collection
}

func NewPaymentRepo(db *mongo.Database) *PaymentRepo {
	return &PaymentRepo{
		collection: db.Collection("payments"),
	}
}

func (r *PaymentRepo) CreatePayment(req *carwash.CreatePaymentRequest) (*carwash.CreatePaymentResponse, error) {
	req.TransactionId = fmt.Sprintf("%011d", rand.Intn(1000000000000))
	payment := bson.M{
		"booking_id":     req.BookingId,
		"amount":         req.Amount,
		"status":         req.Status,
		"payment_method": req.PaymentMethod,
		"transaction_id": req.TransactionId,
		"created_at":     time.Now().Format(time.RFC3339),
		"updated_at":     time.Now().Format(time.RFC3339),
	}

	_, err := r.collection.InsertOne(context.TODO(), payment)
	if err != nil {
		return nil, err
	}

	return &carwash.CreatePaymentResponse{}, nil
}

func (r *PaymentRepo) GetPayment(req *carwash.GetPaymentRequest) (*carwash.GetPaymentResponse, error) {
	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": id}

	var result bson.M
	err = r.collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("Payment not found")
		}
		return nil, err
	}

	paymentID, ok := result["_id"].(primitive.ObjectID)
	if !ok {
		return nil, errors.New("Failed to convert _id to ObjectID")
	}

	payment := &carwash.Payment{
		Id:            paymentID.Hex(),
		BookingId:     getString(result["booking_id"]),
		Amount:        getFloat32(result["amount"]),
		Status:        getString(result["status"]),
		PaymentMethod: getString(result["payment_method"]),
		TransactionId: getString(result["transaction_id"]),
		CreatedAt:     getString(result["created_at"]),
		UpdatedAt:     getString(result["updated_at"]),
	}

	return &carwash.GetPaymentResponse{
		Payment: payment,
	}, nil
}

func (r *PaymentRepo) ListPayments(req *carwash.ListPaymentsRequest) (*carwash.ListPaymentsResponse, error) {

	filter := bson.M{}
	if req.BookingId != "" {
		filter["booking_id"] = req.BookingId
	}
	if req.Amount != 0 {
		filter["amount"] = req.Amount
	}
	if req.Status != "" {
		filter["status"] = req.Status
	}
	if req.PaymentMethod != "" {
		filter["payment_method"] = req.PaymentMethod
	}
	if req.TransactionId != "" {
		filter["transaction_id"] = req.TransactionId
	}

	options := options.Find()
	if req.Limit != 0 {
		options.SetLimit(int64(req.Limit))
	}
	if req.Offset != 0 {
		options.SetSkip(int64(req.Offset))
	}

	cursor, err := r.collection.Find(context.TODO(), filter, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var payments []*carwash.Payment
	for cursor.Next(context.TODO()) {
		var result bson.M
		err := cursor.Decode(&result)
		if err != nil {
			return nil, err
		}

		paymentID, ok := result["_id"].(primitive.ObjectID)
		if !ok {
			return nil, errors.New("Failed to convert _id to ObjectID")
		}

		payment := &carwash.Payment{
			Id:            paymentID.Hex(),
			BookingId:     getString(result["booking_id"]),
			Amount:        float32(result["amount"].(float64)),
			Status:        getString(result["status"]),
			PaymentMethod: getString(result["payment_method"]),
			TransactionId: getString(result["transaction_id"]),
			CreatedAt:     getString(result["created_at"]),
			UpdatedAt:     getString(result["updated_at"]),
		}
		payments = append(payments, payment)
	}

	return &carwash.ListPaymentsResponse{
		Payments: payments,
	}, nil
}

func (r *PaymentRepo) UpdatePayment(req *carwash.UpdatePaymentRequest) (*carwash.UpdatePaymentResponse, error) {
	filter := bson.M{"transaction_id": req.Id}
	update := bson.M{"$set": bson.M{"status": req.Status, "updated_at": time.Now()}}

	_, err := r.collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}

	return &carwash.UpdatePaymentResponse{}, nil

}
