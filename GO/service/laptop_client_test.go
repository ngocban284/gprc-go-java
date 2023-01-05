package service_test

import (
	"context"
	"io"
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
	serverAddress := startTestLaptopServer(t, laptopStore, nil, nil)
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

func TestLaptopClientSearchLaptop(t *testing.T) {
	t.Parallel()

	filter := &pb.Filter{
		MaxPriceUsd: 2000,
		MinCpuCores: 4,
		MinCpuGhz:   2.2,
		MinRam:      &pb.Memory{Value: 8, Unit: pb.Memory_GIGABYTE},
	}

	laptopStore := service.NewInMemoryLaptopStore()
	expectedID := make(map[string]bool)

	for i := 0; i < 6; i++ {
		laptop := sample.NewLaptop()

		switch i {
		case 0:
			laptop.PriceUsd = 2500
		case 1:
			laptop.Cpu.NumberCores = 2
		case 2:
			laptop.Cpu.MinGhz = 2.0
		case 3:
			laptop.Ram = &pb.Memory{Value: 4096, Unit: pb.Memory_MEGABYTE}
		case 4:
			laptop.PriceUsd = 1999
			laptop.Cpu.NumberCores = 4
			laptop.Cpu.MinGhz = 2.5
			laptop.Cpu.MaxGhz = 4.5
			laptop.Ram = &pb.Memory{Value: 16, Unit: pb.Memory_GIGABYTE}
			expectedID[laptop.Id] = true
		case 5:
			laptop.PriceUsd = 2000
			laptop.Cpu.NumberCores = 6
			laptop.Cpu.MinGhz = 2.8
			laptop.Cpu.MaxGhz = 5.0
			laptop.Ram = &pb.Memory{Value: 64, Unit: pb.Memory_GIGABYTE}
			expectedID[laptop.Id] = true
		default:
			panic("should not reach here")
		}

		err := laptopStore.Save(laptop)
		require.NoError(t, err)
	}

	serverAddress := startTestLaptopServer(t, laptopStore, nil, nil)
	laptopClient := newTestLaptopClient(t, serverAddress)

	req := &pb.SearchLaptopRequest{
		Filter: filter,
	}
	stream, err := laptopClient.SearchLaptop(context.Background(), req)
	require.NoError(t, err)
	found := 0
	for {
		res, err := stream.Recv()
		if err != nil {
			break
		}
		require.NoError(t, err)
		require.Contains(t, expectedID, res.Laptop.Id)
		found++
	}
	require.Equal(t, len(expectedID), found)
}

func TestClientRateLaptop(t *testing.T) {
	t.Parallel()

	laptopStore := service.NewInMemoryLaptopStore()
	ratingStore := service.NewInMemoryRatingStore()

	laptop := sample.NewLaptop()
	err := laptopStore.Save(laptop)
	require.NoError(t, err)

	serverAddress := startTestLaptopServer(t, laptopStore, nil, ratingStore)
	laptopClient := newTestLaptopClient(t, serverAddress)

	stream, err := laptopClient.RateLaptop(context.Background())
	require.NoError(t, err)

	scores := []float64{8, 7.5, 10}
	avarage := []float64{8, 7.75, 8.5}
	n := len(scores)

	for i := 0; i < n; i++ {
		req := &pb.RateLaptopRequest{
			LaptopId: laptop.Id,
			Score:    scores[i],
		}
		err := stream.Send(req)
		require.NoError(t, err)
	}

	err = stream.CloseSend()
	require.NoError(t, err)

	for resCount := 0; ; resCount++ {
		res, err := stream.Recv()
		if err == io.EOF {
			require.Equal(t, n, resCount)
			return
		}
		require.NoError(t, err)
		require.Equal(t, laptop.Id, res.LaptopId)
		require.Equal(t, uint32(resCount+1), res.RatedCount)
		require.Equal(t, avarage[resCount], res.AverageCore)
	}

}

func startTestLaptopServer(
	t *testing.T,
	laptopStore service.LaptopStore,
	imageStore service.ImageStore, ratingStore service.RatingStore) string {

	laptopServer := service.NewLaptopServer(laptopStore, imageStore, ratingStore)

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
