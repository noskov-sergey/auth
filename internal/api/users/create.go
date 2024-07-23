package users

import (
	"context"
	"fmt"
	"github.com/noskov-sergey/auth/internal/converter"
	desc "github.com/noskov-sergey/auth/pkg/user_v1"
)

func (i *Implementation) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	id, err := i.usecase.Create(ctx, converter.ToServiceFromUser(req))
	if err != nil {
		return nil, fmt.Errorf("usecase create: %w", err)
	}

	return &desc.CreateResponse{
		Id: int64(id),
	}, nil
}
