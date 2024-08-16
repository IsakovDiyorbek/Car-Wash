package mongo

import (
	"context"
	"errors"
	"time"

	"log/slog"

	pb "github.com/exam-5/Car-Wash-Booking-Service/genproto/carwash"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProviderRepo struct {
	collection *mongo.Collection
}

func NewProviderRepo(db *mongo.Database) *ProviderRepo {
	return &ProviderRepo{
		collection: db.Collection("providers"),
	}
}
func (r *ProviderRepo) CreateProvider(req *pb.CreateProviderRequest) (*pb.CreateProviderResponse, error) {
	var availabilityBsonArray []bson.M
	if req.Availability != nil {
		for _, avail := range req.Availability {
			if avail != nil {
				availabilityBsonArray = append(availabilityBsonArray, bson.M{
					"day":        avail.Day,
					"start_time": avail.StartTime,
					"end_time":   avail.EndTime,
				})
			}
		}
	}

	provider := bson.M{
		"user_id":        req.UserId,
		"company_name":   req.CompanyName,
		"description":    req.Description,
		"service_id":     req.ServiceId,
		"availability":   availabilityBsonArray,
		"location":       req.Location,
		"average_rating": req.AverageRating,
		"created_at":     time.Now().Format(time.RFC3339),
		"updated_at":     time.Now().Format(time.RFC3339),
	}

	_, err := r.collection.InsertOne(context.TODO(), provider)
	if err != nil {
		slog.Error("Error creating provider", err)
		return nil, err
	}




	return &pb.CreateProviderResponse{}, nil
}

func (r *ProviderRepo) GetProvider(req *pb.GetProviderRequest) (*pb.GetProviderResponse, error) {
	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		slog.Error("Invalid ObjectId format", err)
		return nil, err
	}

	var res bson.M
	filter := bson.M{"_id": id}

	err = r.collection.FindOne(context.TODO(), filter).Decode(&res)
	if err != nil {
		slog.Error("Error getting provider", err)
		return nil, err
	}

	getStringArray := func(data interface{}) []string {
		if arr, ok := data.(primitive.A); ok {
			strArr := make([]string, len(arr))
			for i, v := range arr {
				if s, ok := v.(string); ok {
					strArr[i] = s
				}
			}
			return strArr
		}
		return []string{}
	}
	providerID, ok := res["_id"].(primitive.ObjectID)
	if !ok {
		return nil, errors.New("Failed to convert _id to ObjectID")
	}

	provider := pb.Provider{
		Id:            providerID.Hex(),
		UserId:        getString(res["user_id"]),
		CompanyName:   getString(res["company_name"]),
		Description:   getString(res["description"]),
		ServiceId:     getStringArray(res["service_id"]),
		AverageRating: float32(getFloat64(res["average_rating"])),
		CreatedAt:     getString(res["created_at"]),
		UpdatedAt:     getString(res["updated_at"]),
	}

	if availability, ok := res["availability"].(primitive.A); ok {
		for _, a := range availability {
			if av, ok := a.(bson.M); ok {
				provider.Availability = append(provider.Availability, &pb.Availability{
					Day:       getString(av["day"]),
					StartTime: getString(av["start_time"]),
					EndTime:   getString(av["end_time"]),
				})
			}
		}
	}

	if location, ok := res["location"].(bson.M); ok {
		provider.Location = &pb.GeoPoing{
			Latitude:  getFloat64(location["latitude"]),
			Longitude: getFloat64(location["longitude"]),
		}
	}

	return &pb.GetProviderResponse{Provider: &provider}, nil
}

func (r *ProviderRepo) UpdateProvider(req *pb.UpdateProviderRequest) (*pb.UpdateProviderResponse, error) {
	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		slog.Error("Invalid ObjectId format", err)
		return nil, err
	}

	filter := bson.M{"_id": id}

	updateFields := bson.M{}
	if req.UserId != "" {
		updateFields["user_id"] = req.UserId
	}
	if req.CompanyName != "" {
		updateFields["company_name"] = req.CompanyName
	}
	if req.Description != "" {
		updateFields["description"] = req.Description
	}
	if len(req.ServiceId) > 0 {
		updateFields["service_id"] = req.ServiceId
	}
	if len(req.Availability) > 0 {
		updateFields["availability"] = req.Availability
	}
	if req.Location != nil {
		updateFields["location"] = req.Location
	}

	if req.AverageRating > 0 {
		updateFields["average_rating"] = req.AverageRating
	}

	if len(updateFields) > 0 {
		updateFields["updated_at"] = time.Now().Format(time.RFC3339)
		update := bson.M{
			"$set": updateFields,
		}

		_, err = r.collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			slog.Error("Error updating provider", err)
			return nil, err
		}
	}

	return &pb.UpdateProviderResponse{}, nil
}

