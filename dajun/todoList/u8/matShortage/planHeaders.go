package matShortage

import (
	"github.com/360EntSecGroup-Skylar/excelize"
)

// 读取 输入计划的 表头（第一行：日期，第二行，星期）

// 输入计划的表头
func getHeader(workbook *excelize.File) [][]string {
	//  取得计划的表头，后续写入缺料表
	sheetName := "Sheet2"
	// 本 worksheet 写入的 helper
	readHelper := sheetReader(workbook, sheetName)

	// 第一行、第二行的读取 helper
	getColumn := rowReader(readHelper, 1)
	getColumn2 := rowReader(readHelper, 2)

	//  表头有两行
	header := [][]string{}
	for col := 0; col < TotalColumn; col++ {
		idx := col + StartColumn

		h1 := getColumn(idx)
		h2 := getColumn2(idx)

		oneDay := []string{h1, h2}
		header = append(header, oneDay)
	}

	return header
}
