package users

import (
	"context"
	"fmt"
	"github.com/noskov-sergey/auth/internal/converter"
	desc "github.com/noskov-sergey/auth/pkg/user_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	err := i.usecase.Update(ctx, converter.ToUpdateServiceFromUser(req))
	if err != nil {
		return nil, fmt.Errorf("usecase update: %w", err)
	}

	return &emptypb.Empty{}, nil
}
