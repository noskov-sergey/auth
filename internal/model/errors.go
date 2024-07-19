package model

import "errors"

var (
	NotFoundErr = errors.New("not found")

	UnauthorizedErr     = errors.New("users is not authorized")
	InvalidArgumentsErr = errors.New("invalid arguments")
	PermissionDeniedErr = errors.New("permission denied")
)
