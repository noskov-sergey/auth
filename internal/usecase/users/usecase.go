package users

import (
	"context"
	"github.com/noskov-sergey/auth/internal/model"
)

type Repository interface {
	Get(ctx context.Context, filter model.UserFilter) (*model.User, error)
	Create(ctx context.Context, u model.User) (int, error)
	Update(ctx context.Context, u model.User) error
	Delete(ctx context.Context, id int) error
}

type useCase struct {
	rep Repository
}

func New(rep Repository) *useCase {
	return &useCase{rep: rep}
}
