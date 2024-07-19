package users

import (
	"context"
	"fmt"
	"github.com/noskov-sergey/auth/internal/model"
)

func (u *useCase) Create(ctx context.Context, user model.User) (int, error) {
	id, err := u.rep.Create(ctx, user)
	if err != nil {
		return 0, fmt.Errorf("repository create: %w", err)
	}

	return id, nil
}
