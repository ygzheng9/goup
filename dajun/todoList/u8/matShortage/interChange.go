package matShortage

import (
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
)

// InterchangeMat 替代料的
type InterchangeMat struct {
	SrcMat  string
	DestMat string
}

// InterchangeMatSlice 一组替代料
type InterchangeMatSlice []InterchangeMat

func loadInterchangeMat(workbook *excelize.File) InterchangeMatSlice {
	//  tab 名字固定
	sheetName := "替代物料清单"
	// 本 worksheet 写入的 helper
	readHelper := sheetReader(workbook, sheetName)

	results := InterchangeMatSlice{}
	for r := 2; r < MaxLoop; r++ {
		// 从第二行开始
		rowData := rowReader(readHelper, r)

		// 第一列：替代前物料
		srcMat := strings.TrimSpace(rowData(0))
		if len(srcMat) == 0 {
			break
		}

		// 第四列：替代后物料
		destMat := strings.TrimSpace(rowData(3))

		w := InterchangeMat{
			SrcMat:  srcMat,
			DestMat: destMat,
		}
		results = append(results, w)
	}

	return results
}
