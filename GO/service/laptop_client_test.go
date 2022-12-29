package service_test

import (
	"context"
	"net"
	"pcbook/sample"
	"pcbook/service"
	"testing"

	pb "pcbook/generateProto"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func TestLaptopClientCreateLaptop(t *testing.T) {
	t.Parallel()

	laptopServer, serverAddress := startTestLaptopServer(t)
	laptopClient := newTestLaptopClient(t, serverAddress)

	laptop := sample.NewLaptop()
	expectedId := laptop.Id

	req := &pb.CreateLaptopRequest{
		Laptop: laptop,
	}
	res, err := laptopClient.CreateLaptop(context.Background(), req)
}

func startTestLaptopServer(t *testing.T) (*service.LaptopServer, string) {
	laptopServer := service.NewLaptopServer(service.NewInMemoryLaptopStore())

	grpcServer := grpc.NewServer()

	pb.RegisterLaptopServiceServer(grpcServer, laptopServer)

	listener, err := net.Listen("tcp", ":0")

	require.NoError(t, err, "cannot start test laptop server")

	go grpcServer.Serve(listener)

	return laptopServer, listener.Addr().String()

}

func newTestLaptopClient(t *testing.T, serverAddress string) pb.LaptopServiceClient {
	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())

	require.NoError(t, err, "cannot connect to test laptop server")

	return pb.NewLaptopServiceClient(conn)
}
