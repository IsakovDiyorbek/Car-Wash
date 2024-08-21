package service

import (
	"context"

	"github.com/Car-Wash/Car-Wash-Auth-Service/genproto/user"
	"github.com/Car-Wash/Car-Wash-Auth-Service/storage"
)

type UserService struct {
	stg storage.StorageI
	user.UnimplementedUserServiceServer
}

func NewUserService(stg storage.StorageI) *UserService {
	return &UserService{stg: stg}
}

func (u *UserService) GetProfile(ctx context.Context, req *user.GetProfileRequest) (*user.GetProfileResponse, error) {
	return u.stg.User().GetProfile(ctx, req)
}

func (u *UserService) UpdateProfile(ctx context.Context, req *user.UpdateProfileRequest) (*user.UpdateProfileResponse, error) {
	return u.stg.User().UpdateProfile(ctx, req)
}

func (u *UserService) ChangePassword(ctx context.Context, req *user.ChangePasswordRequest) (*user.ChangePasswordResponse, error) {
	return u.stg.User().ChangePassword(ctx, req)
}

func (u *UserService) GetAllUsers(ctx context.Context, req *user.GetAllUsersRequest) (*user.GetAllUsersResponse, error) {
	return u.stg.User().GetAllUsers(ctx, req)
}
