package main

import (
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

const (
	secretKey     = "secret"
	tokenDuration = 15 * time.Minute
)

func accessibleRoles() map[string][]string {
	return map[string][]string{
		"/pb.AuthService/Login":          {"guest"},
		"/pb.LaptopService/CreateLaptop": {"admin"},
		"/pb.LaptopService/UploadImage":  {"admin"},
		"/pb.LaptopService/RateLaptop":   {"user"},
		"/pb.LaptopService/GetAll":       {"user"},
		"/pb.LaptopService/SearchLaptop": {"user"},
	}
}

// hello world
func main() {
	port := flag.Int("port", 0, "port to listen on")
	flag.Parse()
	log.Print("starting server on port: ", *port)

	jwtManager := service.NewJWTManager(secretKey, tokenDuration)

	authServer := service.NewAuthServer(
		service.NewInMemoryUserStore(),
		jwtManager,
	)

	laptopServer := service.NewLaptopServer(
		service.NewInMemoryLaptopStore(),
		service.NewDiskImageStore("img"),
		service.NewInMemoryRatingStore(),
	)

	interceptor := service.NewAuthInterceptor(jwtManager, accessibleRoles)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.Unary()),
		grpc.StreamInterceptor(interceptor.Stream()),
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
