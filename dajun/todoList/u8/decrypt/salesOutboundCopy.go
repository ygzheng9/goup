// 复制 excel 文件内容

package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
)

// 解析上载的文件
func copyOutboundFile() error {
	// 从文件中读取
	// xlsx, err := excelize.OpenReader(file)

	// 文件名是固定死的
	xlsx, err := excelize.OpenFile("./Book1.xlsx")
	if err != nil {
		return err
	}

	rows := xlsx.GetRows("Sheet1")
	fmt.Printf("total rows: %d\n", len(rows))

	// 输出文件
	xlsxOut := excelize.NewFile()

	// 最大行数
	const maxRow = 1000

	for idx, row := range rows {
		// 料号、新批次号为空，表示是最后一行
		if idx >= maxRow || len(row[1]) == 0 || len(row[2]) == 0 {
			break
		}

		lineNum := strconv.Itoa(idx + 1)
		xlsxOut.SetCellValue("Sheet1", "A"+lineNum, strings.TrimSpace(row[0]))
		xlsxOut.SetCellValue("Sheet1", "B"+lineNum, strings.TrimSpace(row[1]))
		xlsxOut.SetCellValue("Sheet1", "C"+lineNum, strings.TrimSpace(row[2]))
	}

	// 文件名是固定死的
	err = xlsxOut.SaveAs("./Book1_out.xlsx")
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func main() {
	err := copyOutboundFile()
	if err != nil {
		fmt.Printf("err: %+v\n", err)
	}
}
