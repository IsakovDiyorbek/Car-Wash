package clients

import (
	"fmt"
	"log"

	"github.com/Car-Wash/Car-Wash-Api-Gateway/genproto/auth"
	"github.com/Car-Wash/Car-Wash-Api-Gateway/genproto/carwash"
	kafka "github.com/Car-Wash/Car-Wash-Api-Gateway/kafka/producer"
	"github.com/casbin/casbin/v2"

	u "github.com/Car-Wash/Car-Wash-Api-Gateway/genproto/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	Service      carwash.ServicesServiceClient
	Provider     carwash.ProviderServiceClient
	Booking      carwash.BookingServiceClient
	Review       carwash.ReviewServiceClient
	Payment      carwash.PaymentServiceClient
	Notification carwash.NotificationServiceClient
	User         u.UserServiceClient
	Auth         auth.AuthServiceClient
	Kafka        kafka.KafkaProducer
	Enforcer     *casbin.Enforcer
}

func NewClient() *Client {
	booking, err := grpc.NewClient(fmt.Sprintf("booking_service%s", ":7777"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Error while NEwclient: ", err.Error())
	}
	user, err := grpc.NewClient(fmt.Sprintf("auth_service%s", ":9999"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Error while NEwclient: ", err.Error())
	}

	enforcer, err := casbin.NewEnforcer("config/model.conf", "config/policy.csv")
	if err != nil {
		log.Fatal("Error while NEwclient: ", err.Error())
	}
	kafkaProducer, err := kafka.NewKafkaProducer([]string{"kafka:9092"})
	if err != nil {
		return nil
	}

	return &Client{
		Service:      carwash.NewServicesServiceClient(booking),
		Provider:     carwash.NewProviderServiceClient(booking),
		Booking:      carwash.NewBookingServiceClient(booking),
		Review:       carwash.NewReviewServiceClient(booking),
		Payment:      carwash.NewPaymentServiceClient(booking),
		Notification: carwash.NewNotificationServiceClient(booking),
		Auth:         auth.NewAuthServiceClient(user),
		User:         u.NewUserServiceClient(user),
		Enforcer:     enforcer,
		Kafka:        kafkaProducer,
	}
}
