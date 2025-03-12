package usecase

import (
	"fmt"
	"net/http"
)

//nolint:revive
type UseCaseError interface {
	ErrorCode() string
	ErrorMsg() string
	Error() string
	HTTPStatusCode() int
}

type InternalServerError struct {
	Err error
}

func (e InternalServerError) ErrorCode() string {
	return "INTERNAL_SERVER_ERROR"
}

func (e InternalServerError) ErrorMsg() string {
	return "Internal Server Error"
}

func (e InternalServerError) Error() string {
	return fmt.Sprintf("internal server error: %v", e.Err)
}

func (e InternalServerError) HTTPStatusCode() int {
	return http.StatusInternalServerError
}

type DuplicatedResourceError struct {
	Name     string
	Resource string
}

func (e DuplicatedResourceError) ErrorCode() string {
	return "DUPLICATED_RESOURCE"
}

func (e DuplicatedResourceError) ErrorMsg() string {
	return fmt.Sprintf("%s %s already exists", e.Resource, e.Name)
}

func (e DuplicatedResourceError) Error() string {
	return fmt.Sprintf("%s %s already exists", e.Resource, e.Name)
}

func (e DuplicatedResourceError) HTTPStatusCode() int {
	return http.StatusBadRequest
}

type NotFoundError struct {
	Resource string
	ID       interface{}
}

func (e NotFoundError) ErrorCode() string {
	return "NOT_FOUND"
}

func (e NotFoundError) ErrorMsg() string {
	return fmt.Sprintf("%s %v not found", e.Resource, e.ID)
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("%s %v not found", e.Resource, e.ID)
}

func (e NotFoundError) HTTPStatusCode() int {
	return http.StatusNotFound
}
