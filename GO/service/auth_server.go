package service

import (
	"context"
	pb "pcbook/generateProto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServer struct {
	pb.UnimplementedAuthServiceServer
	userStore  UserStore
	jwtManager *JWTManager
}

func NewAuthServer(userStore UserStore, jwtManager *JWTManager) *AuthServer {
	return &AuthServer{
		userStore:  userStore,
		jwtManager: jwtManager,
	}
}

func (s *AuthServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, err := s.userStore.Find(req.GetUsername())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot find user: %v", err)
	}
	if user == nil || user.IsCorrectPassword(req.GetPassword()) {
		return nil, status.Errorf(codes.NotFound, "invalid username or password")
	}

	token, err := s.jwtManager.GenerateToken(user.UserName, user.Role)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot generate token: %v", err)
	}

	res := &pb.LoginResponse{
		AccessToken: token,
	}

	return res, nil
}
