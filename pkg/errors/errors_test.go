package errors

import (
	stderrors "errors"
	"testing"

	pkgerrors "github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func createNotFoundWrapErr() error {
	return ErrNotFound(stderrors.New("not found"))
}

func TestErrors(t *testing.T) {
	t.Run("not found error test", func(t *testing.T) {
		err1 := stderrors.New("not found error")
		newErr := ErrNotFound(err1)
		assert.Equal(t, "not found error: not found error", newErr.Error())
	})
	t.Run("not found raw error test", func(t *testing.T) {
		err1 := stderrors.New("not found error")
		newErr := ErrNotFoundf(err1, "err1 error %s", "test")
		assert.Equal(t, "not found error: err1 error test: not found error", newErr.Error())
	})
	t.Run("nested cause test", func(t *testing.T) {
		err := createNotFoundWrapErr()
		switch pkgerrors.Cause(err).(type) {
		case notFoundError:
		default:
			t.Error("expect to equal InvalidError type")
		}
	})
	t.Run("nested not found error test", func(t *testing.T) {
		err := createNotFoundWrapErr()
		if !ErrorIs(err, NotFoundError) {
			t.Error("not equal")
		}
	})
}
