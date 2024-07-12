package main

import (
	"context"
	"flag"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/noskov-sergey/auth/internal/config"
	desc "github.com/noskov-sergey/auth/pkg/user_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"net"
	"time"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

type server struct {
	desc.UnimplementedUserV1Server
	pool *pgxpool.Pool
}

func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	log.Printf("User id: %v", req.GetId())

	builder := sq.Select(
		"id",
		"name",
		"email",
		"role",
		"created_at",
		"updated_at",
	).
		PlaceholderFormat(sq.Dollar).
		From("users").
		Where(sq.Eq{"id": req.GetId()})

	sqlQuery, args, err := builder.ToSql()
	if err != nil {
		log.Fatalf("to sql: %v", err)
	}

	type user struct {
		id        int
		name      string
		email     string
		role      string
		createdAt time.Time
		updatedAt time.Time
	}

	var User user

	if err = s.pool.QueryRow(ctx, sqlQuery, args...).Scan(&User.id, &User.name, &User.email, &User.role, &User.createdAt, &User.updatedAt); err != nil {
		log.Fatalf("query row scan: %v", err)
		return nil, fmt.Errorf("query row scan: %w", err)
	}

	return &desc.GetResponse{
		User: &desc.User{
			Id:        int64(User.id),
			Name:      User.name,
			Email:     User.email,
			CreatedAt: timestamppb.New(User.createdAt),
			UpdatedAt: timestamppb.New(User.updatedAt),
		},
	}, nil
}

func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	log.Printf("Get CreateRequest: %v", req.GetName())

	builder := sq.Insert("users").
		PlaceholderFormat(sq.Dollar).
		Columns("name", "email", "role").
		Values(req.GetName(), req.GetEmail(), int(req.GetRole())).
		Suffix("RETURNING id")

	sqlQuery, args, err := builder.ToSql()
	if err != nil {
		log.Fatalf("to sql: %v", err)
	}

	var insertedID int

	if err = s.pool.QueryRow(ctx, sqlQuery, args...).Scan(&insertedID); err != nil {
		log.Printf("query row: %v", err)
		return nil, fmt.Errorf("query row: %w", err)
	}

	return &desc.CreateResponse{
		Id: int64(insertedID),
	}, nil
}

func (s *server) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	log.Printf("UpdateMethod - User id: %v, rename to: %s, new email: %s", req.GetId(), req.GetName(), req.GetMail())
	builder := sq.
		Update("users").
		PlaceholderFormat(sq.Dollar).
		SetMap(map[string]any{
			"name":       req.Name,
			"email":      req.Mail,
			"updated_at": time.Now(),
		}).
		Where("id = ?", req.Id)

	sqlQuery, args, err := builder.ToSql()

	if _, err = s.pool.Exec(ctx, sqlQuery, args...); err != nil {
		log.Printf("exec row: %v", err)
		return nil, fmt.Errorf("exec row: %w", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *server) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	log.Printf("DeleteMethod - User id: %v", req.GetId())
	builder := sq.
		Delete("users").
		PlaceholderFormat(sq.Dollar).
		Where("id = ?", int(req.Id))

	sqlQuery, args, err := builder.ToSql()
	if err != nil {
		log.Fatalf("to sql: %v", err)
	}

	if _, err = s.pool.Exec(ctx, sqlQuery, args...); err != nil {
		log.Printf("exec row: %v", err)
		return nil, fmt.Errorf("exec row: %w", err)
	}

	return &emptypb.Empty{}, nil
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

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterUserV1Server(s, &server{pool: pool})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
