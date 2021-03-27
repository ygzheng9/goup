package matShortage

import (
	"fmt"
	"testing"

	"github.com/360EntSecGroup-Skylar/excelize"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/jmoiron/sqlx"
)

func init() {
	fmt.Println("init in matShortage_test.go")

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

	SetupDB(db)
	fmt.Printf("Setup u8 Complete. file: matShortage_test.go \n")
}

func Test_Stack(t *testing.T) {
	a := OneLevel{
		InvCode: "adsf",
		BaseQty: 23.12,
	}

	s := NewStack()

	fmt.Printf("isEmpty: %t\n", s.isEmpty())

	s.push(a)
	s.dump()

	s.push(a)
	s.dump()

	b := OneLevel{
		InvCode: "hello",
		BaseQty: 34.98,
	}

	s.push(b)
	s.dump()

	c := s.pop()
	fmt.Printf("c: %+v\n", c)
	s.dump()

	s.pop()
	s.pop()

	fmt.Printf("isEmpty: %t\n", s.isEmpty())

	d := s.pop()
	fmt.Printf("d: %+v\n", d)
}

func Test_loadBOM(t *testing.T) {
	items, err := loadBOM()
	if err != nil {
		t.Errorf("%+v\n", err)
	}

	fmt.Printf("total: %d\n", len(items))
	// for _, v := range items {
	// 	fmt.Printf("%s, %s", v.ParentInv, v.ChildInv)
	// }

	subs := findAllSubs("D13-C12040000-000", items)
	for _, v := range subs {
		fmt.Printf("%s, %.4f\n", v.InvCode, v.BaseQty)
	}

}

func Test_doProcess(t *testing.T) {
	inFile := "E:/99.localDev/easypy/u8/1.xlsx"
	interFile := "E:/99.localDev/easypy/u8/缺料表替代物料调整清单.xlsx"
	outFile := "E:/99.localDev/easypy/u8/1_out.xlsx"

	StartProcess(inFile, outFile, interFile)
}

func Test_xxx(t *testing.T) {
	// 列 从 0 开始
	fmt.Printf("%s%d ", excelize.ToAlphaString(0), 1)
	fmt.Printf("%s%d ", excelize.ToAlphaString(1), 2)

}

func Test_Interchagne(t *testing.T) {
	inFile := "E:/99.localDev/easypy/u8/缺料表替代物料调整清单.xlsx"

	xlsx, err := excelize.OpenFile(inFile)
	if err != nil {
		fmt.Println("输入文件错误。")
		fmt.Printf("%+v\n", err)
		return
	}

	items := loadInterchangeMat(xlsx)

	fmt.Printf("items: %d\n", len(items))
	for _, i := range items {
		fmt.Printf("%s -> %s\n", i.SrcMat, i.DestMat)
	}
}
