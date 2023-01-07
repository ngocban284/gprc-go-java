package main

import (
	"flag"
	"fmt"
	"log"
	"pcbook/client"
	pb "pcbook/generateProto"
	"pcbook/sample"
	"strings"
	"time"

	"google.golang.org/grpc"
)

func testCreateLaptop(laptopClient *client.LaptopClient) {
	laptopClient.CreateLaptop(sample.NewLaptop())
}

func testSearchLaptop(laptopClient *client.LaptopClient) {
	for i := 0; i < 10; i++ {
		laptopClient.CreateLaptop(sample.NewLaptop())
	}

	filter := &pb.Filter{
		MaxPriceUsd: 3000,
		MinCpuCores: 4,
		MinCpuGhz:   2.5,
		MinRam:      &pb.Memory{Value: 8, Unit: pb.Memory_GIGABYTE},
	}

	laptopClient.SearchLaptop(filter)
}

func testUploadImage(laptopClient *client.LaptopClient) {
	laptop := sample.NewLaptop()
	laptopClient.CreateLaptop(laptop)
	laptopClient.UploadImage(laptop.GetId(), "jpg", "tmp/kho-hieu.jpg")
}

func testRateLaptop(laptopClient *client.LaptopClient) {
	n := 3
	laptopIDs := make([]string, n)

	for i := 0; i < n; i++ {
		laptop := sample.NewLaptop()
		laptopIDs[i] = laptop.GetId()
		laptopClient.CreateLaptop(laptop)
	}

	scores := make([]float64, n)
	for {
		fmt.Print("rate laptop (y/n)? ")
		var answer string
		fmt.Scan(&answer)

		if strings.ToLower(answer) != "y" {
			break
		}

		for i := 0; i < n; i++ {
			scores[i] = sample.RandomLaptopScore()
		}

		err := laptopClient.RateLaptop(laptopIDs, scores)
		if err != nil {
			log.Fatal(err)
		}
	}
}

const (
	username        = "admin"
	password        = "admin"
	refreshDuration = 30 * time.Second
)

func authMethods() map[string]bool {
	laptopServicePath := "/pcbook.LaptopService/"
	return map[string]bool{
		laptopServicePath + "CreateLaptop": true,
		laptopServicePath + "UploadImage":  true,
		laptopServicePath + "RateLaptop":   true,
	}
}

func main() {
	serverAddress := flag.String("serverAddress", "", "server address")
	flag.Parse()
	log.Print("dial server: ", *serverAddress)

	conn1, err := grpc.Dial(*serverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatal("cannot dial server: ", err)
	}
	defer conn1.Close()

	authClient := client.NewAuthClient(conn1, username, password)
	interceptor, err := client.NewAuthInterceptorClient(authClient, authMethods(), refreshDuration)
	if err != nil {
		log.Fatal(err)
	}

	conn2, err := grpc.Dial(
		*serverAddress,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(interceptor.Unary()),
		grpc.WithStreamInterceptor(interceptor.Stream()),
	)
	if err != nil {
		log.Fatal("cannot dial server: ", err)
	}
	defer conn2.Close()

	laptopClient := client.NewLaptopClient(conn2)
	testRateLaptop(laptopClient)
}
