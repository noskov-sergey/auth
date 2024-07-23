package users

import (
	"context"
	"fmt"
	"github.com/noskov-sergey/auth/internal/model"
)

func (u *UseCase) Get(ctx context.Context, filter model.UserFilter) (*model.User, error) {
	user, err := u.rep.Get(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("repository get: %w", err)
	}

	return user, nil
}
