package config

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
)

// AppConfig 程序的配置
type AppConfig struct {
	Port        string `json:"port"`
	Release     bool   `json:"release"`
	Schedule    bool   `json:"schedule"`
	TemplateDir string `json:"templateDir"`
	DB          struct {
		User     string `json:"user"`
		Password string `json:"password"`
		Server   string `json:"server"`
		Schema   string `json:"schema"`
	} `json:"db"`
	U8 struct {
		User     string `json:"user"`
		Password string `json:"password"`
		Server   string `json:"server"`
		Schema   string `json:"schema"`
		IsAvl    bool   `json:"available"`
	} `json:"u8"`
}

const (
	passwordKey = "MYGOD"
)

var appConfig AppConfig

// LoadConfiguration 加载配置文件
func LoadConfiguration(file string) (AppConfig, error) {
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		return appConfig, err
	}
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&appConfig)

	// 文件中存的密文
	var data []byte
	data, _ = base64.StdEncoding.DecodeString(appConfig.DB.Password)
	appConfig.DB.Password = string(data)

	data, _ = base64.StdEncoding.DecodeString(appConfig.U8.Password)
	appConfig.U8.Password = string(data)

	return appConfig, err
}

// ConnectDB 初始化数据库连接
func ConnectDB() *sqlx.DB {
	// 连接数据库
	// connStr := fmt.Sprintf("%s:%s@%s/%s?parseTime=true",
	// 	appConfig.DB.User, appConfig.DB.Password, appConfig.DB.Server, appConfig.DB.Schema)

	connStr := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true",
		"dangerUser", "!@#456QWErty", "47.92.175.81:3366", "world")

	// connStr := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true",
	// 	"root", "mysql", "127.0.0.1:3306", "world")

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

	return db
}

// ConnectU8 初始化U8数据库连接
func ConnectU8() *sqlx.DB {
	// 连接数据库
	connString := fmt.Sprintf("odbc:server=%s;database=%s;user id=%s;password=%s;encrypt=disable",
		appConfig.U8.Server, appConfig.U8.Schema, appConfig.U8.User, appConfig.U8.Password)

	db, err := sqlx.Open("mssql", connString)
	if err != nil {
		panic(err)
	}
	if db == nil {
		log.Printf("u8 is nil.\n")
	}

	// Test the connection to the database
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	return db
}
