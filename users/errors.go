package users

import "errors"

//UsersDomainError - domain errors for user logic validations
type UsersDomainError struct {
	err error
}

var (
	ErrNotFound          error = NewDomainError("user not found")
	ErrInternalError     error = NewDomainError("error")
	ErrInvalidData       error = NewDomainError("invalid data")
	ErrUserAlreadyExists error = NewDomainError("already exists")
)

func NewDomainError(msg string) UsersDomainError {
	return UsersDomainError{err: errors.New(msg)}
}

func (e UsersDomainError) Error() string {
	return e.err.Error()
}
