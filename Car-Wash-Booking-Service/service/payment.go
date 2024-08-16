package service




import (
	"context"

	"github.com/exam-5/Car-Wash-Booking-Service/genproto/carwash"
	"github.com/exam-5/Car-Wash-Booking-Service/storage"
)



type PaymentService struct{
	storage storage.StorageI
	carwash.UnimplementedPaymentServiceServer
}

func NewPaymentService(storage storage.StorageI) *PaymentService{
	return &PaymentService{storage: storage}
}

func (s *PaymentService) CreatePayment(ctx context.Context,req *carwash.CreatePaymentRequest) (*carwash.CreatePaymentResponse, error) {
	return s.storage.Payment().CreatePayment(req)
}

func (s *PaymentService) GetPayment(ctx context.Context, req *carwash.GetPaymentRequest) (*carwash.GetPaymentResponse, error) {
	return s.storage.Payment().GetPayment(req)
}


func (s *PaymentService) ListPayments(ctx context.Context, req *carwash.ListPaymentsRequest) (*carwash.ListPaymentsResponse, error) {
	return s.storage.Payment().ListPayments(req)
}

func (s *PaymentService) UpdatePayment(ctx context.Context, req *carwash.UpdatePaymentRequest) (*carwash.UpdatePaymentResponse, error) {
	return s.storage.Payment().UpdatePayment(req)
}