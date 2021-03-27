package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path"
	"strings"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/jmoiron/sqlx"

	"pickup/dajun/todoList/u8/matShortage"
)

func init() {
	// fmt.Println("init in main.go")

	connString := fmt.Sprintf("odbc:server=%s;database=%s;user id=%s;password=%s;encrypt=disable",
		"10.10.10.5", "UFDATA_001_2012", "sa", "Z&9I^0x)9D*6")

	// connString := fmt.Sprintf("odbc:server=%s;database=%s;user id=%s;password=%s;encrypt=disable",
	// 	"serveru8", "UFDATA_128_2012", "sa", "P@ss12345")

	db, err := sqlx.Open("mssql", connString)
	if err != nil {
		fmt.Printf("\nerr: %#v\n", err)
		fmt.Printf("can not connect u8")

		return
	}

	matShortage.SetupDB(db)
	// fmt.Printf("Setup u8 Complete. file: main.go \n")
}

func main() {
	inFile := flag.String("i", "./1.xlsx", "请输入文件名")
	flag.Parse()

	file := *inFile
	if file[0] != '.' {
		file = "./" + file
	}
	fmt.Println("输入文件: ", file)

	//获取文件名带后缀
	nameWithSuffix := path.Base(file)
	//获取文件后缀
	suffix := path.Ext(nameWithSuffix)
	// fmt.Println("suffix =", suffix)
	//获取文件名
	nameOnly := strings.TrimSuffix(nameWithSuffix, suffix)
	// fmt.Println("nameOnly =", nameOnly)

	outFile := fmt.Sprintf("./%s_out%s", nameOnly, suffix)
	fmt.Println("输出文件: ", outFile)

	interFile := "./缺料表替代物料调整清单.xlsx"
	matShortage.StartProcess(file, outFile, interFile)

	fmt.Print("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
