package models

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"pickup/dajun/todoList/u8"
	"testing"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/jmoiron/sqlx"
)

func Test_ScanAllTodos(t *testing.T) {
	err := ScanAllTodos()
	if err != nil {
		t.Errorf("ScanAllTodos: %#v\n", err)
		return
	}
}

func Test_Path(t *testing.T) {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	fmt.Printf("path: " + path)
}

func Test_ScanEvent(t *testing.T) {
	err := ScanEvent(3)
	if err != nil {
		t.Errorf("ScanAllTodos: %#v\n", err)
		return
	}
}

func Test_print(t *testing.T) {
	uploadDir := "C:/tmp/attachments"
	a := fmt.Sprintf("\n%s/at_%d", uploadDir, 2)
	t.Log(a)
}

func Test_ScanMatRule(t *testing.T) {
	err := ScanMatRule(0)
	if err != nil {
		t.Errorf("ScanMatRule: %#v\n", err)
		return
	}
}

func Test_ListMatNotice(t *testing.T) {
	connString := fmt.Sprintf("odbc:server=%s;database=%s;user id=%s;password=%s;encrypt=disable",
		"10.10.10.9", "UFDATA_001_2012", "sa", "Z&9I^0x)9D*6")

	db, err := sqlx.Open("mssql", connString)
	if err != nil {
		fmt.Printf("\nerr: %#v\n", err)
		fmt.Printf("can not connect u8")

		return
	}

	u8.SetupDB(db)

	items, err := ListMatNotice()
	if err != nil {
		t.Errorf("ScanMatRule: %#v\n", err)
		return
	}
	for _, i := range items {
		fmt.Printf("%+v\n", i)
	}
}
