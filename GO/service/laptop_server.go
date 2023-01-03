package service

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log"
	pb "pcbook/generateProto"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const maxImageSize = 1 << 20 // 1MB

type LaptopServer struct {
	pb.UnimplementedLaptopServiceServer
	laptopStore LaptopStore
	imageStore  ImageStore
}

func NewLaptopServer(laptopStore LaptopStore, imageStore ImageStore) *LaptopServer {
	return &LaptopServer{
		laptopStore: laptopStore,
		imageStore:  imageStore,
	}
}

func (s *LaptopServer) CreateLaptop(
	ctx context.Context,
	req *pb.CreateLaptopRequest,
) (

	*pb.CreateLaptopResponse,
	error) {
	laptop := req.GetLaptop()
	log.Print("Received a request to create a laptop with ID: ", laptop.Id)

	if len(laptop.Id) > 0 {
		_, err := uuid.Parse(laptop.Id)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "laptop id is not a valid UUID : %v", err)
		}
	} else {
		id, err := uuid.NewRandom()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "cannot generate a new laptop ID: %v", err)
		}
		laptop.Id = id.String()

	}

	if ctx.Err() == context.Canceled {
		log.Print("request is canceled")
		return nil, status.Error(codes.Canceled, "request is canceled")
	}

	if ctx.Err() == context.DeadlineExceeded {
		log.Print("deadline is exceeded")
		return nil, status.Error(codes.DeadlineExceeded, "deadline is exceeded")
	}

	// save laptop to store
	err := s.laptopStore.Save(laptop)

	if err != nil {
		code := codes.Internal
		if errors.Is(err, ErrAlreadyExists) {
			code = codes.AlreadyExists
		}
		return nil, status.Errorf(code, "cannot save laptop to store: %v", err)
	}
	log.Print("saved laptop with id: ", laptop.Id)

	res := &pb.CreateLaptopResponse{
		Id: laptop.Id,
	}
	return res, nil
}

func (server *LaptopServer) SearchLaptop(
	req *pb.SearchLaptopRequest,
	stream pb.LaptopService_SearchLaptopServer,
) error {
	filter := req.GetFilter()
	log.Printf("receive a search-laptop request with filter: %v", filter)

	err := server.laptopStore.Search(
		stream.Context(),
		filter,
		func(laptop *pb.Laptop) error {
			res := &pb.SearchLaptopResponse{Laptop: laptop}
			err := stream.Send(res)
			if err != nil {
				return err
			}

			log.Printf("sent laptop with id: %s", laptop.GetId())
			return nil
		},
	)

	if err != nil {
		return status.Errorf(codes.Internal, "unexpected error: %v", err)
	}

	return nil
}

func (server *LaptopServer) UploadImage(stream pb.LaptopService_UploadImageServer) error {
	req, err := stream.Recv()
	if err != nil {
		log.Print("cannot receive image info: ", err)
		return status.Errorf(codes.Unknown, "cannot receive image info: %v", err)
	}

	laptopID := req.GetImageInfo().GetLaptopId()
	imageType := req.GetImageInfo().GetImageType()
	log.Printf("receive a request with  laptopID: %s, imageType: %s", laptopID, imageType)

	// check if laptop exists
	laptop, err := server.laptopStore.Find(laptopID)
	if err != nil {
		log.Print("cannot find laptop with id: ", laptopID)
		return status.Error(codes.NotFound, "cannot find laptop with id: "+laptopID)
	}
	if laptop == nil {
		log.Print("laptop with id not found: ", laptopID)
		return status.Error(codes.InvalidArgument, "laptop with id not found: "+laptopID)
	}

	// receive image data
	imageData := bytes.Buffer{}
	imageSize := 0

	for {
		log.Print("waiting to receive more data...")
		req, err := stream.Recv()
		if err == io.EOF {
			log.Print("finished uploading image data")
			break
		}
		if err != nil {
			log.Print("cannot receive image data: ", err)
			return status.Errorf(codes.Unknown, "cannot receive image data: %v", err)
		}

		chunk := req.GetChunkData()
		size := len(chunk)
		imageSize += size
		if imageSize > maxImageSize {
			log.Print("image size is too large")
			return status.Errorf(codes.InvalidArgument, "image size is too large: %d > %d", imageSize, maxImageSize)
		}

		_, err = imageData.Write(chunk)
		if err != nil {
			log.Print("cannot write chunk data to buffer: ", err)
			return status.Errorf(codes.Internal, "cannot write chunk data to buffer: %v", err)
		}
	}

	// save image to store
	imageID, err := server.imageStore.Save(laptopID, imageType, imageData)
	if err != nil {
		log.Print("cannot save image to store: ", err)
		return status.Errorf(codes.Internal, "cannot save image to store: %v", err)
	}

	// send response
	res := &pb.UploadImageResponse{
		Id:   imageID,
		Size: uint32(imageSize),
	}

	err = stream.SendAndClose(res)
	if err != nil {
		log.Print("cannot send response: ", err)
		return status.Errorf(codes.Internal, "cannot send response: %v", err)
	}

	log.Print("saved image with id: ", imageID)
	return nil
}
