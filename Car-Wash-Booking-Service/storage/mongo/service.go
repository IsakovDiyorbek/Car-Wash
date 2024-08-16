package mongo

import (
	"context"
	"log/slog"
	"time"

	pb "github.com/exam-5/Car-Wash-Booking-Service/genproto/carwash"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ServiceRepo struct {
	mongo *mongo.Collection
}

func NewServiceRepo(db *mongo.Database) *ServiceRepo {
	return &ServiceRepo{
		mongo: db.Collection("services"),
	}
}

func (s *ServiceRepo) CreateService(req *pb.CreateServiceRequest) (*pb.CreateServiceResponse, error) {
	service := pb.Services{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Duration:    req.Duration,

		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}

	_, err := s.mongo.InsertOne(context.TODO(), service)
	if err != nil {
		slog.Info("Error creating service", err)
		return nil, err

	}
	return &pb.CreateServiceResponse{}, nil
}

func (s *ServiceRepo) GetService(req *pb.GetServiceRequest) (*pb.GetServiceResponse, error) {
	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		slog.Info("Invalid ObjectId format", err)
		return nil, err
	}
	filter := bson.M{"_id": id}
	var bson bson.M

	err = s.mongo.FindOne(context.TODO(), filter).Decode(&bson)
	if err != nil {
		slog.Info("Error getting service", err)
		return nil, err
	}

	service := &pb.Services{
		Id:          bson["_id"].(primitive.ObjectID).Hex(),
		Name:        bson["name"].(string),
		Description: bson["description"].(string),
		Price:       getFloat32(bson["price"]),
		Duration:    bson["duration"].(int32),
		CreatedAt:   bson["createdat"].(string),
		UpdatedAt:   bson["updatedat"].(string),
	}
	return &pb.GetServiceResponse{Service: service}, nil
}
func (s *ServiceRepo) UpdateService(req *pb.UpdateServiceRequest) (*pb.UpdateServiceResponse, error) {
	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		slog.Info("Invalid ObjectId format", err)
		return nil, err
	}

	filter := bson.M{"_id": id}

	updateFields := bson.M{}
	if req.Name != "" {
		updateFields["name"] = req.Name
	}
	if req.Description != "" {
		updateFields["description"] = req.GetDescription()
	}
	if req.Price != 0 {
		updateFields["price"] = req.Price
	}
	if req.Duration > 0 {
		updateFields["duration"] = req.Duration
	}

	if len(updateFields) > 0 {
		updateFields["updatedat"] = time.Now().Format(time.RFC3339)

		update := bson.M{
			"$set": updateFields,
		}

		_, err = s.mongo.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			slog.Info("Error updating service", err)
			return nil, err
		}
	}

	return &pb.UpdateServiceResponse{}, nil
}

func (s *ServiceRepo) DeleteService(req *pb.DeleteServiceRequest) (*pb.DeleteServiceResponse, error) {

	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		slog.Info("Invalid ObjectId format", err)
		return nil, err
	}

	filter := bson.M{"_id": id}

	_, err = s.mongo.DeleteOne(context.TODO(), filter)
	if err != nil {
		slog.Info("Error deleting service", err)
		return nil, err
	}

	return &pb.DeleteServiceResponse{}, nil
}
func (s *ServiceRepo) ListServices(req *pb.ListServicesRequest) (*pb.ListServicesResponse, error) {
	filter := bson.M{}

	if req.Name != "" {
		filter["name"] = bson.M{"$regex": req.Name, "$options": "i"}
	}
	if req.Description != "" {
		filter["description"] = bson.M{"$regex": req.Description, "$options": "i"}
	}
	if req.Price > 0 {
		filter["price"] = req.Price
	}
	if req.Duration > 0 {
		filter["duration"] = req.Duration
	}

	findOptions := options.Find()
	findOptions.SetLimit(int64(req.Limit))
	findOptions.SetSkip(int64(req.Offset))

	cursor, err := s.mongo.Find(context.TODO(), filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var services []*pb.Services
	for cursor.Next(context.TODO()) {
		var res bson.M
		if err := cursor.Decode(&res); err != nil {
			return nil, err
		}

		service := &pb.Services{
			Id:        res["_id"].(primitive.ObjectID).Hex(),
			Name:      res["name"].(string),
			Price:     getFloat32(res["price"]),
			Duration:  res["duration"].(int32),
			CreatedAt: res["createdat"].(string),
			UpdatedAt: res["updatedat"].(string),
		}
		services = append(services, service)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return &pb.ListServicesResponse{
		Services: services,
	}, nil
}

func (s *ServiceRepo) SearchServices(req *pb.SearchServicesRequest) (*pb.SearchServicesResponse, error) {
	var services []*pb.Services

	filter := bson.M{}
	if req.Name != "" {
		filter["name"] = req.Name
	}
	if req.Description != "" {
		filter["description"] = req.Description
	}

	cursor, err := s.mongo.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var res bson.M
		if err := cursor.Decode(&res); err != nil {
			return nil, err
		}

		service := pb.Services{
			Id:        res["_id"].(primitive.ObjectID).Hex(),
			Name:      res["name"].(string),
			Price:     getFloat32(res["price"]),
			Duration:  res["duration"].(int32),
			CreatedAt: res["createdat"].(string),
			UpdatedAt: res["updatedat"].(string),
		}

		services = append(services, &service)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return &pb.SearchServicesResponse{Services: services}, nil
}
