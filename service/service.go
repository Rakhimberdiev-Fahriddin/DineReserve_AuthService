package service

import (
	pb "auth-service/generated/auth_service"
	"auth-service/storage/postgres"
	"context"
	"log/slog"
)

type AuthService struct {
	pb.UnimplementedAuthServiceServer
	User   *postgres.UserRepo
	Logger *slog.Logger
}

func NewAuthService(user *postgres.UserRepo) *AuthService {
	return &AuthService{User: user}
}

func (a *AuthService) GetUserProfile(ctx context.Context, in *pb.GetUserProfileRequest) (*pb.GetUserProfileResponse, error) {
	a.Logger.Info("gRPC method GetUserProfile")
	resp, err := a.User.GetUserProfile(in.Username)
	if err != nil {
		a.Logger.Error("Error getting user profile:", "error", err.Error())
		return nil, err
	}

	return resp, nil
}

func (a *AuthService) UpdateUserProfile(ctx context.Context, in *pb.UpdateUserProfileRequest) (*pb.UpdateUserProfileResponse, error) {
	a.Logger.Info("gRPC method UpdateUserProfile")
	resp, err := a.User.UpdateUserProfile(in)

	if err != nil {
		a.Logger.Error("Error updating user profile:", "error", err.Error())
		return &pb.UpdateUserProfileResponse{
			Message: resp.Message,
		}, err
	}

	return &pb.UpdateUserProfileResponse{
		Message: resp.Message,
	}, nil
}

func (a *AuthService) LogoutUser(ctx context.Context, in *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	a.Logger.Info("gRPC method LogoutUser")
	resp, err := a.User.LogoutUser(in.UserId)

	if err != nil {
		a.Logger.Error("Error logging out user:", "error", err.Error())
		return &pb.LogoutResponse{
			Message: resp.Message,
		}, err
	}

	return &pb.LogoutResponse{
		Message: resp.Message,
	}, nil
}
