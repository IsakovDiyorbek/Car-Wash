package service


import (
	"context"

	"github.com/exam-5/Car-Wash-Booking-Service/genproto/carwash"
	"github.com/exam-5/Car-Wash-Booking-Service/storage"
)



type ReviewService struct{
	storage storage.StorageI
	carwash.UnimplementedReviewServiceServer
}

func NewReviewService(storage storage.StorageI) *ReviewService{
	return &ReviewService{storage: storage}
}

func (s *ReviewService) CreateReview(ctx context.Context,req *carwash.CreateReviewRequest) (*carwash.CreateReviewResponse, error) {
	return s.storage.Review().CreateReview(req)
}

func (s *ReviewService) GetReview(ctx context.Context, req *carwash.GetReviewRequest) (*carwash.GetReviewResponse, error) {
	return s.storage.Review().GetReview(req)
}

func (s *ReviewService) UpdateReview(ctx context.Context, req *carwash.UpdateReviewRequest) (*carwash.UpdateReviewResponse, error) {
	return s.storage.Review().UpdateReview(req)
}

func (s *ReviewService) DeleteReview(ctx context.Context, req *carwash.DeleteReviewRequest) (*carwash.DeleteReviewResponse, error) {
	return s.storage.Review().DeleteReview(req)
}

func (s *ReviewService) ListReviews(ctx context.Context, req *carwash.ListReviewsRequest) (*carwash.ListReviewsResponse, error) {
	return s.storage.Review().ListReviews(req)
}