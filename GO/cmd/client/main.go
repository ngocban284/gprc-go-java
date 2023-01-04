package main

import (
	"bufio"
	"context"
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"
	pb "pcbook/generateProto"
	"pcbook/sample"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func createLaptop(laptopClient pb.LaptopServiceClient, laptop *pb.Laptop) {
	req := &pb.CreateLaptopRequest{
		Laptop: laptop,
	}
	// set timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := laptopClient.CreateLaptop(ctx, req)
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

func searchLaptop(laptopClient pb.LaptopServiceClient, filter *pb.Filter) {
	log.Print("searching for laptop with filter: ", filter)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.SearchLaptopRequest{
		Filter: filter,
	}
	stream, err := laptopClient.SearchLaptop(ctx, req)
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

func testCreateLaptop(laptopClient pb.LaptopServiceClient) {
	for i := 0; i < 10; i++ {
		createLaptop(laptopClient, sample.NewLaptop())
	}
}

func uploadImage(laptopClient pb.LaptopServiceClient, laptopID string, imageType string, imagePath string) {
	file, err := os.Open(imagePath)
	if err != nil {
		log.Fatal("cannot open image file: ", err)
	}
	defer file.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stream, err := laptopClient.UploadImage(ctx)
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

func testUploadImage(laptopClient pb.LaptopServiceClient) {
	laptop := sample.NewLaptop()
	createLaptop(laptopClient, laptop)
	uploadImage(laptopClient, laptop.GetId(), "jpeg", "tmp/kho-hieu.jpeg")
}

func testSearchLaptop(laptopClient pb.LaptopServiceClient) {
	for i := 0; i < 10; i++ {
		createLaptop(laptopClient, sample.NewLaptop())
	}

	filter := &pb.Filter{
		MaxPriceUsd: 1000,
		MinCpuCores: 4,
		MinCpuGhz:   2.5,
		MinRam:      &pb.Memory{Value: 8, Unit: pb.Memory_GIGABYTE},
	}
	searchLaptop(laptopClient, filter)
}

func main() {
	serverAddress := flag.String("serverAddress", "", "server address")
	flag.Parse()
	log.Print("dial server: ", *serverAddress)

	conn, err := grpc.Dial(*serverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatal("cannot dial server: ", err)
	}
	defer conn.Close()

	laptopClient := pb.NewLaptopServiceClient(conn)
	for i := 0; i < 10; i++ {
		createLaptop(laptopClient, sample.NewLaptop())
	}

	// filter := &pb.Filter{
	// 	MaxPriceUsd: 1000,
	// 	MinCpuCores: 4,
	// 	MinCpuGhz:   2.5,
	// 	MinRam:      &pb.Memory{Value: 8, Unit: pb.Memory_GIGABYTE},
	// }
	// searchLaptop(laptopClient, filter)
	testUploadImage(laptopClient)

}
