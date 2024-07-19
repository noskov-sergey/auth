package users

import (
	"context"
	"github.com/noskov-sergey/auth/internal/model"
	desc "github.com/noskov-sergey/auth/pkg/user_v1"
)

type Usecase interface {
	Get(ctx context.Context, filter model.UserFilter) (*model.User, error)
	Create(ctx context.Context, user model.User) (int, error)
	Update(ctx context.Context, user model.User) error
	Delete(ctx context.Context, id int) error
}

type Implementation struct {
	desc.UnimplementedUserV1Server
	usecase Usecase
}

func New(u Usecase) *Implementation {
	return &Implementation{
		desc.UnimplementedUserV1Server{},
		u,
	}
}
