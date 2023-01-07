gen_grpc:
	protoc -I=proto --go-grpc_out=generateProto proto/*.proto
gen:
	protoc -I=proto --go_out=generateProto proto/*.proto
clear:
	rm -rf generateProto/*.pb.go
server:
	go run cmd/server/main.go -port 8080
client:
	go run cmd/client/main.go -serverAddress 0.0.0.0:8080
cert:
	cd ./cert; ./gen.sh; cd ..
.PHONY: gen_grpc gen clear server client cert