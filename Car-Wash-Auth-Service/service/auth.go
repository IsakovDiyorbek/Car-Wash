package service

import (
	"context"

	"github.com/Car-Wash/Car-Wash-Auth-Service/genproto/auth"
	"github.com/Car-Wash/Car-Wash-Auth-Service/storage"
)

type AuthService struct {
	stg storage.StorageI
	auth.UnimplementedAuthServiceServer
}

func NewAuthService(stg storage.StorageI) *AuthService {
	return &AuthService{stg: stg}
}

func (a *AuthService) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	return a.stg.Auth().Register(ctx, req)
}

func (a *AuthService) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	return a.stg.Auth().Login(ctx, req)
}

func (a *AuthService) Logout(ctx context.Context, req *auth.LogoutRequest) (*auth.LogoutResponse, error) {
	return a.stg.Auth().Logout(ctx, req)
}

func (a *AuthService) ResetPassword(ctx context.Context, req *auth.ResetPasswordRequest) (*auth.ResetPasswordResponse, error) {
	return a.stg.Auth().ResetPassword(ctx, req)
}
