package service

import (
	"context"

	"github.com/Car-Wash/Car-Wash-Booking-Service/genproto/carwash"
	"github.com/Car-Wash/Car-Wash-Booking-Service/storage"
)

type ProviderService struct {
	storage storage.StorageI
	carwash.UnimplementedProviderServiceServer
}

func NewProviderService(storage storage.StorageI) *ProviderService {
	return &ProviderService{storage: storage}
}

func (s *ProviderService) CreateProvider(ctx context.Context, req *carwash.CreateProviderRequest) (*carwash.CreateProviderResponse, error) {
	return s.storage.Provider().CreateProvider(req)
}

func (s *ProviderService) GetProvider(ctx context.Context, req *carwash.GetProviderRequest) (*carwash.GetProviderResponse, error) {
	return s.storage.Provider().GetProvider(req)
}

func (s *ProviderService) UpdateProvider(ctx context.Context, req *carwash.UpdateProviderRequest) (*carwash.UpdateProviderResponse, error) {
	return s.storage.Provider().UpdateProvider(req)
}

func (s *ProviderService) DeleteProvider(ctx context.Context, req *carwash.DeleteProviderRequest) (*carwash.DeleteProviderResponse, error) {
	return s.storage.Provider().DeleteProvider(req)
}

func (s *ProviderService) ListProviders(ctx context.Context, req *carwash.ListProvidersRequest) (*carwash.ListProvidersResponse, error) {
	return s.storage.Provider().ListProviders(req)
}

func (s *ProviderService) SearchProviders(ctx context.Context, req *carwash.SearchProvidersRequest) (*carwash.SearchProvidersResponse, error) {
	return s.storage.Provider().SearchProviders(req)
}
