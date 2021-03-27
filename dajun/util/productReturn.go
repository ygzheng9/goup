package util

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"html/template"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/xuri/excelize"
)

// GenProcutReturnSQL 从文件读取数据，生成 sql 脚本
func GenProcutReturnSQL(sqlTemplate, dataFile string) {
	// load template
	// sqlTemplate := "./productReturn.tmpl"
	t1, err := loadTemplate(sqlTemplate)
	checkError(err, "loadTemplate")

	// load from file
	// dataFile := "./input.xlsx"
	// file, err := os.Open(dataFile)
	// checkError(err, "Open File")

	// items, err := loadExcel(file)
	// checkError(err, "loadData")

	// dataFile := "./productReturn.csv"
	file, err := os.Open(dataFile)
	checkError(err, "Open File")

	items, err := loadCSV(file)
	checkError(err, "loadCSV")

	// 对 file 的每一行，生成
	var total bytes.Buffer
	for _, item := range items {
		var buf bytes.Buffer
		if err := t1.Execute(&buf, item); err != nil {
			panic(err)
		}

		total.Write(buf.Bytes())
	}

	fmt.Printf("total: %d\n", len(items))
	fmt.Printf("%s\n", total.String())
}

func checkError(e error, source string) {
	if e != nil {
		fmt.Printf("ERR from: %s\n", source)
		panic(e)
	}
}

// 模板文件，和可执行文件，放到相同的目录下
func loadTemplate(fileName string) (*template.Template, error) {
	t1, err := template.New(fileName).ParseFiles(fileName)
	if err != nil {
		fmt.Printf("parse template error: %+v\n", err)
		return nil, err
	}
	return t1, nil
}

// ItemT 和上载文件的列对应
type ItemT struct {
	OrderNum string `json:"orderNum"`
	RowNo    string `json:"rowNo"`
	Quantity string `json:"quantity"`
}

// padLeft 左边补零
func padLeft(input, left string) string {
	l := strings.TrimSpace(left)
	i := strings.TrimSpace(input)
	start := len(i)

	tmp := l + i
	return tmp[start:]
}

// 从 csv 中读取文件
func loadCSV(file io.Reader) ([]ItemT, error) {
	// 从文件中读取
	r := csv.NewReader(file)

	items := []ItemT{}
	line := 0
	for {
		line++
		record, err := r.Read()
		// 第一行是头信息
		if line == 1 {
			continue
		}

		// 先检查 EOF
		if err == io.EOF {
			break
		}

		// 再检查不等于 nil
		if err != nil {
			return nil, err
		}

		fmt.Println(record)

		// 记录信息
		item := ItemT{
			OrderNum: padLeft(record[0], "0000000000"),
			RowNo:    strings.TrimSpace(record[1]),
			Quantity: strings.TrimSpace(record[2]),
		}

		// 读到了空值，表示结束
		if len(item.RowNo) == 0 || len(item.Quantity) == 0 {
			break
		}

		items = append(items, item)
	}

	return items, nil
}

// 解析上载的文件
func loadExcel(file io.Reader) ([]ItemT, error) {
	items := []ItemT{}

	// 从文件中读取
	xlsx, err := excelize.OpenReader(file)
	if err != nil {
		return nil, err
	}

	// sheet 的名字是固定的，只能是 Sheet1
	index := xlsx.GetSheetIndex("Sheet1")

	// Get all the rows in a sheet.
	rows := xlsx.GetRows("sheet" + strconv.Itoa(index))
	// fmt.Printf("total rows: %d\n", len(rows))

	for idx, row := range rows {
		// 第一行，表头，跳过
		if idx == 0 {
			continue
		}

		// 单据号 行号 数量
		if len(row[0]) == 0 || len(row[1]) == 0 || len(row[2]) == 0 {
			break
		}

		i := ItemT{
			OrderNum: padLeft(row[0], "0000000000"),
			RowNo:    strings.TrimSpace(row[1]),
			Quantity: strings.TrimSpace(row[2]),
		}

		items = append(items, i)
	}

	return items, nil
}
