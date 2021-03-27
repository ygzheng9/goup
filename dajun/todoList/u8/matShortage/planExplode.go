package matShortage

import (
	"fmt"
	"sort"

	"github.com/360EntSecGroup-Skylar/excelize"
)

// 按照 BOM 结构，把 母件的累计需求，展开到子件需求；
// 再按照子件汇总
// 保存：汇总后的子件需求
// 保存：母件单层 BOM 展开后的结果

// expandOneLevel 计算每个母件，每天，累计需求，展开后的子件需求
func expandOneLevel(parentReqs []PlanItem, boms []BOM) []PlanItem {
	// 循环1：每天的计划
	// 循环2：一天计划中，循环每颗料；展开BOM，计算子件需求；

	planItems := parentReqs

	// 取得母件的下级子件
	for i := range planItems {
		invCode := planItems[i].invCode
		planItems[i].subs = findAllSubs(invCode, boms)
	}

	// 根据下级子件，计算每天的子件需求量（累计量）
	// 对每一个母件
	for i := range planItems {
		// 后面需要修改 当前元素，所以使用下标，
		// 并且使用引用，方便书写
		item := &planItems[i]

		// 清空子件需求
		item.accSubsReq = [][]StockLevel{}

		// 对母件每一天的累计需求量
		for _, accReq := range item.accQtys {
			// 该母件的每个子件的耗量 * 母件的累计需求
			oneDay := []StockLevel{}
			for _, s := range item.subs {
				a := StockLevel{
					InvCode: s.InvCode,
					Qty:     s.BaseQty * accReq,
				}
				oneDay = append(oneDay, a)
			}
			// 每天都有子件需求
			item.accSubsReq = append(item.accSubsReq, oneDay)
		}
	}
	return planItems
}

// groupByMat 所有母件的子件需求，按照子件汇总
func groupByMat(planItems []PlanItem) []StockLevelSlice {
	// 参数：按照日累计需求，展开单层BOM的需求量
	// 按照子件汇总（合并母件），按天显示子件总需求

	// 首先找到所有母件的所有子件
	allInvs := []string{}
	for _, p := range planItems {
		// 每个母件的 子件料号
		for _, s := range p.subs {
			allInvs = append(allInvs, s.InvCode)
		}
	}
	// 所有母件的子件料号
	distInvs := distString(allInvs)
	// 排序
	sort.Strings(distInvs)

	// 按天，汇总这些子件的需求
	dayCnt := len(planItems[0].accQtys)
	// 每一天的子件汇总需求
	allDays := []StockLevelSlice{}
	for d := 0; d < dayCnt; d++ {
		// 对于每一天
		oneDay := StockLevelSlice{}
		for _, inv := range distInvs {
			// 这一天，所有子件的需求
			oneMat := StockLevel{
				InvCode: inv,
				Qty:     0,
			}
			// 累计所有母件中，该子件的需求
			for _, p := range planItems {
				for _, mat := range p.accSubsReq[d] {
					if inv == mat.InvCode {
						oneMat.Qty += mat.Qty
						break
					}
				}
			}
			oneDay = append(oneDay, oneMat)
		}

		// 每天的子件汇总需求，合并到总需求中
		allDays = append(allDays, oneDay)
	}

	return allDays
}

// saveMatReq 保存子件需求
func saveMatReq(workbook *excelize.File, matReqByDay []StockLevelSlice) {
	// 保存子件的需求（相对于母件的累计需求）
	// 计划天数
	dayCnt := len(matReqByDay)
	if dayCnt == 0 {
		fmt.Println("保存子件需求：无子件可保存")
		return
	}

	sheetName := "子件总需求"
	workbook.NewSheet(sheetName)
	// 本 worksheet 写入的 helper
	writeHelper := sheetWriter(workbook, sheetName)

	// 第一行的写入 helper
	setColumn := rowWriter(writeHelper, 1)

	// 第一行，第一列
	setColumn(0, "料号")
	for d := 0; d < dayCnt; d++ {
		// 第一行，从第二列开始；列从 0 开始，也即：0->A, 1->B, 2->C
		setColumn(d+1, fmt.Sprintf("D%d", d+1))
	}

	// 每一天的料号清单是一样的，所以取第一天的即可
	invList := []string{}
	for _, i := range matReqByDay[0] {
		invList = append(invList, i.InvCode)
	}
	// 排序
	sort.Strings(invList)

	for r, inv := range invList {
		// 第一行是表头，从第二行开始，是正文的料号
		rowIdx := r + 2

		// 本行写入的 Helper
		setColumn = rowWriter(writeHelper, rowIdx)

		// 第一列 料号
		setColumn(0, inv)

		// 后面每一列，是这天的累计需求量
		for d, matList := range matReqByDay {
			for _, mat := range matList {
				if inv == mat.InvCode {
					// 从第二列开始
					setColumn(d+1, fmt.Sprintf("%.2f", mat.Qty))
					break
				}
			}
		}
	}
}

// saveBOMData 保存 母件的 BOM一阶展开
func saveBOMData(workbook *excelize.File, planItems []PlanItem, matInfos []MatInfo) {
	//  保存 bom 信息
	sheetName := "BOM"
	workbook.NewSheet(sheetName)

	// 本 worksheet 写入的 helper
	writeHelper := sheetWriter(workbook, sheetName)

	//  写入表头
	headers := []string{"子件料号", "子件名称", "用量", "规格", "外购", "自制", "委外", "最小起订量", "提前期", "MCCode"}
	// axis := ""
	// 第一行写入的 helper
	setColumn := rowWriter(writeHelper, 1)
	for i, h := range headers {
		// 第一列是母件代码，子件代码从第二列开始，
		setColumn(i+1, h)
	}

	// 当前行号，第一行已经是表头，正文从第二行开始
	rowIdx := 2
	for _, p := range planItems {
		// 本行的写入 helper
		setColumn = rowWriter(writeHelper, rowIdx)

		//  每一颗母件
		invCode := p.invCode

		// 母件 第一列：code
		setColumn(0, invCode)

		for _, mat := range matInfos {
			if mat.InvCode == invCode {
				// 第二列
				setColumn(1, mat.InvName)

				// 第三列
				setColumn(2, mat.InvStd)

				break
			}
		}

		//  对该母件的所有子件
		for _, part := range p.subs {
			//  子件的基础数据
			for _, mat := range matInfos {
				if mat.InvCode == part.InvCode {
					//  母件下一行
					rowIdx++

					// 本行的写入 helper
					setColumn = rowWriter(writeHelper, rowIdx)

					//  从 第二列 开始
					setColumn(1, mat.InvCode)
					setColumn(2, mat.InvName)
					setColumn(3, part.BaseQty)
					setColumn(4, mat.InvStd)
					setColumn(5, mat.Purchase)
					setColumn(6, mat.SelfMade)
					setColumn(7, mat.Outsourcing)
					setColumn(8, mat.Moq)
					setColumn(9, mat.Leadtime)
					setColumn(10, mat.McCode)
					break
				}
			}
		}

		// 下一个母件
		rowIdx++
	}
}
