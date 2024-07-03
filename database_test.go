package tools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMysqlDatabase(t *testing.T) {
	//test connect to local database
	cfg := DBConfiguration{
		DbType:         Mysql,
		Host:           "127.0.0.1",
		Port:           "3306",
		DBName:         "tgdata",
		Username:       "root",
		Password:       "root",
		Logging:        true,
		SessionName:    "test",
		ConnectTimeOut: 30,
		MaxOpenConn:    30,
		MaxIdleConn:    30,
		Schema:         "public",
	}
	conn, _ := ConnectDB(cfg)
	defer conn.Close()

	err := conn.Ping()
	assert.Nil(t, err, "error should nil")

	//test create logger
	logger := CreateLogger(true)
	assert.NotNil(t, logger, "object should initiated")

	//setup gorm db
	gormDB, err := NewGormDB(Mysql, conn, logger)
	assert.Nil(t, err, "error should nil on gorm db creation")

	assert.NotNil(t, gormDB, "gorm database should initiated")
}

func TestPostgreDatabase(t *testing.T) {
	//test connect to local database
	cfg := DBConfiguration{
		DbType:         Postgresql,
		Host:           "127.0.0.1",
		Port:           "5432",
		DBName:         "postgres",
		Username:       "postgres",
		Password:       "postgres",
		Logging:        true,
		SessionName:    "test",
		ConnectTimeOut: 30,
		MaxOpenConn:    30,
		MaxIdleConn:    30,
		Schema:         "public",
	}
	conn, _ := ConnectDB(cfg)
	defer conn.Close()

	err := conn.Ping()
	assert.Nil(t, err, "error should nil")

	assert.Nil(t, err, "error should nil on local database connection")

	//test create logger
	logger := CreateLogger(true)
	assert.NotNil(t, logger, "object should initiated")

	//setup gorm db
	gormDB, err := NewGormDB(Postgresql, conn, logger)
	assert.Nil(t, err, "error should nil on gorm db creation")

	assert.NotNil(t, gormDB, "gorm database should initiated")
}

func TestFailConnection(t *testing.T) {
	//test connect to local database
	cfg := DBConfiguration{
		DbType:         Postgresql,
		Host:           "127.0.0.1",
		Port:           "5432",
		DBName:         "postgres",
		Username:       "",
		Password:       "",
		Logging:        true,
		SessionName:    "test",
		ConnectTimeOut: 30,
		MaxOpenConn:    30,
		MaxIdleConn:    30,
		Schema:         "public",
	}
	conn, _ := ConnectDB(cfg)
	err := conn.Ping()
	assert.NotNil(t, err, "error should exist because fail connect to database")
}
