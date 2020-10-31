package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"

	db2 "key-value-store/db"
	"key-value-store/server"
	"key-value-store/test"
)

const ADDR = "localhost:8080"

type keyValuesSuite struct {
	suite.Suite
	server *server.EchoServer
	db     db2.DBTX
}

func (s *keyValuesSuite) SetupSuite() {
	s.db = test.NewTestDB()
	s.server = server.StartEchoServer(ADDR, test.NewTestDB(), false)
}

func Test_KeyValuesController(t *testing.T) {
	suite.Run(t, new(keyValuesSuite))
}

func (s *keyValuesSuite) Test_Put_KeyValue() {
	s.Run("when no key provided", func() {
		body, status := s.request("PUT", "")
		s.Equal(status, 404)
		s.Equal(body["message"], "Not Found")
	})

	s.Run("when no value provided", func() {
		body, status := s.request("PUT", "key")
		s.Equal(status, 400)
		s.Equal(body["message"], "value is not provided")
	})

	s.Run("when key is int", func() {
		body, status := s.request("PUT", "key", echo.Map{"value": 1})
		s.Equal(status, 400)
		s.Equal(body["message"], "value must be string")
	})
}

func (s *keyValuesSuite) Test_Get_KeyValue() {
	s.Run("when key is not provided", func() {
		body, status := s.request("GET", "")
		s.Equal(status, 404)
		s.Equal(body["message"], "Not Found")
	})

	s.Run("when key doesn't exist in database", func() {
		body, status := s.request("GET", "UNKNOWN")
		s.Equal(status, 404)
		s.Equal(body["message"], "данный ключ отсутствует в базе")
	})

	s.Run("when key exists in database", func() {
		ctx := context.TODO()
		key := "key"
		value := "value"
		_, err := s.db.ExecContext(ctx, `insert into key_values (k, v) values ($1, $2)`, key, value)
		s.NoError(err)

		body, status := s.request("GET", key)
		s.Equal(status, 200)
		s.Equal(body["value"], value)

		s.T().Cleanup(func() {
			_, err := s.db.ExecContext(ctx, `delete from key_values`)
			s.NoError(err)
		})
	})
}

func (s *keyValuesSuite) request(method, key string, requestBody ...echo.Map) (echo.Map, int) {
	if len(requestBody) == 0 {
		requestBody = []echo.Map{nil}
	}

	b, err := json.Marshal(requestBody[0])
	s.NoError(err)

	req, err := http.NewRequest(method, "http://"+ADDR+"/kv/"+key, bytes.NewReader(b))
	s.NoError(err)

	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	s.NoError(err)

	var body echo.Map
	s.NoError(json.NewDecoder(res.Body).Decode(&body))
	return body, res.StatusCode
}
