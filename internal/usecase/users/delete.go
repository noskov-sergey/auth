package users

import (
	"context"
	"fmt"
)

func (u *useCase) Delete(ctx context.Context, id int) error {
	err := u.rep.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("repository delete: %w", err)
	}

	return nil
}
