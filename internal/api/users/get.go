package users

import (
	"context"
	"fmt"
	"github.com/noskov-sergey/auth/internal/converter"
	desc "github.com/noskov-sergey/auth/pkg/user_v1"
)

func (i *Implementation) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	u, err := i.usecase.Get(ctx, converter.ToUserFilterFromUser(req))
	if err != nil {
		return nil, fmt.Errorf("usecase get: %w", err)
	}

	return &desc.GetResponse{
		User: converter.ToUserFromService(u),
	}, nil
}
