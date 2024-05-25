package main

import (
	"context"
	"fmt"
	desc "github.com/noskov-sergey/auth/pkg/user_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"net"
	"time"
)

const grpcPort = 50051

type server struct {
	desc.UnimplementedUserV1Server
}

func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	log.Printf("User id: %v", req.GetId())

	return &desc.GetResponse{
		User: &desc.User{
			Id:        1,
			Name:      "testName",
			Email:     "testName@mail.ru",
			Role:      desc.Enum.Enum(1),
			CreatedAt: timestamppb.New(time.Now()),
			UpdatedAt: timestamppb.New(time.Now()),
		},
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterUserV1Server(s, &server{})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
