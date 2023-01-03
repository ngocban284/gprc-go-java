package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	pb "pcbook/generateProto"
	"pcbook/service"

	"google.golang.org/grpc"
)

// hello world
func main() {
	port := flag.Int("port", 0, "port to listen on")
	flag.Parse()
	log.Print("starting server on port: ", *port)

	laptopServer := service.NewLaptopServer(
		service.NewInMemoryLaptopStore(),
		service.NewDiskImageStore("img"),
	)

	grpcServer := grpc.NewServer()
	pb.RegisterLaptopServiceServer(grpcServer, laptopServer)

	address := fmt.Sprintf("0.0.0.0:%d", *port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
