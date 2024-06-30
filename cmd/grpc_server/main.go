package main

import (
	"context"
	"flag"
	"github.com/noskov-sergey/auth/internal/config"
	desc "github.com/noskov-sergey/auth/pkg/user_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"math/rand"
	"net"
	"strconv"
	"time"
)

const (
	grpcPort = 50051
	dbDSN    = "host=localhost port=54321 dbname=auth user=auth-user password=auth-password"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

type server struct {
	desc.UnimplementedUserV1Server
}

func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	log.Printf("User id: %v", req.GetId())

	return &desc.GetResponse{
		User: &desc.User{
			Id:        req.GetId(),
			Name:      "testName" + strconv.Itoa(int(req.GetId())),
			Email:     "testName@mail.ru",
			Role:      desc.Enum.Enum(1),
			CreatedAt: timestamppb.New(time.Now()),
			UpdatedAt: timestamppb.New(time.Now()),
		},
	}, nil
}

func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	i := rand.Intn(100)

	log.Printf("CreateRequest: %v", i)

	return &desc.CreateResponse{
		Id: int64(i),
	}, nil
}

func (s *server) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	log.Printf("UpdateMethod - User id: %v, rename to: %s, new email: %s", req.GetId(), req.GetName(), req.GetMail())

	return &emptypb.Empty{}, nil
}

func (s *server) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	log.Printf("DeleteMethod - User id: %v", req.GetId())

	return &emptypb.Empty{}, nil
}

func main() {
	flag.Parse()

	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	grpcConfig, err := config.NewGPRCConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	lis, err := net.Listen("tcp", grpcConfig.Address())
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
