package models

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"time"
)

// DataHead 测试报告头
type DataHead struct {
	ID          int64  `db:"ID" json:"id" form:"id"`
	Type        string `db:"TYP" json:"type" form:"type"`
	Product     string `db:"PRODCT" json:"product" form:"product"`
	ProductSN   string `db:"PRODCT_SN" json:"product_sn" form:"product_sn"`
	TestingDate string `db:"TESTING_DTE" json:"testing_date" form:"testing_date"`
	EquipNo     string `db:"EQUIP_NUM" json:"equip_no" form:"equip_no"`
	Operator    string `db:"OP" json:"operator" form:"operator"`
	Result      string `db:"RESULT" json:"result" form:"result"`
	CreateDate  string `db:"CRE_DTE" json:"create_date" form:"create_date"`
}

// DataItem 测试项目
type DataItem struct {
	ID         int64  `db:"ID" json:"id" form:"id"`
	HeadID     int64  `db:"HEAD_ID" json:"head_id" form:"head_id"`
	Seq        string `db:"SEQ" json:"seq" form:"seq"`
	Name       string `db:"NME" json:"name" form:"name"`
	Celling    string `db:"CELLING" json:"celling" form:"celling"`
	Floor      string `db:"FLR" json:"floor" form:"floor"`
	Unit       string `db:"UNIT" json:"unit" form:"unit"`
	Value      string `db:"VAL" json:"value" form:"value"`
	Result     string `db:"RESULT" json:"result" form:"result"`
	CreateDate string `db:"CRE_DTE" json:"create_date" form:"create_date"`
}

// TestingDataParse 把测试数据保存到数据库
func TestingDataParse(b []byte) error {
	// 从文件中读取
	r := csv.NewReader(bytes.NewReader(b))

	header := DataHead{}
	items := []DataItem{}

	line := 0
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("cvs read error: %+v\n", err)
			break
		}

		// fmt.Println(record)

		// 头信息有两行
		if line == 0 {
			// 第一行，头信息
			header.Type = record[0]
			header.Product = record[2]
			header.ProductSN = record[4]
			header.TestingDate = record[6]
		} else if line == 1 {
			// 第二行，头信息
			header.EquipNo = record[2]
			header.Operator = record[4]
			header.Result = record[6]
		} else if line == 2 {
			// 测试项目的头，跳过
		} else {
			// 行信息
			item := DataItem{
				Seq:     record[0],
				Name:    record[1],
				Celling: record[2],
				Floor:   record[3],
				Unit:    record[4],
				Value:   record[5],
				Result:  record[6],
			}
			items = append(items, item)
		}

		line++
	}

	fmt.Printf("test data: %d lines.\n", line)

	now := time.Now().Format("2006-01-02 15:04:05")
	header.CreateDate = now

	// 保存头信息
	insHead := ` insert into T_TESTINGHD (TYP,PRODCT,PRODCT_SN,TESTING_DTE,EQUIP_NUM,OP,RESULT,CRE_DTE)
						values (:TYP,:PRODCT,:PRODCT_SN,:TESTING_DTE,:EQUIP_NUM,:OP,:RESULT,:CRE_DTE)`
	res, err := db.NamedExec(insHead, header)
	if err != nil {
		fmt.Printf("head error: %+v\n", err)
		return err
	}

	// 取得测试头的ID
	headID, err := res.LastInsertId()
	if err != nil {
		fmt.Printf("headID error: %+v\n", err)
		return err
	}

	// 保存行信息
	insItem := ` insert into T_TESTINGITM (HEAD_ID,SEQ,NME,CELLING,FLR,UNIT,VAL,RESULT,CRE_DTE)
				values (:HEAD_ID,:SEQ,:NME,:CELLING,:FLR,:UNIT,:VAL,:RESULT,:CRE_DTE) `
	for _, i := range items {
		i.HeadID = headID
		i.CreateDate = now

		_, err := db.NamedExec(insItem, i)
		if err != nil {
			fmt.Printf("item error: %+v\n", err)
			return err
		}
	}

	return nil
}

// TestingDataHeadFind 取得测试数据的头信息
func TestingDataHeadFind(cond string) ([]DataHead, error) {
	cmd := `select ID, ifnull(TYP,'') TYP, ifnull(PRODCT,'') PRODCT, ifnull(PRODCT_SN,'') PRODCT_SN, ifnull(TESTING_DTE,'') TESTING_DTE, ifnull(EQUIP_NUM,'') EQUIP_NUM, ifnull(OP,'') OP, ifnull(RESULT,'') RESULT, ifnull(CRE_DTE,'') CRE_DTE from T_TESTINGHD `

	if cond != "" {
		cmd = cmd + cond
	}

	// fmt.Printf("cmd: %s\n", cmd)

	items := []DataHead{}
	err := db.Select(&items, cmd)
	return items, err
}

// TestingDataItemFind 根据 headid 取得测试报告明细
func TestingDataItemFind(headID int) ([]DataItem, error) {
	cmd := `select ID, ifnull(HEAD_ID,'') HEAD_ID, ifnull(SEQ,'') SEQ, ifnull(NME,'') NME, ifnull(CELLING,'') CELLING, ifnull(FLR,'') FLR, ifnull(UNIT,'') UNIT, ifnull(VAL,'') VAL, ifnull(RESULT,'') RESULT, ifnull(CRE_DTE,'') CRE_DTE
	from T_TESTINGITM where HEAD_ID = ` + strconv.Itoa(headID)

	items := []DataItem{}
	err := db.Select(&items, cmd)
	return items, err
}

// TestingDataHeadDelete 根据对象的 ID 删除
func TestingDataHeadDelete(id int) error {
	headID := strconv.Itoa(id)

	// 删除头
	sqlCmd := `delete from T_TESTINGHD where ID=` + headID
	db.Exec(sqlCmd)

	// 删除行
	sqlCmd = `delete from T_TESTINGITM where HEAD_ID=` + headID
	db.Exec(sqlCmd)

	return nil
}
