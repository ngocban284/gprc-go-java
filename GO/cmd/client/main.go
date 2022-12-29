package main

import (
	"context"
	"flag"
	"log"
	pb "pcbook/generateProto"
	"pcbook/sample"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// hello world
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
	laptop := sample.NewLaptop()

	req := &pb.CreateLaptopRequest{
		Laptop: laptop,
	}
	res, err := laptopClient.CreateLaptop(context.Background(), req)
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
