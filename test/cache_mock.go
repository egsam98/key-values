package test

import (
	"github.com/hashicorp/golang-lru/simplelru"
	"github.com/stretchr/testify/mock"
)

var _ simplelru.LRUCache = (*CacheMock)(nil)

// Заглушка для LRU Cache
type CacheMock struct {
	mock.Mock
}

func (c *CacheMock) Add(key, value interface{}) bool {
	return false
}

func (c *CacheMock) Get(key interface{}) (value interface{}, ok bool) {
	args := c.Called(key)
	return args.Get(0), args.Bool(1)
}

func (c *CacheMock) Contains(key interface{}) (ok bool) {
	panic("implement me")
}

func (c *CacheMock) Peek(key interface{}) (value interface{}, ok bool) {
	panic("implement me")
}

func (c *CacheMock) Remove(key interface{}) bool {
	panic("implement me")
}

func (c *CacheMock) RemoveOldest() (interface{}, interface{}, bool) {
	panic("implement me")
}

func (c *CacheMock) GetOldest() (interface{}, interface{}, bool) {
	panic("implement me")
}

func (c *CacheMock) Keys() []interface{} {
	panic("implement me")
}

func (c CacheMock) Len() int {
	panic("implement me")
}

func (c *CacheMock) Purge() {
	panic("implement me")
}

func (c *CacheMock) Resize(int) int {
	panic("implement me")
}
