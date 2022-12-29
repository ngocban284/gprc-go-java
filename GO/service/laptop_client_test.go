package service_test

import (
	"context"
	"net"
	"pcbook/sample"
	"pcbook/serializer"
	"pcbook/service"
	"testing"

	pb "pcbook/generateProto"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func TestLaptopClientCreateLaptop(t *testing.T) {
	t.Parallel()

	laptopStore := service.NewInMemoryLaptopStore()
	serverAddress := startTestLaptopServer(t, laptopStore)
	laptopClient := newTestLaptopClient(t, serverAddress)

	laptop := sample.NewLaptop()
	expectedId := laptop.Id

	req := &pb.CreateLaptopRequest{
		Laptop: laptop,
	}
	res, err := laptopClient.CreateLaptop(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, res.Id)
	require.Equal(t, expectedId, res.Id)

	other, err := laptopStore.Find(res.Id)
	require.NoError(t, err)
	require.NotNil(t, other)

	requireSameLaptop(t, laptop, other)
}

func startTestLaptopServer(t *testing.T, laptopStore service.LaptopStore) string {
	laptopServer := service.NewLaptopServer(laptopStore)

	grpcServer := grpc.NewServer()

	pb.RegisterLaptopServiceServer(grpcServer, laptopServer)

	listener, err := net.Listen("tcp", ":0")

	require.NoError(t, err, "cannot start test laptop server")

	go grpcServer.Serve(listener)

	return listener.Addr().String()

}

func newTestLaptopClient(t *testing.T, serverAddress string) pb.LaptopServiceClient {
	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())

	require.NoError(t, err, "cannot connect to test laptop server")

	return pb.NewLaptopServiceClient(conn)
}

func requireSameLaptop(t *testing.T, expected *pb.Laptop, actual *pb.Laptop) {
	json1, err := serializer.ProtobufToJSON(expected)
	require.NoError(t, err)

	json2, err := serializer.ProtobufToJSON(expected)
	require.NoError(t, err)

	require.Equal(t, json1, json2)
}
