package matShortage

import (
	"fmt"
	"sort"

	"github.com/360EntSecGroup-Skylar/excelize"
)

// calcInvDiff 计算子件总需求，和 当前可用库存的差异
func calcInvDiff(reqsByDay []StockLevelSlice, totalStocks StockLevelSlice) []StockLevelSlice {
	//  现有库存，和每天子件需求 的差异
	//  dailyItems: 每天的子件总需求 [[{invCode, qty}]]
	//  totalStocks: 当前可用库存 [{invCode, qty}]

	if len(reqsByDay) == 0 {
		return []StockLevelSlice{}
	}

	//  每天的子件清单都一样，所以取第一天的子件清单
	allInvs := []string{}
	for _, i := range reqsByDay[0] {
		allInvs = append(allInvs, i.InvCode)
	}
	// 排序
	sort.Strings(allInvs)
	sort.Sort(totalStocks)

	allDays := []StockLevelSlice{}
	for _, matList := range reqsByDay {
		sort.Sort(matList)

		//  每天, dailyItems 中是 每一天的需求
		oneDay := []StockLevel{}

		for _, i := range matList {
			// 每颗料的需求，负数
			oneMat := StockLevel{
				InvCode: i.InvCode,
				Qty:     -1 * i.Qty,
			}

			// 每颗料的可用库存，正数
			for _, s := range totalStocks {
				//  totalStocks 是可用库存
				if i.InvCode == s.InvCode {
					oneMat.Qty += s.Qty
					break
				}
			}
			oneDay = append(oneDay, oneMat)
		}
		allDays = append(allDays, oneDay)
	}
	return allDays
}

//  saveInvDiff 差异保存
func saveInvDiff(workbook *excelize.File, diffsByDay []StockLevelSlice, totalStocks StockLevelSlice,
	matInfos []MatInfo, headers [][]string) {
	//  差异保存
	//  items: 每天，每个子件，的库存差异 [[{invCode, qty}]]
	//  totalStocks； 总可用库存 [{invCode,qty}]
	//  matInfos: 物料基本细信息
	//  headers: 计划的表头，两行，20 列

	//  保存子件的需求（相对于母件的累计需求）
	//  计划天数
	dayCnt := len(diffsByDay)
	if dayCnt == 0 {
		fmt.Print("保存库存差异：无数据可保存")
		return
	}

	sheetName := "子件库存缺口"
	workbook.NewSheet(sheetName)
	// 本 worksheet 写入的 helper
	writeHelper := sheetWriter(workbook, sheetName)

	// 表头，第一行
	// 第一行写入的 helper
	setColumn := rowWriter(writeHelper, 1)

	columns := []string{"料号", "名称", "规格", "外购", "自制", "委外", "最小起订量", "提前期", "MCCode", "现有库存"}
	offsetCol := len(columns)
	for idx, name := range columns {
		// column 从 0 开始，也即： 0 --> A
		setColumn(idx, name)
	}

	// 日期表头，长度是计划的天数，每天里面两个值
	// 第二行的写入 helper
	setColumn2 := rowWriter(writeHelper, 2)
	for d, h := range headers {
		// 计算日期列的偏移量
		col := d + offsetCol

		// 每个日期，同一列，有两行
		// 第一行
		setColumn(col, h[0])

		// 第二行
		setColumn2(col, h[1])
	}

	//  取得物料号
	invList := []string{}
	for _, i := range diffsByDay[0] {
		invList = append(invList, i.InvCode)
	}
	// 排序
	sort.Strings(invList)

	for r, inv := range invList {
		//  每行一个料号
		// 第一行、第二行是表头，具体内容从第三行开始
		// 当前行的写入 helper
		setColumn = rowWriter(writeHelper, r+3)

		//  第一列 料号
		setColumn(0, inv)

		// 基本信息
		for _, mat := range matInfos {
			if mat.InvCode == inv {
				setColumn(1, mat.InvName)
				setColumn(2, mat.InvStd)
				setColumn(3, mat.Purchase)
				setColumn(4, mat.SelfMade)
				setColumn(5, mat.Outsourcing)
				setColumn(6, mat.Moq)
				setColumn(7, mat.Leadtime)
				setColumn(8, mat.McCode)
				break
			}
		}

		// 总可用库存
		for _, s := range totalStocks {
			if s.InvCode == inv {
				setColumn(9, fmt.Sprintf("%.2f", s.Qty))
				break
			}
		}

		//  后面每一列，是库存缺口
		// 延续之前的列序号
		for day, matList := range diffsByDay {
			for _, mat := range matList {
				if inv == mat.InvCode {
					setColumn(day+offsetCol, fmt.Sprintf("%.2f", mat.Qty))
					break
				}
			}
		}
	}
}
