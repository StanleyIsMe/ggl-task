package repository

import (
	"errors"
)

var ErrDataNotFound = errors.New("data not found")
var ErrInvalidData = errors.New("invalid data")
