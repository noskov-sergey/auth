package users

import (
	"context"
	"fmt"
	"github.com/noskov-sergey/auth/internal/model"
)

func (u *useCase) Update(ctx context.Context, user model.User) error {
	err := u.rep.Update(ctx, user)
	if err != nil {
		return fmt.Errorf("repository update: %w", err)
	}

	return nil
}
