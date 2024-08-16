package mongo

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	pb "github.com/exam-5/Car-Wash-Booking-Service/genproto/carwash"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type NotificationManager struct {
	collec *mongo.Collection
}

func NewNotificationManager(db *mongo.Database) *NotificationManager {

	return &NotificationManager{
		collec: db.Collection("notifications"),
	}
}

func (n *NotificationManager) AddNotification(req *pb.AddNotificationRequest) (*pb.AddNotificationResponse, error) {
	notification := bson.M{
		"booking_id": req.BookingId,
		"message":    req.Message,
		"is_read":    req.IsRead,
		"created_at": time.Now().Format(time.RFC3339),
		"updated_at": time.Now().Format(time.RFC3339),
	}
	_, err := n.collec.InsertOne(context.Background(), notification)
	if err != nil {
		log.Println("error while inserting", err)
		return nil, err

	}

	return &pb.AddNotificationResponse{}, nil

}

func (n *NotificationManager) GetNotifications(req *pb.GetNotificationsRequest) (*pb.GetNotificationsResponse, error) {
	log.Printf("Received request with BookingId: %s", req.BookingId)

	filter := bson.M{}
	if req.BookingId != "" {
		filter["booking_id"] = req.BookingId
	}

	var notifications []*pb.Notification
	cursor, err := n.collec.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var res bson.M
		if err := cursor.Decode(&res); err != nil {
			return nil, err
		}
		NotificationID, ok := res["_id"].(primitive.ObjectID)
		if !ok {
			return nil, errors.New("Failed to convert _id to ObjectID")
		}

		notification := pb.Notification{
			Id:        NotificationID.Hex(),
			BookingId: getString(res["booking_id"]),
			Message:   getString(res["message"]),
			IsRead:    getBoolean(res["is_read"]),
			CreatedAt: getString(res["created_at"]),
		}
		notifications = append(notifications, &notification)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return &pb.GetNotificationsResponse{Notifications: notifications}, nil
}

func (n *NotificationManager) MarkNotificationAsRead(req *pb.MarkNotificationAsReadRequest) (*pb.MarkNotificationAsReadResponse, error) {

	objectId, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, fmt.Errorf("invalid resource id format: %s", req.Id)
	}

	filter := bson.M{"_id": objectId}
	update := bson.M{"$set": bson.M{"is_read": true}}

	result, err := n.collec.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}

	if result.MatchedCount == 0 {
		return nil, fmt.Errorf("resource not found")
	}

	response := &pb.MarkNotificationAsReadResponse{
		Message: "Resource marked as complete",
		Success: true,
	}

	return response, nil
}
