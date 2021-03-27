package models

import (
	"fmt"
	"log"
	"path"
	"testing"
)

func init() {
	fmt.Println("init in filemgmt_test.go")
	fmt.Println("测试时，只能初始化一次数据库连接")

	// // TODO: 如果设置数据库连接？
	// db, err := sqlx.Open("mysql", "root:mysql@/world?parseTime=true")
	// if err != nil {
	// 	fmt.Printf("can not connect db")
	// 	return
	// }

	// SetupDB(db)
}

func Test_Insert(t *testing.T) {
	var f FileMgmt

	f.AbsolutePath = "C:/localWork"
	f.ID = -1
	f.Depth = 0
	f.FileName = "localWork"
	f.UpdateUser = "zhengyg"
	f.IsDir = "true"
	f.Status = "A"
	f.Content = "项目根目录"

	err := f.Insert()
	fmt.Printf("%#v \n", err)
}

func Test_listFS(t *testing.T) {
	var f FileMgmt

	// f.AbsolutePath = "C:/localWork"
	// f.ID = 16
	// f.Depth = 2

	f.ID = 1

	items, err := f.FindSubs()

	if err != nil {
		log.Printf("%#v\n", err)
	}

	for _, i := range items {
		fmt.Printf("%#v \n", i)
	}

	t.Log("listFS")
}

func Test_removeFile(t *testing.T) {
	var i FileMgmt
	i.ID = 5
	err := i.RemoveFile()

	if err != nil {
		log.Printf("%#v\n", err)
	}
}

func Test_reNameFile(t *testing.T) {
	var i FileMgmt
	err := i.LoadByID(5)

	if err != nil {
		log.Printf("%#v\n", err)
	}

	err = i.RenameFile("babab")
	if err != nil {
		log.Printf("%#v\n", err)
	}
}

func Test_ModifyFile(t *testing.T) {
	var i FileMgmt
	err := i.LoadByID(5)

	if err != nil {
		log.Printf("%#v\n", err)
	}

	err = i.BackupFile()
	if err != nil {
		log.Printf("%#v\n", err)
	}

	err = i.BackupFile()
	if err != nil {
		log.Printf("%#v\n", err)
	}
}

func Test_FindHistory(t *testing.T) {
	var i FileMgmt
	err := i.LoadByID(5)

	if err != nil {
		log.Printf("%#v\n", err)
	}

	items, err := i.FindHistory()
	if err != nil {
		log.Printf("%#v\n", err)
	}

	fmt.Printf("%d \n", len(items))

	for _, i := range items {
		fmt.Printf("%#v \n", i)
	}
}

func Test_GetFullName(t *testing.T) {
	var f FileMgmt

	err := f.LoadByID(7)
	if err != nil {
		fmt.Printf("%#v\n", err)
	}

	full, err := f.GetFullName()
	if err != nil {
		fmt.Printf("%#v\n", err)
	}

	fmt.Printf("full: %s\n", full)
}

func Test_foo(t *testing.T) {
	msg := "good"
	t.Logf("Log: hello %s\n", msg)
	t.Errorf("Error: hello %s\n", msg)
}

func Test_changeFile(t *testing.T) {
	f := FileMgmt{}

	f.LoadByID(10)

	fullPath, err := f.GetFullName()
	if err != nil {
		fmt.Printf("%#v\n", err)
		return
	}
	t.Logf("\n fullname: %s - %s", path.Dir(fullPath), path.Base(fullPath))
}
