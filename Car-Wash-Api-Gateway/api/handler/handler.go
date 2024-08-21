package handler

import (
	"github.com/Car-Wash/Car-Wash-Api-Gateway/clients"
)

type Handler struct {
	Client *clients.Client
}

func NewHandler() *Handler {
	conn := clients.NewClient()
	return &Handler{
		Client: conn,
	}
}