func (r *ProviderRepo) DeleteProvider(req *pb.DeleteProviderRequest) (*pb.DeleteProviderResponse, error) {
	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		slog.Error("Invalid ObjectId format", err)
		return nil, err
	}

	filter := bson.M{"_id": id}

	_, err = r.collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		slog.Error("Error deleting provider", err)
		return nil, err
	}

	return &pb.DeleteProviderResponse{}, nil
}
func (r *ProviderRepo) ListProviders(req *pb.ListProvidersRequest) (*pb.ListProvidersResponse, error) {
	filter := bson.M{}

	if req.UserId != "" {
		filter["user_id"] = req.UserId
	}
	if req.CompanyName != "" {
		filter["company_name"] = req.CompanyName
	}
	if req.Description != "" {
		filter["description"] = req.Description
	}
	if req.AverageRating > 0 {
		filter["average_rating"] = bson.M{"$gte": req.AverageRating}
	}

	findOptions := options.Find()
	findOptions.SetLimit(int64(req.Limit))
	findOptions.SetSkip(int64(req.Offset))

	cursor, err := r.collection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var providers []*pb.Provider
	for cursor.Next(context.TODO()) {
		var res bson.M
		if err := cursor.Decode(&res); err != nil {
			return nil, err
		}

		getStringArray := func(data interface{}) []string {
			if arr, ok := data.(primitive.A); ok {
				strArr := make([]string, len(arr))
				for i, v := range arr {
					if s, ok := v.(string); ok {
						strArr[i] = s
					}
				}
				return strArr
			}
			return []string{}
		}

		reviewID, ok := res["_id"].(primitive.ObjectID)
		if !ok {
			return nil, errors.New("Failed to convert _id to ObjectID")
		}

		provider := pb.Provider{
			Id:            reviewID.Hex(),
			UserId:        getString(res["user_id"]),
			CompanyName:   getString(res["company_name"]),
			Description:   getString(res["description"]),
			ServiceId:     getStringArray(res["service_id"]),
			AverageRating: float32(getFloat64(res["average_rating"])),
			CreatedAt:     getString(res["created_at"]),
			UpdatedAt:     getString(res["updated_at"]),
		}

		if availability, ok := res["availability"].(primitive.A); ok {
			for _, a := range availability {
				if av, ok := a.(bson.M); ok {
					provider.Availability = append(provider.Availability, &pb.Availability{
						Day:       getString(av["day"]),
						StartTime: getString(av["start_time"]),
						EndTime:   getString(av["end_time"]),
					})
				}
			}
		}

		if location, ok := res["location"].(bson.M); ok {
			provider.Location = &pb.GeoPoing{
				Latitude:  getFloat64(location["latitude"]),
				Longitude: getFloat64(location["longitude"]),
			}
		}

		providers = append(providers, &provider)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return &pb.ListProvidersResponse{
		Provider: providers,
	}, nil
}


func (p *ProviderRepo) SearchProviders(req *pb.SearchProvidersRequest) (*pb.SearchProvidersResponse, error) {

	var providers []*pb.Provider

	filter := bson.M{}
	if req.CompanyName != "" {
		filter["company_name"] = req.CompanyName
	}
	if req.Description != "" {
		filter["description"] = req.Description
	}

	cursor, err := p.collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var res bson.M
		if err := cursor.Decode(&res); err != nil {
			return nil, err
		}
		getStringArray := func(data interface{}) []string {
			if arr, ok := data.(primitive.A); ok {
				strArr := make([]string, len(arr))
				for i, v := range arr {
					if s, ok := v.(string); ok {
						strArr[i] = s
					}
				}
				return strArr
			}
			return []string{}
		}

		provider := pb.Provider{
			Id:            res["_id"].(primitive.ObjectID).Hex(),
			UserId:        getString(res["user_id"]),
			CompanyName:   getString(res["company_name"]),
			Description:   getString(res["description"]),
			ServiceId:     getStringArray(res["service_id"]),
			AverageRating: float32(getFloat64(res["average_rating"])),
			CreatedAt:     getString(res["created_at"]),
			UpdatedAt:     getString(res["updated_at"]),
		}

		providers = append(providers, &provider)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}


	return &pb.SearchProvidersResponse{Providers: providers}, nil
}