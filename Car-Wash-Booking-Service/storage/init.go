package storage

import (
	pb "github.com/exam-5/Car-Wash-Booking-Service/genproto/carwash"
)


type StorageI interface{
	Provider() ProviderI
	Review() Review
	Service() Service
	Booking() BookingI
	Payment() PaymentI
	Notification() NotificationI
}
type BookingI interface {
	CreateBooking(req *pb.CreateBookingRequest) (*pb.CreateBookingResponse, error)
	GetBooking(req *pb.GetBookingRequest) (*pb.GetBookingResponse, error)
	UpdateBooking(req *pb.UpdateBookingRequest) (*pb.UpdateBookingResponse, error)
	DeleteBooking(req *pb.DeleteBookingRequest) (*pb.DeleteBookingResponse, error)
	ListBookings(req *pb.ListBookingsRequest) (*pb.ListBookingsResponse, error)
}

type PaymentI interface {
	CreatePayment(*pb.CreatePaymentRequest) (*pb.CreatePaymentResponse, error)
	GetPayment(*pb.GetPaymentRequest) (*pb.GetPaymentResponse, error)
	ListPayments(*pb.ListPaymentsRequest) (*pb.ListPaymentsResponse, error)
	UpdatePayment(*pb.UpdatePaymentRequest) (*pb.UpdatePaymentResponse, error)

}

type ProviderI interface {
	CreateProvider(*pb.CreateProviderRequest) (*pb.CreateProviderResponse, error)
	GetProvider(*pb.GetProviderRequest) (*pb.GetProviderResponse, error)
	UpdateProvider(*pb.UpdateProviderRequest) (*pb.UpdateProviderResponse, error)
	DeleteProvider(*pb.DeleteProviderRequest) (*pb.DeleteProviderResponse, error)
	ListProviders(*pb.ListProvidersRequest) (*pb.ListProvidersResponse, error)
	SearchProviders(*pb.SearchProvidersRequest) (*pb.SearchProvidersResponse, error)

}

type Review interface {
	CreateReview(*pb.CreateReviewRequest) (*pb.CreateReviewResponse, error)
	GetReview(*pb.GetReviewRequest) (*pb.GetReviewResponse, error)
	UpdateReview(*pb.UpdateReviewRequest) (*pb.UpdateReviewResponse, error)
	DeleteReview(*pb.DeleteReviewRequest) (*pb.DeleteReviewResponse, error)
	ListReviews(*pb.ListReviewsRequest) (*pb.ListReviewsResponse, error)
}

type Service interface {
	CreateService(*pb.CreateServiceRequest) (*pb.CreateServiceResponse, error)
	GetService(*pb.GetServiceRequest) (*pb.GetServiceResponse, error)
	UpdateService(*pb.UpdateServiceRequest) (*pb.UpdateServiceResponse, error)
	DeleteService(*pb.DeleteServiceRequest) (*pb.DeleteServiceResponse, error)
	ListServices(*pb.ListServicesRequest) (*pb.ListServicesResponse, error)
	SearchServices(*pb.SearchServicesRequest) (*pb.SearchServicesResponse, error)

}


type NotificationI interface {
	AddNotification(req *pb.AddNotificationRequest) (*pb.AddNotificationResponse, error)
	GetNotifications(req *pb.GetNotificationsRequest) (*pb.GetNotificationsResponse, error)
	MarkNotificationAsRead(req *pb.MarkNotificationAsReadRequest) (*pb.MarkNotificationAsReadResponse, error)
}