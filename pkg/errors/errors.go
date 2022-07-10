package errors

import (
	pkgerrors "github.com/pkg/errors"
)

const (
	notFoundErrMsg     = "not found error"
	internalErrMsg     = "internal error"
	invalidErrMsg      = "invalid request"
	permissionErrMsg   = "permission error"
	unauthorizedErrMsg = "unauthorized error"
)

var (
	NotFoundError     = notFoundError{}
	InternalError     = internalError{}
	InvalidError      = invalidError{}
	PermissionError   = permissionError{}
	UnauthorizedError = unauthorizedError{}
)

type notFoundError struct{}

func (notFoundError) Error() string { return notFoundErrMsg }

type internalError struct{}

func (internalError) Error() string { return internalErrMsg }

type invalidError struct{}

func (invalidError) Error() string { return invalidErrMsg }

type permissionError struct{}

func (permissionError) Error() string { return permissionErrMsg }

type unauthorizedError struct{}

func (unauthorizedError) Error() string { return unauthorizedErrMsg }

// ErrNotFound is an extend of codes.NotFound
func ErrNotFound(err error) error {
	return wrapErr(NotFoundError, err)
}

// ErrNotFoundf is an extend of codes.Invalid
func ErrNotFoundf(err error, format string, args ...interface{}) error {
	return wrapErrf(NotFoundError, err, format, args...)
}

func ErrNotFoundRaw(message string) error {
	return wrapErr(NotFoundError, pkgerrors.New(message))
}

// ErrInternal is an extend of codes.Internal
func ErrInternal(err error) error {
	return wrapErr(InternalError, err)
}

// ErrInternalf is an extend of codes.Internal
func ErrInternalf(err error, format string, args ...interface{}) error {
	return wrapErrf(InternalError, err, format, args...)
}

func ErrInternalRaw(message string) error {
	return wrapErr(InternalError, pkgerrors.New(message))
}

// ErrInvalid is an extend of codes.Invalid
func ErrInvalid(err error) error {
	return wrapErr(InvalidError, err)
}

// ErrInvalidf is an extend of codes.Invalid
func ErrInvalidf(err error, format string, args ...interface{}) error {
	return wrapErrf(InvalidError, err, format, args...)
}

func ErrInvalidRaw(message string) error {
	return wrapErr(InvalidError, pkgerrors.New(message))
}

// ErrPermissionDenied is an extend of codes.PermissionDenied
func ErrPermissionDenied(err error) error {
	return wrapErr(PermissionError, err)
}

// ErrPermissionDeniedf is an extend of codes.PermissionDenied
func ErrPermissionDeniedf(
	err error, format string, args ...interface{}) error {
	return wrapErrf(PermissionError, err, format, args...)
}

func ErrPermissionDeniedRaw(message string) error {
	return wrapErr(PermissionError, pkgerrors.New(message))
}

func ErrUnauthorized(err error) error {
	return wrapErr(UnauthorizedError, err)
}

func ErrUnauthorizedf(err error, format string, args ...interface{}) error {
	return wrapErrf(UnauthorizedError, err, format, args...)
}

func ErrUnauthorizedRaw(message string) error {
	return wrapErr(UnauthorizedError, pkgerrors.New(message))
}

func ErrorIs(err error, target error) bool {
	return pkgerrors.Is(err, target)
}

func wrapErr(wrapedErr error, err error) error {
	if err == nil {
		return pkgerrors.Wrap(wrapedErr, "")
	}
	return pkgerrors.Wrap(wrapedErr, err.Error())
}

func wrapErrf(
	wrapedErr error, err error, format string, args ...interface{}) error {
	e := pkgerrors.WithMessagef(wrapedErr, format, args...)
	if err == nil {
		return pkgerrors.Wrap(wrapedErr, e.Error())
	}
	return pkgerrors.Wrap(e, err.Error())
}
