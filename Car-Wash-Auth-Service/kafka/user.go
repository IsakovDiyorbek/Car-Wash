package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Car-Wash/Car-Wash-Auth-Service/genproto/user"
	"github.com/Car-Wash/Car-Wash-Auth-Service/service"
)

func Change(User *service.UserService) func(message []byte) {
	return func(message []byte) {
		var req user.ChangePasswordRequest
		if err := json.Unmarshal(message, &req); err != nil {
			log.Printf("Cannot unmarshal JSON: %v", err)
			return
		}
		resp, err := User.ChangePassword(context.Background(), &req)
		if err != nil {
			log.Printf("Cannot user register via Kafka: %v", err)
			return
		}
		log.Printf("Register user via Kafka: %+v", resp)
	}
}
