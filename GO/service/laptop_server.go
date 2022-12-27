package service

import (
	"context"
	"log"
	pb "pcbook/generateProto"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LaptopServer struct {
}

func NewLaptopServer() *LaptopServer {
	return &LaptopServer{}
}

func (s *LaptopServer) CreateLaptop(
	ctx context.Context,
	req *pb.CreateLaptopRequest,
	opts ...grpc.CallOption) (
	*pb.CreateLaptopResponse,
	error) {
	laptop := req.GetLaptop()
	log.Print("Received a request to create a laptop with ID: ", laptop.Id)

	if len(laptop.Id) > 0 {
		_, err := uuid.Parse(laptop.Id)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "laptop id is not a valid UUID : %v", err)
		}
	} else {
		id, err := uuid.NewRandom()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "cannot generate a new laptop ID: %v", err)
		}
		laptop.Id = id.String()

	}
	return &pb.CreateLaptopResponse{
		Id: laptop.Id,
	}, nil

}
