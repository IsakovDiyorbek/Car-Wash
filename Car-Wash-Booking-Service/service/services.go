package service

import (
	"context"

	"github.com/exam-5/Car-Wash-Booking-Service/genproto/carwash"
	"github.com/exam-5/Car-Wash-Booking-Service/storage"
)



type ServicesService struct{
	storage storage.StorageI
	carwash.UnimplementedServicesServiceServer
}

func NewServicesService(storage storage.StorageI) *ServicesService{
	return &ServicesService{storage: storage}
}

func (s *ServicesService) CreateService(ctx context.Context,req *carwash.CreateServiceRequest) (*carwash.CreateServiceResponse, error) {
	return s.storage.Service().CreateService(req)
}

func (s *ServicesService) GetService(ctx context.Context, req *carwash.GetServiceRequest) (*carwash.GetServiceResponse, error) {
	return s.storage.Service().GetService(req)
}

func (s *ServicesService) UpdateService(ctx context.Context, req *carwash.UpdateServiceRequest) (*carwash.UpdateServiceResponse, error) {
	return s.storage.Service().UpdateService(req)
}

func (s *ServicesService) DeleteService(ctx context.Context, req *carwash.DeleteServiceRequest) (*carwash.DeleteServiceResponse, error) {
	return s.storage.Service().DeleteService(req)
}

func (s *ServicesService) ListServices(ctx context.Context, req *carwash.ListServicesRequest) (*carwash.ListServicesResponse, error) {
	return s.storage.Service().ListServices(req)
}

func (s *ServicesService) SearchServices(ctx context.Context, req *carwash.SearchServicesRequest) (*carwash.SearchServicesResponse, error) {
	return s.storage.Service().SearchServices(req)
}

func (s *ServicesService) GetPopularService(ctx context.Context, req *carwash.PopularServiceRequest) (*carwash.PopularServicesResponse, error) {
	return s.storage.Service().GetPopularService(req)
}