package main

import (
	"fmt"

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
	inFile := "./forecast.xlsx"
	outFile := "./forecast_out.xlsx"

	from, to, err := matShortage.GetDuration(inFile)
	if err != nil {
		fmt.Printf("输入文件错误：%s", inFile)
		fmt.Printf("err: %+v\n", err)
		return
	}

	fmt.Printf("from: %s, to: %s\n", from, to)
	matShortage.CalcForecastDiff(from, to, outFile)
}
