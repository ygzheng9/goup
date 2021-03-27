package config

import (
	"fmt"
	"log"
	"testing"

	_ "github.com/go-sql-driver/mysql"

	"github.com/jmoiron/sqlx"
)

func Test_LoadConfiguration(t *testing.T) {
	appConfig, err := LoadConfiguration("C:/localWork/goTest/todoList.json")
	if err != nil {
		t.Errorf("%+v\n", err)
		return
	}

	fmt.Printf("%+v\n", appConfig)
}

func Test_Conn(t *testing.T) {
	// 连接数据库
	connStr := fmt.Sprintf("%s:%s@(%s)/%s?parseTime=true",
		"demo", "1234", "127.0.0.1", "world")

	// fmt.Println(connStr)
	db, err := sqlx.Open("mysql", connStr)

	if err != nil {
		panic(err)
	}
	if db == nil {
		log.Printf("db is nil.\n")
	}

	// Test the connection to the database
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("connected.")
}
