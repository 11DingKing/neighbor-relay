package domain

import "errors"

var (
	ErrNotFound     = errors.New("not found")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrConflict     = errors.New("conflict")
	ErrInvalid      = errors.New("invalid request")
	ErrExpired      = errors.New("expired")
	ErrCancelled    = errors.New("cancelled")
)

type ErrorCode string

const (
	CodeNotFound     ErrorCode = "not_found"
	CodeUnauthorized ErrorCode = "unauthorized"
	CodeForbidden    ErrorCode = "forbidden"
	CodeConflict     ErrorCode = "conflict"
	CodeInvalid      ErrorCode = "invalid_request"
	CodeInternal     ErrorCode = "internal_error"
)
