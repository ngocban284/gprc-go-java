gen:
	protoc -I=proto --go_out=generateProto proto/*.proto

clear:
	rm -rf generateProto/*.pb.go
run:
	go run main.go
