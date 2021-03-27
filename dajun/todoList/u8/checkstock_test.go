package u8

import (
	"bytes"
	"fmt"
	"strconv"
	"testing"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/jmoiron/sqlx"
)

func init() {
	fmt.Println("init in checkstock_test.go")

	// connString := fmt.Sprintf("odbc:server=%s;database=%s;user id=%s;password=%s;encrypt=disable",
	// 	"10.10.10.9", "UFDATA_001_2012", "sa", "Z&9I^0x)9D*6")

	connString := fmt.Sprintf("odbc:server=%s;database=%s;user id=%s;password=%s;encrypt=disable",
		"serveru8", "UFDATA_128_2012", "sa", "P@ss12345")

	db, err := sqlx.Open("mssql", connString)
	if err != nil {
		fmt.Printf("\nerr: %#v\n", err)
		fmt.Printf("can not connect u8")

		return
	}

	SetupDB(db)
	fmt.Printf("Setup u8 Complete. file: checkstock_test.go \n")
}

func Test_init(t *testing.T) {
	fmt.Println("ok")
}

func Test_DoCheckStock(t *testing.T) {
	qty, err := DoCheckStock("999-002000011-000", ">= 500")
	if err != nil {
		t.Errorf("\n%#v\n", err)
	}

	fmt.Printf("qty: %f\n", qty)
}

func Test_un(t *testing.T) {
	arr := []uint8{55, 54, 48, 48, 48, 46, 48, 48, 48, 48, 48, 48}

	var buffer bytes.Buffer
	for _, v := range arr {
		s := rune(v)
		buffer.WriteString(string(s))
	}

	qtyStr := buffer.String()
	fmt.Println(qtyStr)

	result, err := strconv.ParseFloat(qtyStr, 64)
	if err != nil {
		t.Errorf("\n%#v\n", err)
	}

	fmt.Printf("%f", result)

}

func Test_GetCurrentStock(t *testing.T) {
	// codes := []string{"D17-C12342021-000", "D16-S10701012-000", "999-007030106-000", "999-007010015-000"}
	codes := []string{"D17-C15342035-000"}
	items, err := GetCurrentStock(codes)
	if err != nil {
		t.Errorf("\n%#v\n", err)
	}
	for _, i := range items {
		fmt.Printf("%+v\n", i)
	}
}
