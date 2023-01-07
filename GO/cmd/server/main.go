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
	laptopServicePath := "/pcbook.LaptopService/"
	return map[string][]string{
		laptopServicePath + "CreateLaptop": {"admin"},
		laptopServicePath + "UploadImage":  {"admin"},
		laptopServicePath + "RateLaptop":   {"admin", "user"},
	}
}

func createUser(userStore service.UserStore, username, password, role string) error {
	user, err := service.NewUser(username, password, role)
	if err != nil {
		return err
	}
	return userStore.Save(user)
}

func seedUsers(userStore service.UserStore) error {
	err := createUser(userStore, "admin", "admin", "admin")
	if err != nil {
		return err
	}
	return createUser(userStore, "user1", "user1", "user")
}

// hello world
func main() {
	port := flag.Int("port", 0, "port to listen on")
	flag.Parse()
	log.Print("starting server on port: ", *port)

	userStore := service.NewInMemoryUserStore()
	err := seedUsers(userStore)
	if err != nil {
		log.Fatal("cannot seed users: ", err)
	}

	jwtManager := service.NewJWTManager(secretKey, tokenDuration)

	authServer := service.NewAuthServer(
		userStore,
		jwtManager,
	)

	laptopServer := service.NewLaptopServer(
		service.NewInMemoryLaptopStore(),
		service.NewDiskImageStore("img"),
		service.NewInMemoryRatingStore(),
	)

	interceptor := service.NewAuthInterceptor(jwtManager, accessibleRoles())

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
