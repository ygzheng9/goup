package matShortage

import (
	"fmt"
	"os"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

// SetupDB 把数据库连接传入
func SetupDB(dbConn *sqlx.DB) {
	db = dbConn
}

// GetDB 返回 U8 数据库连接
func GetDB() *sqlx.DB {
	return db
}

// StartProcess 主程序
func StartProcess(inFile, outFile, interFile string) {
	// inFile := "E:/99.localDev/easypy/u8/1.xlsx"
	xlsx, err := excelize.OpenFile(inFile)
	if err != nil {
		fmt.Println("输入文件错误。")
		fmt.Printf("%+v\n", err)
		return
	}

	xlsxInter, err := excelize.OpenFile(interFile)
	if err != nil {
		fmt.Println("替代料文件错误。")
		fmt.Printf("%+v\n", err)
		return
	}

	fmt.Println("计算中....")

	interMats := loadInterchangeMat(xlsxInter)

	whs := getWarehouses(xlsx)
	// fmt.Printf("warehouse: %d\n", len(whs))

	currStocks, err := getCurrentStock(whs)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}
	// fmt.Printf("stocks: %d\n", len(currStocks))

	openStocks, err := getOpenStock()
	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}
	// fmt.Printf("openStocks: %d\n", len(openStocks))

	midStocks, err := getMiddleStock()
	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}
	// fmt.Printf("midStocks: %d\n", len(midStocks))

	rawStocks := mergeStocks([]StockLevelSlice{currStocks, openStocks, midStocks})
	// fmt.Printf("allStocks: %d\n", len(allStocks))

	// 替代料替换
	allStocks := replaceStock(interMats, rawStocks)

	headers := getHeader(xlsx)
	// fmt.Printf("headers: %d\n", len(headers))

	accReqs := calcAccReq(xlsx)
	// fmt.Printf("planItems: %d\n", len(accReqs))

	matInfos, err := getInvInfo()
	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}
	// fmt.Printf("matInfos: %d\n", len(matInfos))

	rawBoms, err := loadBOM()
	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}
	// fmt.Printf("boms: %d\n", len(boms))

	boms := replacePart(interMats, rawBoms)

	planExplodes := expandOneLevel(accReqs, boms)
	// fmt.Printf("planExplodes: %d\n", len(planExplodes))

	byMats := groupByMat(planExplodes)
	// fmt.Printf("byMats: %d\n", len(byMats))

	diffs := calcInvDiff(byMats, allStocks)
	// fmt.Printf("diffs: %d\n", len(diffs))

	// for _, v := range matInfos {
	// 	fmt.Printf("%s, %s\n", v.InvCode, v.InvName)
	// }

	// 准备输出文件
	xlsxOut := excelize.NewFile()

	saveStockLevel(xlsxOut, "现存量", currStocks)
	saveStockLevel(xlsxOut, "未入库量", openStocks)
	saveStockLevel(xlsxOut, "备料仓差异量", midStocks)
	saveStockLevel(xlsxOut, "总可用量(原始)", rawStocks)
	saveStockLevel(xlsxOut, "总可用量(替代)", allStocks)

	saveAccReq(xlsxOut, accReqs)
	// fmt.Println("saved: accReqs")

	saveBOMData(xlsxOut, accReqs, matInfos)
	// fmt.Println("saved: boms")

	saveMatReq(xlsxOut, byMats)
	// fmt.Println("saved: byMats")

	saveInvDiff(xlsxOut, diffs, allStocks, matInfos, headers)
	// fmt.Println("saved: diff")

	// 删除默认的 tab
	xlsxOut.DeleteSheet("Sheet1")

	fmt.Println("保存中....")

	// 输出文件
	// outFile := "E:/99.localDev/easypy/u8/1_out.xlsx"

	// 保存前，删除上次文件
	os.Remove(outFile)
	err = xlsxOut.SaveAs(outFile)
	if err != nil {
		fmt.Printf("保存文件失败：%s.  \n%+v\n", outFile, err)
		return
	}

	fmt.Println("完成")
}
