package services

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/suite"

	"key-value-store/db"
	"key-value-store/services"
	"key-value-store/test"
)

type keyValueServiceSuite struct {
	suite.Suite
	db             db.DBTX
	service        *services.KeyValueService
	getByKey       func(key string) (*string, error)
	putKey         func(key string, value *string) error
	cleanKeyValues func()
}

func (s *keyValueServiceSuite) SetupSuite() {
	ctx := context.TODO()
	s.db = test.NewTestDB()
	s.service = services.NewKeyValueService(db.New(s.db))
	s.getByKey = func(key string) (*string, error) {
		return s.service.GetByKey(ctx, key)
	}
	s.putKey = func(key string, value *string) error {
		return s.service.PutKey(ctx, key, value)
	}
	s.cleanKeyValues = func() {
		_, err := s.db.ExecContext(ctx, `delete from key_values`)
		s.NoError(err)
	}
}

func (s *keyValueServiceSuite) AfterTest(_, _ string) {
	s.cleanKeyValues()
}

func Test_KeyValueService(t *testing.T) {
	suite.Run(t, new(keyValueServiceSuite))
}

func (s *keyValueServiceSuite) Test_GetByKey() {
	s.Run("when key doesn't exist in database", func() {
		_, err := s.getByKey("unknown")
		s.Error(err, sql.ErrNoRows)
	})

	s.Run("when key exists in database", func() {
		key := "key"
		value := "value"
		s.NoError(s.putKey(key, &value))

		expectedValue, err := s.getByKey(key)
		s.NoError(err)
		s.Equal(expectedValue, &value)
	})
}

func (s *keyValueServiceSuite) Test_PutKey_Success() {
	key := "key"
	s.NoError(s.putKey(key, nil))

	row := s.db.QueryRowContext(context.TODO(), `select 1 from key_values where k = $1`, key)
	var one int
	s.NoError(row.Scan(&one))
	s.Equal(one, 1)
}
