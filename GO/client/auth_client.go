package client

import (
	"context"
	pb "pcbook/generateProto"
	"time"

	"google.golang.org/grpc"
)

type AuthClient struct {
	service  pb.AuthServiceClient
	userName string
	password string
}

func NewAuthClient(cc *grpc.ClientConn, userName, password string) *AuthClient {
	return &AuthClient{
		service:  pb.NewAuthServiceClient(cc),
		userName: userName,
		password: password,
	}
}

func (c *AuthClient) Login() (string, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.LoginRequest{
		Username: c.userName,
		Password: c.password,
	}

	res, err := c.service.Login(ctx, req)
	if err != nil {
		return "", err
	}

	return res.AccessToken, nil
}
