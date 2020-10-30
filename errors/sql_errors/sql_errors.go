// Ошибки SQL-запросов

package sql_errors

import (
	"errors"
	"strings"
)

var ErrUnique = &SQLError{errors.New("UNIQUE constraint failed")}

type SQLError struct {
	error
}

func (sqlErr *SQLError) Is(err error) bool {
	return strings.Contains(err.Error(), sqlErr.Error())
}
