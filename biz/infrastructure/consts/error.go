package consts

import (
	"google.golang.org/grpc/status"
)

var (
	ErrNotFound        = status.Error(10001, "no such element")
	ErrInvalidObjectId = status.Error(10002, "invalid objectId")
)
