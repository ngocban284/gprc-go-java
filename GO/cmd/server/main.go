package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	pb "pcbook/generateProto"
	"pcbook/service"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func unaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	log.Print("---> intercepting unary method: ", info.FullMethod)
	return handler(ctx, req)
}

func streamInterceptor(
	srv interface{},
	stream grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	log.Print("---> intercepting stream method: ", info.FullMethod)
	return handler(srv, stream)
}

const (
	secretKey     = "secret"
	tokenDuration = 15 * time.Minute
)

// hello world
func main() {
	port := flag.Int("port", 0, "port to listen on")
	flag.Parse()
	log.Print("starting server on port: ", *port)

	authServer := service.NewAuthServer(
		service.NewInMemoryUserStore(),
		service.NewJWTManager(secretKey, tokenDuration),
	)

	laptopServer := service.NewLaptopServer(
		service.NewInMemoryLaptopStore(),
		service.NewDiskImageStore("img"),
		service.NewInMemoryRatingStore(),
	)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(unaryInterceptor),
		grpc.StreamInterceptor(streamInterceptor),
	)

	pb.RegisterAuthServiceServer(grpcServer, authServer)
	pb.RegisterLaptopServiceServer(grpcServer, laptopServer)
	reflection.Register(grpcServer)

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
