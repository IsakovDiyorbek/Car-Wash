package mongo

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	pb "github.com/exam-5/Car-Wash-Booking-Service/genproto/carwash"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ReviewRepo struct {
	collection *mongo.Collection
}

func NewReviewRepo(db *mongo.Database) *ReviewRepo {
	return &ReviewRepo{
		collection: db.Collection("reviews"),
	}
}

func (r *ReviewRepo) CreateReview(req *pb.CreateReviewRequest) (*pb.CreateReviewResponse, error) {
	review := bson.M{
		"booking_id":  req.BookingId,
		"user_id":     req.UserId,
		"provider_id": req.ProviderId,
		"rating":      req.Rating,
		"comment":     req.Comment,
		"created_at":  time.Now().Format(time.RFC3339),
		"updated_at":  time.Now().Format(time.RFC3339),
	}

	_, err := r.collection.InsertOne(context.TODO(), review)
	if err != nil {
		slog.Error("Error creating review", err)
		return nil, err
	}
	return &pb.CreateReviewResponse{}, nil
}

func (r *ReviewRepo) GetReview(req *pb.GetReviewRequest) (*pb.GetReviewResponse, error) {
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
			return nil, errors.New("Review not found")
		}
		return nil, err
	}


	reviewID, ok := result["_id"].(primitive.ObjectID)
	if !ok {
		return nil, errors.New("Failed to convert _id to ObjectID")
	}



	review := &pb.Review{
		Id:         reviewID.Hex(), 
		BookingId:  getString(result["booking_id"]),
		UserId:     getString(result["user_id"]),
		ProviderId: getString(result["provider_id"]),
		Rating:     getFloat32(result["rating"]),
		Comment:    getString(result["comment"]),
		CreatedAt:  getString(result["created_at"]),
		UpdatedAt:  getString(result["updated_at"]),
	}

	return &pb.GetReviewResponse{
		Review: review,
	}, nil
}





func (r *ReviewRepo) DeleteReview(req *pb.DeleteReviewRequest) (*pb.DeleteReviewResponse, error) {
	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		slog.Error("Invalid ObjectId format", err)
		return nil, err
	}
	filter := bson.M{"_id": id}

	_, err = r.collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		slog.Error("Error deleting review", err)
		return nil, err
	}

	return &pb.DeleteReviewResponse{}, nil
}

func (r *ReviewRepo) ListReviews(req *pb.ListReviewsRequest) (*pb.ListReviewsResponse, error) {
	filter := bson.M{}

	if req.BookingId != "" {
		filter["booking_id"] = req.BookingId
	}
	if req.ProviderId != "" {
		filter["provider_id"] = req.ProviderId
	}
	if req.Rating != 0 {
		filter["rating"] = req.Rating
	}
	if req.Comment != "" {
		filter["comment"] = bson.M{"$regex": req.Comment, "$options": "i"}
	}

	findOptions := options.Find()
	findOptions.SetLimit(int64(req.Limit))
	findOptions.SetSkip(int64(req.Offset))

	cursor, err := r.collection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var reviews []*pb.Review
	for cursor.Next(context.TODO()) {
		var res bson.M
		if err := cursor.Decode(&res); err != nil {
			return nil, err
		}
		
		reviewID, ok := res["_id"].(primitive.ObjectID)
		if !ok {
			return nil, errors.New("Failed to convert _id to ObjectID")
		}


		review := &pb.Review{
			Id:         reviewID.Hex(),
			BookingId:  getString(res["booking_id"]),
			UserId:     getString(res["user_id"]),
			ProviderId: getString(res["provider_id"]),
			Rating:     getFloat32(res["rating"]),
			Comment:    getString(res["comment"]),
			CreatedAt:  getString(res["created_at"]),
			UpdatedAt:  getString(res["updated_at"]),
		}

		reviews = append(reviews, review)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return &pb.ListReviewsResponse{
		Reviews: reviews,
	}, nil
}

func (r *ReviewRepo) UpdateReview(req *pb.UpdateReviewRequest) (*pb.UpdateReviewResponse, error) {
	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, fmt.Errorf("invalid ObjectId format: %v", err)
	}

	filter := bson.M{"_id": id}

	updateFields := bson.M{}

	if req.BookingId != "" {
		updateFields["booking_id"] = req.BookingId
	}
	if req.UserId != "" {
		updateFields["user_id"] = req.UserId
	}
	if req.ProviderId != "" {
		updateFields["provider_id"] = req.ProviderId
	}
	if req.Rating != 0 {
		updateFields["rating"] = req.Rating
	}
	if req.Comment != "" {
		updateFields["comment"] = req.Comment
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
		return nil, fmt.Errorf("error updating review: %v", err)
	}

	return &pb.UpdateReviewResponse{}, nil
}





