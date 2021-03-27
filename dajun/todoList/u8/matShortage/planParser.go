package matShortage

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize"
)

// 根据输入的每天的需求，计算天累计需求

// StartColumn 日计划开始的列号，从 0 开始
const StartColumn = 3

// TotalColumn 日计划的数据列数
const TotalColumn = 25

// MaxLoop 最大循环次数
const MaxLoop = 5000

// PlanItem 输入的计划数据
type PlanItem struct {
	// 料号
	invCode string
	// 每天的需求
	currQtys []float64

	// 累计到今天的需求
	accQtys []float64

	// 料号的第一层子件（展开到外购，委外，自制领用）
	subs []OneLevel

	// 子件的累计需求
	accSubsReq [][]StockLevel
}

// PlanItemSlice 一组物料的计划信息
type PlanItemSlice []PlanItem

// 实现 sort.Interface
func (c PlanItemSlice) Len() int {
	return len(c)
}
func (c PlanItemSlice) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c PlanItemSlice) Less(i, j int) bool {
	return c[i].invCode < c[j].invCode
}

// 从文件中读取 每天的计划，折算成累计的计划
func calcAccReq(workbook *excelize.File) []PlanItem {
	// 从文件中读取 每天的计划，折算成累计的计划
	sheetName := "Sheet2"
	// 本 worksheet 写入的 helper
	readHelper := sheetReader(workbook, sheetName)

	planList := PlanItemSlice{}
	// 第一列 A：料号；从第四列 D 开始，当天计划信息
	// 计划数据从 第三行 开始
	for row := 3; row < MaxLoop; row++ {
		// 本行的读取 helper
		getColumn := rowReader(readHelper, row)

		// 第一列为空，则跳出循环
		invCode := getColumn(0)
		if len(invCode) == 0 {
			break
		}

		oneMat := PlanItem{}
		// 当前料号
		oneMat.invCode = invCode

		// 计划数据：
		for i := 0; i < TotalColumn; i++ {
			col := StartColumn + i
			qty, err := strconv.ParseFloat(getColumn(col), 64)
			if err != nil {
				qty = 0
			}

			oneMat.currQtys = append(oneMat.currQtys, qty)
			oneMat.accQtys = append(oneMat.accQtys, qty)
		}

		// 计算累计量
		cnt := len(oneMat.currQtys)
		for i := 1; i < cnt; i++ {
			oneMat.accQtys[i] = oneMat.currQtys[i] + oneMat.accQtys[i-1]
		}
		planList = append(planList, oneMat)
	}

	// 排序
	sort.Sort(planList)
	return planList
}

// # 把母件的累计需求，保存到文件
func saveAccReq(workbook *excelize.File, planItems []PlanItem) {
	if len(planItems) == 0 {
		fmt.Println("累计需求：无数据传入")
		return
	}

	//  计划天数
	dayCnt := len(planItems[0].accQtys)
	if dayCnt == 0 {
		fmt.Println("累计需求：无累计需求")
		return
	}

	sheetName := "母件累计需求"
	workbook.NewSheet(sheetName)
	// 本 worksheet 写入的 helper
	writeHelper := sheetWriter(workbook, sheetName)

	// 第一行写入的 helper
	setColumn := rowWriter(writeHelper, 1)

	// 第一行
	for c := 0; c <= dayCnt; c++ {
		if c == 0 {
			// 第一行、第一列
			setColumn(c, "料号")
		} else {
			// 从第二列开始，内容是：D1，D2，D3,...
			setColumn(c, fmt.Sprintf("D%d", c))
		}
	}

	for idx, p := range planItems {
		// 从第二行开始（第一行是表头）
		// 本行写入 helper
		setColumn = rowWriter(writeHelper, idx+2)

		//  第一列 料号
		setColumn(0, p.invCode)

		// 后面每一列，是累计需求量
		for c := 0; c < dayCnt; c++ {
			// 第一列是料号，从第二列开始
			setColumn(c+1, fmt.Sprintf("%.02f", p.accQtys[c]))
		}
	}
}
