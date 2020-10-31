package test

import (
	"context"
	"io/ioutil"
	"strings"
	"sync"

	"key-value-store/db"
	"key-value-store/utils"
)

const (
	DatabaseName = "db_test.sqlite"
	SchemaDir    = "schema/"
)

var once sync.Once
var inst db.DBTX

func NewTestDB() db.DBTX {
	root := utils.RootDir()
	once.Do(func() {
		dbtx, err := db.NewSQLite(root + DatabaseName)
		if err != nil {
			panic(err)
		}
		setupFromSchema(dbtx, root+SchemaDir)
		inst = dbtx
	})
	return inst
}

func setupFromSchema(dbtx db.DBTX, schemaPath string) {
	files, err := ioutil.ReadDir(schemaPath)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		b, err := ioutil.ReadFile(schemaPath + file.Name())
		if err != nil {
			panic(err)
		}
		DDL := strings.Split(string(b), "--")[0]
		_, err = dbtx.ExecContext(context.TODO(), DDL)
		if err != nil {
			panic(err)
		}
	}
}
