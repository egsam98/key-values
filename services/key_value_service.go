package services

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"

	"key-value-store/db"
)

// Сервис взаимодействия с хранилищем ключ-значения
type KeyValueService struct {
	q *db.Queries
}

func NewKeyValueService(queries *db.Queries) *KeyValueService {
	return &KeyValueService{q: queries}
}

func (self *KeyValueService) GetByKey(ctx context.Context, key string) (*string, error) {
	nullStr, err := self.q.GetByKey(ctx, key)
	if err != nil {
		return nil, err
	}
	if !nullStr.Valid {
		return nil, nil
	}
	return &nullStr.String, nil
}

func (self *KeyValueService) PutKey(ctx context.Context, key string, value *string) error {
	nullStr := sql.NullString{
		String: "",
		Valid:  false,
	}
	if value != nil {
		nullStr.String = *value
		nullStr.Valid = true
	}

	params := db.PutKeyParams{K: key, V: nullStr}
	if err := self.q.PutKey(ctx, params); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
