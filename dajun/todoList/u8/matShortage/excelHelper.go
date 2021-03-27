package matShortage

import (
	"fmt"

	"github.com/360EntSecGroup-Skylar/excelize"
)

// distString 字符串去重
func distString(items []string) []string {
	distInvs := []string{}
	for _, i := range items {
		flag := true
		for _, j := range distInvs {
			if i == j {
				flag = false
				break
			}
		}
		if flag {
			distInvs = append(distInvs, i)
		}
	}

	return distInvs
}

// 根据坐标，返回单元格名字
// r 从 1 开始，也即：1 代表 第1行，2 代表 第2行； A1, A2
// c 从 0 开始，也即 0 -> A, 1 -> B
func cellAxis(r int) func(c int) string {
	return func(c int) string {
		return fmt.Sprintf("%s%d", excelize.ToAlphaString(c), r)
	}
}

// sheet 的写入的辅助函数
type sheetWriterT func(axis string, value interface{})

func sheetWriter(workbook *excelize.File, sheetName string) sheetWriterT {
	return func(axis string, value interface{}) {
		workbook.SetCellValue(sheetName, axis, value)
	}
}

// 设置 workbook，sheetName，以及 rowIndex，每次只需要再输入 colIndex，value
// rowIndex 从 1 开始； A1, B1
// colIndex 从 0 开始； 0 -> A, 1 -> B
func rowWriter(shWriter sheetWriterT, r int) func(c int, value interface{}) {
	axis := cellAxis(r)
	return func(c int, value interface{}) {
		shWriter(axis(c), value)
	}
}

// sheet 的读取的辅助函数
type sheetReaderT func(axis string) string

func sheetReader(workbook *excelize.File, sheetName string) sheetReaderT {
	return func(axis string) string {
		return workbook.GetCellValue(sheetName, axis)
	}
}

// 设置 workbook，sheetName，以及 rowIndex，每次只需要再输入 colIndex
// rowIndex 从 1 开始； A1, B1
// colIndex 从 0 开始； 0 -> A, 1 -> B
func rowReader(shReader sheetReaderT, r int) func(c int) string {
	axis := cellAxis(r)
	return func(c int) string {
		return shReader(axis(c))
	}
}
