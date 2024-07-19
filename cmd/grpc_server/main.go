package main

import (
	"context"
	"flag"
	"github.com/jackc/pgx/v5/pgxpool"
	userApi "github.com/noskov-sergey/auth/internal/api/users"
	"github.com/noskov-sergey/auth/internal/config"
	userRepository "github.com/noskov-sergey/auth/internal/repository/users"
	userUsecase "github.com/noskov-sergey/auth/internal/usecase/users"
	desc "github.com/noskov-sergey/auth/pkg/user_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

func main() {
	ctx := context.Background()

	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	cfg, err := config.NewGPRCConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	pgConfig, err := config.NewPGConfig()
	if err != nil {
		log.Fatalf("failed to get pg config: %v", err)
	}

	lis, err := net.Listen("tcp", cfg.Address())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	pool, err := pgxpool.New(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	fileRep := userRepository.NewUserRepository(pool)
	usecase := userUsecase.New(fileRep)

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterUserV1Server(s, userApi.New(usecase))

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
