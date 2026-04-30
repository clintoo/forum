package services

import "fmt"

type ValidationError struct{ Msg string }

type ConflictError struct{ Msg string }

type NotFoundError struct{ Msg string }

type InternalError struct {
	Err error
	Msg string
}

func (e *NotFoundError) Error() string { return e.Msg }

func (e *ConflictError) Error() string { return e.Msg }

func (e *ValidationError) Error() string { return e.Msg }

func (e *InternalError) Error() string {
	if e == nil {
		return ""
	}

	if e.Msg == "" {
		if e.Err != nil {
			return e.Err.Error()
		}

		return "internal error"
	}

	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Msg, e.Err)
	}

	return e.Msg
}

func NewValidationError(msg string) error {
	return &ValidationError{Msg: msg}
}

func NewConflictError(msg string) error {
	return &ConflictError{Msg: msg}
}

func NewNotFoundError(msg string) error {
	return &NotFoundError{Msg: msg}
}

func NewInternalError(msg string, err error) error {
	return &InternalError{Err: err, Msg: msg}
}
