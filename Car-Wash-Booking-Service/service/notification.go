package service

import (
	"context"

	"github.com/Car-Wash/Car-Wash-Booking-Service/genproto/carwash"
	"github.com/Car-Wash/Car-Wash-Booking-Service/storage"
)

type NotificationService struct {
	storage storage.StorageI
	carwash.UnimplementedNotificationServiceServer
}

func NewNotificationService(storage storage.StorageI) *NotificationService {
	return &NotificationService{storage: storage}
}

func (s *NotificationService) AddNotification(ctx context.Context, req *carwash.AddNotificationRequest) (*carwash.AddNotificationResponse, error) {
	return s.storage.Notification().AddNotification(req)
}

func (s *NotificationService) GetNotifications(ctx context.Context, req *carwash.GetNotificationsRequest) (*carwash.GetNotificationsResponse, error) {
	return s.storage.Notification().GetNotifications(req)
}

func (s *NotificationService) MarkNotificationAsRead(ctx context.Context, req *carwash.MarkNotificationAsReadRequest) (*carwash.MarkNotificationAsReadResponse, error) {
	return s.storage.Notification().MarkNotificationAsRead(req)
}
