package client

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	pb "pcbook/generateProto"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LaptopClient struct {
	service pb.LaptopServiceClient
}

func NewLaptopClient(cc *grpc.ClientConn) *LaptopClient {
	return &LaptopClient{
		service: pb.NewLaptopServiceClient(cc),
	}
}

func (laptopClient *LaptopClient) CreateLaptop(laptop *pb.Laptop) {
	req := &pb.CreateLaptopRequest{
		Laptop: laptop,
	}
	// set timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := laptopClient.service.CreateLaptop(ctx, req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			if st.Code() == codes.AlreadyExists {
				log.Print("laptop already exists")

			} else {
				log.Fatal("cannot create laptop: ", err)
			}
			return
		}
	}
	log.Print("create laptop with id: ", res.Id)
}

func (laptopClient *LaptopClient) SearchLaptop(filter *pb.Filter) {
	log.Print("searching for laptop with filter: ", filter)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.SearchLaptopRequest{
		Filter: filter,
	}
	stream, err := laptopClient.service.SearchLaptop(ctx, req)
	if err != nil {
		log.Fatal("cannot search laptop: ", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			log.Print("stream is closed by server", err)
			return
		}
		if err != nil {
			log.Fatal("cannot receive response: ", err)
		}
		laptop := res.GetLaptop()
		log.Print("-found a laptop: ", laptop.GetId())
		log.Print("  +price: ", laptop.GetPriceUsd())
		log.Print("  +brand name: ", laptop.GetName())
		log.Print("  +cpu cores: ", laptop.GetCpu().GetNumberCores())
		log.Print("  +cpu: ", laptop.GetCpu().GetMinGhz())
		log.Print("  +ram: ", laptop.GetRam().GetValue(), laptop.GetRam().GetUnit())

	}

}

func (laptopClient *LaptopClient) UploadImage(laptopID string, imageType string, imagePath string) {
	file, err := os.Open(imagePath)
	if err != nil {
		log.Fatal("cannot open image file: ", err)
	}
	defer file.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stream, err := laptopClient.service.UploadImage(ctx)
	if err != nil {
		log.Fatal("cannot upload image: ", err)
	}

	req := &pb.UploadImageRequest{
		Data: &pb.UploadImageRequest_ImageInfo{
			ImageInfo: &pb.ImageInfo{
				LaptopId:  laptopID,
				ImageType: filepath.Ext(imagePath),
			},
		},
	}
	// log.Print("sending image info: ", req)
	err = stream.Send(req)
	if err != nil {
		log.Fatal("cannot send image info: ", err)
	}

	reader := bufio.NewReader(file)
	buffer := make([]byte, 1024)

	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("cannot read chunk to buffer: ", err)
		}

		req := &pb.UploadImageRequest{
			Data: &pb.UploadImageRequest_ChunkData{
				ChunkData: buffer[:n],
			},
		}
		// log.Print("sending chunk: ", req)
		err = stream.Send(req)
		if err != nil {
			log.Fatal("cannot send chunk to server: ", err, stream.RecvMsg(nil))
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal("cannot receive response: ", err)
	}

	log.Printf("image uploaded with id: %s, size: %d", res.GetId(), res.GetSize())

}

func (laptopClient *LaptopClient) RateLaptop(laptopIDs []string, scores []float64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stream, err := laptopClient.service.RateLaptop(ctx)
	if err != nil {
		return fmt.Errorf("cannot rate laptop: %w", err)
	}

	waitResponse := make(chan error)
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				log.Print("stream is closed by server")
				waitResponse <- nil
				return
			}
			if err != nil {
				waitResponse <- fmt.Errorf("cannot receive response: %w", err)
				return
			}
			log.Print("received response: ", res)
		}
	}()

	for i, laptopID := range laptopIDs {
		req := &pb.RateLaptopRequest{
			LaptopId: laptopID,
			Score:    scores[i],
		}
		log.Print("sending request: ", req)
		err := stream.Send(req)
		if err != nil {
			return fmt.Errorf("cannot send request: %v - %v", err, stream.RecvMsg(nil))
		}

		log.Print("sent request: ", req)
	}

	err = stream.CloseSend()
	if err != nil {
		return fmt.Errorf("cannot close stream: %w", err)
	}

	return <-waitResponse
}
