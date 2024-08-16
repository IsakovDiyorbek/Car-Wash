package test

import (
	"context"
	"testing"

	"github.com/exam-5/Car-Wash-Booking-Service/genproto/carwash"
	m "github.com/exam-5/Car-Wash-Booking-Service/storage/mongo"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setup(t *testing.T) (*mongo.Database, *mongo.Collection, *mongo.Client) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	db := client.Database("car_wash_test")
	collection := db.Collection("providers")

	return db, collection, client
}

func teardown(t *testing.T, client *mongo.Client) {
	err := client.Disconnect(context.TODO())
	if err != nil {
		t.Fatalf("Failed to disconnect from MongoDB: %v", err)
	}
}

func TestCreateProvider(t *testing.T) {
	db, collection, client := setup(t)
	defer teardown(t, client)

	repo := m.NewProviderRepo(db)

	req := &carwash.CreateProviderRequest{
		UserId:      "user1",
		CompanyName: "Company A",
		Description: "Description A",
		ServiceId:   []string{"service1"},
		Availability: []*carwash.Availability{
			{Day: "Monday", StartTime: "09:00", EndTime: "17:00"},
		},
		Location:      &carwash.GeoPoing{Latitude: 40.7128, Longitude: -74.0060},
		AverageRating: 4.5,
	}

	resp, err := repo.CreateProvider(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	// Verify the provider is created
	filter := bson.M{"user_id": req.UserId}
	var result bson.M
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	assert.NoError(t, err)
	assert.Equal(t, req.UserId, result["user_id"])
}

func TestGetProvider(t *testing.T) {
	db, collection, client := setup(t)
	defer teardown(t, client)

	repo := m.NewProviderRepo(db)

	// Create a provider to retrieve
	createReq := &carwash.CreateProviderRequest{
		UserId:      "user2",
		CompanyName: "Company B",
		Description: "Description B",
		ServiceId:   []string{"service2"},
		Availability: []*carwash.Availability{
			{Day: "Tuesday", StartTime: "10:00", EndTime: "18:00"},
		},
		Location:      &carwash.GeoPoing{Latitude: 34.0522, Longitude: -118.2437},
		AverageRating: 4.7,
	}

	createResp, err := repo.CreateProvider(createReq)
	assert.NoError(t, err)
	assert.NotNil(t, createResp)

	// Retrieve the provider
	filter := bson.M{"user_id": createReq.UserId}
	var result bson.M
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	assert.NoError(t, err)
	providerID, _ := result["_id"].(primitive.ObjectID)
	getReq := &carwash.GetProviderRequest{Id: providerID.Hex()}
	getResp, err := repo.GetProvider(getReq)
	assert.NoError(t, err)
	assert.NotNil(t, getResp)
	assert.Equal(t, createReq.UserId, getResp.Provider.UserId)
}

func TestUpdateProvider(t *testing.T) {
	db, collection, client := setup(t)
	defer teardown(t, client)

	repo := m.NewProviderRepo(db)

	// Create a provider to update
	createReq := &carwash.CreateProviderRequest{
		UserId:      "user3",
		CompanyName: "Company C",
		Description: "Description C",
		ServiceId:   []string{"service3"},
		Availability: []*carwash.Availability{
			{Day: "Wednesday", StartTime: "08:00", EndTime: "16:00"},
		},
		Location:      &carwash.GeoPoing{Latitude: 37.7749, Longitude: -122.4194},
		AverageRating: 4.8,
	}

	createResp, err := repo.CreateProvider(createReq)
	assert.NoError(t, err)
	assert.NotNil(t, createResp)

	// Update the provider

	updateReq := &carwash.UpdateProviderRequest{
		Id:            "66bd9994b671b4f10b3ff5cq",
		Description:   "Updated Description",
		AverageRating: 4.9,
	}

	updateResp, err := repo.UpdateProvider(updateReq)
	assert.NoError(t, err)
	assert.NotNil(t, updateResp)

	// Verify the provider is updated
	filter := bson.M{"_id": "66bd91294b671b4f10b3ff5cd"}
	var result bson.M
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	assert.NoError(t, err)
	assert.Equal(t, updateReq.Description, result["description"])
	assert.Equal(t, updateReq.AverageRating, result["average_rating"])
}

func TestDeleteProvider(t *testing.T) {
	db, _, client := setup(t)
	defer teardown(t, client)

	repo := m.NewProviderRepo(db)

	// Create a provider to delete
	createReq := &carwash.CreateProviderRequest{
		UserId:      "user4",
		CompanyName: "Company D",
		Description: "Description D",
		ServiceId:   []string{"service4"},
		Availability: []*carwash.Availability{
			{Day: "Thursday", StartTime: "11:00", EndTime: "19:00"},
		},
		Location:      &carwash.GeoPoing{Latitude: 40.7306, Longitude: -73.9352},
		AverageRating: 5.0,
	}

	createResp, err := repo.CreateProvider(createReq)
	assert.NoError(t, err)
	assert.NotNil(t, createResp)

	// Delete the provider

	deleteReq := &carwash.DeleteProviderRequest{Id: "66bd9994b671b4f10b3ff5cd"}
	deleteResp, err := repo.DeleteProvider(deleteReq)
	assert.NoError(t, err)
	assert.NotNil(t, deleteResp)

	assert.Error(t, err)
}

func TestListProviders(t *testing.T) {
	db, _, client := setup(t)
	defer teardown(t, client)

	repo := m.NewProviderRepo(db)

	// Create providers to list
	for i := 0; i < 5; i++ {
		createReq := &carwash.CreateProviderRequest{
			UserId:      "user" + string(i),
			CompanyName: "Company " + string(i),
			Description: "Description " + string(i),
			ServiceId:   []string{"service" + string(i)},
			Availability: []*carwash.Availability{
				{Day: "Friday", StartTime: "09:00", EndTime: "17:00"},
			},
			Location:      &carwash.GeoPoing{Latitude: 42.3601, Longitude: -71.0589},
			AverageRating: 4.0 + float32(i),
		}
		_, err := repo.CreateProvider(createReq)
		assert.NoError(t, err)
	}

	// List providers
	listReq := &carwash.ListProvidersRequest{
		Limit:  5,
		Offset: 0,
	}
	listResp, err := repo.ListProviders(listReq)
	assert.NoError(t, err)
	assert.NotNil(t, listResp)
	assert.Len(t, listResp.Provider, 5)
}

func TestSearchProviders(t *testing.T) {
	db, _, client := setup(t)
	defer teardown(t, client)

	repo := m.NewProviderRepo(db)

	// Create providers to search
	for i := 0; i < 5; i++ {
		createReq := &carwash.CreateProviderRequest{
			UserId:      "user" + string(i),
			CompanyName: "Company " + string(i),
			Description: "Description " + string(i),
			ServiceId:   []string{"service" + string(i)},
			Availability: []*carwash.Availability{
				{Day: "Saturday", StartTime: "10:00", EndTime: "16:00"},
			},
			Location:      &carwash.GeoPoing{Latitude: 43.0731, Longitude: -89.4012},
			AverageRating: 4.5 + float32(i),
		}
		_, err := repo.CreateProvider(createReq)
		assert.NoError(t, err)
	}

	// Search providers
	searchReq := &carwash.SearchProvidersRequest{
		CompanyName: "Company 1",
	}
	searchResp, err := repo.SearchProviders(searchReq)
	assert.NoError(t, err)
	assert.NotNil(t, searchResp)
	assert.Len(t, searchResp.Providers, 1)
	assert.Equal(t, "Company 1", searchResp.Providers[0].CompanyName)
}
