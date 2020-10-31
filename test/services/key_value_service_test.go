package services

import (
	"context"
	"database/sql"
	"strconv"
	"sync"
	"testing"

	lru "github.com/hashicorp/golang-lru"
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
	cache, err := lru.New(1024)
	s.NoError(err)
	s.service = services.NewKeyValueService(cache, s.db)
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

func Benchmark_GetByKey_Caching(b *testing.B) {
	cache := &test.CacheMock{}
	testDB := test.NewTestDB()
	service := services.NewKeyValueService(cache, testDB)

	b.Run("when success", func(b *testing.B) {
		key := "key"
		cache.On("Get", key).Return(new(string), true)
		b.ResetTimer()
		_, _ = service.GetByKey(context.TODO(), key)
	})

	b.Run("when page fault occured", func(b *testing.B) {
		ctx := context.TODO()
		cache.On("Get", "").Return(nil, false)
		randomizeKeyValues(service, 500)

		b.ResetTimer()
		_, _ = service.GetByKey(ctx, "")
		b.StopTimer()

		b.Cleanup(func() {
			_, _ = testDB.ExecContext(ctx, `delete from key_values`)
		})
	})
}

func randomizeKeyValues(service *services.KeyValueService, count int) {
	ctx := context.TODO()
	var wg sync.WaitGroup
	wg.Add(count)
	for i := 0; i < count; i++ {
		go func() {
			_ = service.PutKey(ctx, strconv.Itoa(i), nil)
			wg.Done()
		}()
	}
	wg.Wait()
}
