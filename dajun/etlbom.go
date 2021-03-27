package main

import (
	"database/sql"
	"fmt"

	_ "github.com/denisenkom/go-mssqldb"
	"log"
)

func main() {
	connString := fmt.Sprintf("odbc:server=%s;database=%s;user id=%s;password=%s;encrypt=disable",
		"10.10.10.9", "UFDATA_702_2012", "sa", "Z&9I^0x)9D*6")

	db, err := sql.Open("mssql", connString)
	if err != nil {
		log.Fatal("Open connection failed:", err.Error())
		return
	}
	err = db.Ping()
	if err != nil {
		fmt.Print("PING:%s", err)
		return
	}

	fmt.Print("connection ok.")

	rows, err := db.Query("select a.BomId, b.ParentId partId from bom_bom a inner join bom_parent b on a.BomId=b.BomId")

	if err != nil {
		fmt.Println("query bom: ", err)
		return
	}

	for rows.Next() {
		var bomId int
		var partId int
		rows.Scan(&bomId, &partId)
		fmt.Printf("bomId: %d \t partId: %d\n", bomId, partId)
	}

	rows, err = db.Query("select a.BomId, a.componentId from bom_opcomponent a")
	if err != nil {
		fmt.Println("query comps: ", err)
		return
	}

	for rows.Next() {
		var bomId int
		var componentId int
		rows.Scan(&bomId, &componentId)
		fmt.Printf("bomId: %d \t componentId: %d\n", bomId, componentId)
	}

	rows, err = db.Query("select a.partId, a.InvCode from bas_part a")
	if err != nil {
		fmt.Println("query comps: ", err)
		return
	}

	for rows.Next() {
		var partId int
		var InvCode string
		rows.Scan(&partId, &InvCode)
		fmt.Printf("partId: %d \t InvCode: %s\n", partId, InvCode)
	}

}
