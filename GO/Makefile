gen:
	protoc -I=proto --go-grpc_out=generateProto proto/*.proto

clear:
	rm -rf generateProto/*.pb.go
run:
	go run main.go
