package models

import (
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func init() {
	// 设置数据库连接
	db, err := sqlx.Open("mysql", "root:mysql@/world?parseTime=true")
	if err != nil {
		fmt.Printf("can not connect db")
		return
	}

	Setup(db, "C:/localWork/goTest/templates")
	fmt.Printf("SetupDB Complete. file: user_test.go \n")
}

func TestFindAll(t *testing.T) {
	u := User{}
	items, err := u.FindAll()
	if err != nil {
		fmt.Printf("FindAll error: %+v\n", err)
	}

	for idx, i := range items {
		fmt.Printf("%d - %+v\n", idx, i)
	}
}

func TestPrint(t *testing.T) {
	sql := "select * from a"
	name := "hah"
	sql = fmt.Sprintf("%s where 1 = 1 and code like '%%%s%%'", sql, name)
	fmt.Println(sql)
}

func Test_ValidateLogin(t *testing.T) {
	check, err := ValidateLogin("zhengyg@dajuntech.com", "12345")
	if err != nil {
		fmt.Printf("FindAll error: %+v\n", err)
	}
	fmt.Printf("check: %+v\n", check)
}

func Test_ChangePassword(t *testing.T) {
	err := ChangePassword("zhengyg@dajuntech.com", "12345")
	if err != nil {
		fmt.Printf("FindAll error: %+v\n", err)
	}
	fmt.Printf("done.\n")
}
