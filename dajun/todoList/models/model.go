package models

import "github.com/jmoiron/sqlx"

var db *sqlx.DB
var templateDir string

// Setup 设置参数
func Setup(dbConn *sqlx.DB, templDir string) {
	db = dbConn
	templateDir = templDir
}
