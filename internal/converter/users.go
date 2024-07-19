package converter

import (
	"github.com/noskov-sergey/auth/internal/model"
	desc "github.com/noskov-sergey/auth/pkg/user_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToUserFilterFromUser(req *desc.GetRequest) model.UserFilter {
	return model.UserFilter{
		UserID: int(req.GetId()),
	}
}

func ToServiceFromUser(req *desc.CreateRequest) model.User {
	return model.User{
		Name:  req.GetName(),
		Email: req.GetEmail(),
		Role:  int(req.GetRole()),
	}
}

func ToUpdateServiceFromUser(req *desc.UpdateRequest) model.User {
	return model.User{
		ID:    int(req.GetId()),
		Name:  req.GetName(),
		Email: req.GetMail(),
	}
}

func ToUserFromService(user *model.User) *desc.User {
	var updatedAt *timestamppb.Timestamp
	if user.UpdatedAt.Valid {
		updatedAt = timestamppb.New(user.UpdatedAt.Time)
	}

	return &desc.User{
		Id:        int64(user.ID),
		Name:      user.Name,
		Email:     user.Email,
		Role:      toEnum(user.Role),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: updatedAt,
	}
}

func toEnum(role int) *desc.Enum {
	var enumUser desc.Enum = 0
	var enumAdmin desc.Enum = 1
	switch role {
	case 1:
		return &enumAdmin
	default:
		return &enumUser
	}
}
