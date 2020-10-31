package services

import (
	"context"
	"database/sql"

	lru "github.com/hashicorp/golang-lru"
	"github.com/pkg/errors"

	"key-value-store/db"
)

const CacheSize = 1024

// Сервис взаимодействия с хранилищем ключ-значения
type KeyValueService struct {
	q     *db.Queries
	cache *lru.Cache
}

func NewKeyValueService(queries *db.Queries) *KeyValueService {
	cache, err := lru.New(CacheSize)
	if err != nil {
		panic(err)
	}
	return &KeyValueService{q: queries, cache: cache}
}

func (self *KeyValueService) GetByKey(ctx context.Context, key string) (*string, error) {
	if value, ok := self.cache.Get(key); ok {
		return value.(*string), nil
	}

	nullStr, err := self.q.GetByKey(ctx, key)
	if err != nil {
		return nil, err
	}

	var str *string
	if nullStr.Valid {
		str = &nullStr.String
	}
	go self.cache.Add(key, str)
	return str, nil
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

	if err := self.q.PutKey(ctx, db.PutKeyParams{K: key, V: nullStr}); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
