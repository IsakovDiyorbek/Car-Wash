package service

import (
	"context"

	"github.com/exam-5/Car-Wash-Booking-Service/genproto/carwash"
	"github.com/exam-5/Car-Wash-Booking-Service/storage"
)

type BookingService struct {
	storage storage.StorageI
	carwash.UnimplementedBookingServiceServer
}

func NewBookingService(storage storage.StorageI) *BookingService {
	return &BookingService{storage: storage}
}

func (s *BookingService) CreateBooking(ctx context.Context, req *carwash.CreateBookingRequest) (*carwash.CreateBookingResponse, error) {
	return s.storage.Booking().CreateBooking(req)
}

func (s *BookingService) GetBooking(ctx context.Context, req *carwash.GetBookingRequest) (*carwash.GetBookingResponse, error) {
	return s.storage.Booking().GetBooking(req)
}

func (s *BookingService) UpdateBooking(ctx context.Context, req *carwash.UpdateBookingRequest) (*carwash.UpdateBookingResponse, error) {
	return s.storage.Booking().UpdateBooking(req)
}

func (s *BookingService) DeleteBooking(ctx context.Context, req *carwash.DeleteBookingRequest) (*carwash.DeleteBookingResponse, error) {
	return s.storage.Booking().DeleteBooking(req)
}

func (s *BookingService) ListBookings(ctx context.Context, req *carwash.ListBookingsRequest) (*carwash.ListBookingsResponse, error) {
	return s.storage.Booking().ListBookings(req)
}
