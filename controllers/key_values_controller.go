package controllers

import (
	"database/sql"
	"net/http"

	lru "github.com/hashicorp/golang-lru"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"

	"key-value-store/db"
	"key-value-store/errors/sql_errors"
	"key-value-store/services"
)

type KeyValueController struct {
	service *services.KeyValueService
}

func NewKeyValueController(cacheSize int, dbtx db.DBTX) *KeyValueController {
	cache, err := lru.New(cacheSize)
	if err != nil {
		panic(err)
	}
	return &KeyValueController{service: services.NewKeyValueService(cache, dbtx)}
}

func (self *KeyValueController) Get(ctx echo.Context) error {
	value, err := self.service.GetByKey(ctx.Request().Context(), ctx.Param("key"))
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "данный ключ отсутствует в базе")
		}
		return errors.WithStack(err)
	}
	return ctx.JSON(http.StatusOK, echo.Map{"value": value})
}

func (self *KeyValueController) Put(ctx echo.Context) error {
	value, err := self.validatePutBody(ctx)
	if err != nil {
		return err
	}
	if err := self.service.PutKey(ctx.Request().Context(), ctx.Param("key"), value); err != nil {
		if sql_errors.ErrUnique.Is(err) {
			return echo.NewHTTPError(http.StatusBadRequest, "перезапись ключа запрещена")
		}
		return err
	}
	return nil
}

func (self *KeyValueController) validatePutBody(ctx echo.Context) (*string, *echo.HTTPError) {
	body := echo.Map{}
	if err := ctx.Bind(&body); err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, err)
	}
	valueInterface, ok := body["value"]
	if !ok {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "value is not provided")
	}
	if valueInterface == nil {
		return nil, nil
	}
	value, ok := valueInterface.(string)
	if !ok {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "value must be string")
	}
	return &value, nil
}
