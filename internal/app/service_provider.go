package app

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/noskov-sergey/auth/internal/config"
	"github.com/noskov-sergey/platform-common/pkg/closer"
	"log"

	userApi "github.com/noskov-sergey/auth/internal/api/users"
	userRepository "github.com/noskov-sergey/auth/internal/repository/users"
	userUsecase "github.com/noskov-sergey/auth/internal/usecase/users"
)

type serviceProvider struct {
	pgConfig   config.PGConfig
	grpcConfig config.GRPCConfig

	pgPool *pgxpool.Pool
	uRep   *userRepository.UserRepository

	userUsecase *userUsecase.UseCase

	userImpl *userApi.Implementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := config.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %v", err)
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := config.NewGPRCConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %v", err)
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

func (s *serviceProvider) PgPool(ctx context.Context) *pgxpool.Pool {
	if s.pgPool == nil {
		pool, err := pgxpool.New(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to connect to database: %v", err)
		}

		err = pool.Ping(ctx)
		if err != nil {
			log.Fatalf("ping error: %v", err)
		}
		closer.Add(func() error {
			pool.Close()
			return nil
		})

		s.pgPool = pool
	}

	return s.pgPool
}

func (s *serviceProvider) UserRepository(ctx context.Context) *userRepository.UserRepository {
	if s.uRep == nil {
		s.uRep = userRepository.NewUserRepository(s.PgPool(ctx))
	}

	return s.uRep
}

func (s *serviceProvider) UserUsecase(ctx context.Context) userApi.Usecase {
	if s.userUsecase == nil {
		s.userUsecase = userUsecase.New(
			s.UserRepository(ctx),
		)
	}

	return s.userUsecase
}

func (s *serviceProvider) UImpl(ctx context.Context) *userApi.Implementation {
	if s.userImpl == nil {
		s.userImpl = userApi.New(
			s.UserUsecase(ctx),
		)
	}

	return s.userImpl
}
