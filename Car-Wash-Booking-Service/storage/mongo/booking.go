package mongo

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	pb "github.com/Car-Wash/Car-Wash-Booking-Service/genproto/carwash"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BookingRepo struct {
	collection *mongo.Collection
	redis      *redis.Client
}

func NewBookingRepo(db *mongo.Database, redis *redis.Client) *BookingRepo {
	return &BookingRepo{
		collection: db.Collection("bookings"),
		redis:      redis,
	}
}

func (r *BookingRepo) CreateBooking(req *pb.CreateBookingRequest) (*pb.CreateBookingResponse, error) {
	booking := bson.M{
		"user_id":        req.UserId,
		"provider_id":    req.ProviderId,
		"service_id":     req.ServiceId,
		"status":         req.Status,
		"scheduled_time": req.ScheduledTime,
		"location":       bson.M{"latitude": req.Location.Latitude, "longitude": req.Location.Longitude},
		"total_price":    req.TotalPrice,
		"created_at":     time.Now().Format(time.RFC3339),
		"updated_at":     time.Now().Format(time.RFC3339),
	}

	_, err := r.collection.InsertOne(context.TODO(), booking)
	if err != nil {

		slog.Error("Error creating booking", err)
		return nil, err
	}

	_, err = r.redis.ZIncrBy(context.TODO(), "popular_services", 1, req.ServiceId).Result()
	if err != nil {
		slog.Error("Error updating Redis ZSET", err)
		return nil, err
	}

	return &pb.CreateBookingResponse{}, nil
}

func (r *BookingRepo) GetBooking(req *pb.GetBookingRequest) (*pb.GetBookingResponse, error) {
	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		slog.Error("Invalid ObjectId format", err)
		return nil, err
	}
	filter := bson.M{"_id": id}

	var result bson.M
	err = r.collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("Booking not found")
		}
		return nil, err
	}
	reviewID, ok := result["_id"].(primitive.ObjectID)
	if !ok {
		return nil, errors.New("Failed to convert _id to ObjectID")
	}
	booking := &pb.Booking{
		Id:            reviewID.Hex(),
		UserId:        safeString(result["user_id"]),
		ProviderId:    safeString(result["provider_id"]),
		ServiceId:     safeString(result["service_id"]),
		Status:        safeString(result["status"]),
		ScheduledTime: safeString(result["scheduled_time"]),
		Location: &pb.GeoPoint{
			Latitude:  safeFloat64(result["location"].(bson.M)["latitude"]),
			Longitude: safeFloat64(result["location"].(bson.M)["longitude"]),
		},
		TotalPrice: safeFloat32(result["total_price"]),
		CreatedAt:  safeString(result["created_at"]),
		UpdatedAt:  safeString(result["updated_at"]),
	}

	return &pb.GetBookingResponse{
		Booking: booking,
	}, nil
}

func (r *BookingRepo) DeleteBooking(req *pb.DeleteBookingRequest) (*pb.DeleteBookingResponse, error) {
	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		slog.Error("Invalid ObjectId format", err)
		return nil, err
	}
	filter := bson.M{"_id": id}

	_, err = r.collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		slog.Error("Error deleting booking", err)
		return nil, err
	}

	return &pb.DeleteBookingResponse{}, nil
}
func (r *BookingRepo) ListBookings(req *pb.ListBookingsRequest) (*pb.ListBookingsResponse, error) {
	filter := bson.M{}

	if req.UserId != "" {
		filter["user_id"] = req.UserId
	}
	if req.ProviderId != "" {
		filter["provider_id"] = req.ProviderId
	}
	if req.ServiceId != "" {
		filter["service_id"] = req.ServiceId
	}
	if req.Status != "" {
		filter["status"] = req.Status
	}
	if req.ScheduledTime != "" {
		filter["scheduled_time"] = req.ScheduledTime
	}

	findOptions := options.Find()
	findOptions.SetLimit(int64(req.Limit))
	findOptions.SetSkip(int64(req.Offset))

	cursor, err := r.collection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var bookings []*pb.Booking
	for cursor.Next(context.TODO()) {
		var res bson.M
		if err := cursor.Decode(&res); err != nil {
			return nil, err
		}

		bookingID, ok := res["_id"].(primitive.ObjectID)
		if !ok {
			return nil, errors.New("Failed to convert _id to ObjectID")
		}
		booking := &pb.Booking{
			Id:            bookingID.Hex(),
			UserId:        getString(res["user_id"]),
			ProviderId:    getString(res["provider_id"]),
			ServiceId:     getString(res["service_id"]),
			Status:        getString(res["status"]),
			ScheduledTime: getString(res["scheduled_time"]),
			Location: &pb.GeoPoint{
				Latitude:  GetFloat64(res["location"], "latitude"),
				Longitude: GetFloat64(res["location"], "longitude"),
			},
			TotalPrice: getFloat32(res["total_price"]),
			CreatedAt:  getString(res["created_at"]),
			UpdatedAt:  getString(res["updated_at"]),
		}

		bookings = append(bookings, booking)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return &pb.ListBookingsResponse{
		Bookings: bookings,
	}, nil
}

func (r *BookingRepo) UpdateBooking(req *pb.UpdateBookingRequest) (*pb.UpdateBookingResponse, error) {
	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, fmt.Errorf("invalid ObjectId format: %v", err)
	}

	filter := bson.M{"_id": id}

	updateFields := bson.M{}

	if req.UserId != "" {
		updateFields["user_id"] = req.UserId
	}
	if req.ProviderId != "" {
		updateFields["provider_id"] = req.ProviderId
	}
	if req.ServiceId != "" {
		updateFields["service_id"] = req.ServiceId
	}
	if req.Status != "" {
		updateFields["status"] = req.Status
	}
	if req.ScheduledTime != "" {
		updateFields["scheduled_time"] = req.ScheduledTime
	}
	if req.Location != nil {
		updateFields["location"] = bson.M{
			"latitude":  req.Location.Latitude,
			"longitude": req.Location.Longitude,
		}
	}
	if req.TotalPrice > 0 {
		updateFields["total_price"] = req.TotalPrice
	}

	if len(updateFields) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	updateFields["updated_at"] = time.Now().Format(time.RFC3339)
	update := bson.M{
		"$set": updateFields,
	}

	_, err = r.collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, fmt.Errorf("error updating booking: %v", err)
	}

	return &pb.UpdateBookingResponse{}, nil
}
