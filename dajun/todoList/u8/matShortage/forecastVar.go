package matShortage

import (
	"fmt"
	"math"
	"os"
	"sort"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
)

// MatQty 物料数量
type MatQty struct {
	InvCode  string  `db:"invCode"`
	InvName  string  `db:"invName"`
	Quantity float64 `db:"qty"`
}

// MatQtySlice 一组物料清单
type MatQtySlice []MatQty

// 实现 sort.Interface
func (c MatQtySlice) Len() int {
	return len(c)
}
func (c MatQtySlice) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c MatQtySlice) Less(i, j int) bool {
	return c[i].InvCode < c[j].InvCode
}

// 执行 sql，加载物料的数量
func loadMatQty(sqlCmd string) (MatQtySlice, error) {
	items := MatQtySlice{}
	err := db.Select(&items, sqlCmd)

	return items, err
}

// 根据单据日期，加载预测订单行项目
func loadForecast(date string) (MatQtySlice, error) {
	sqlCmd := `
	select isnull(p.InvCode,'') invCode, isnull(i.cInvName, '') invName, isnull(d.FQty, 0) qty
  from mps_forecast h
 inner join mps_forecastdetail d on d.ForecastId = h.ForecastId
 inner join bas_part p on p.PartId = d.PartId
 inner join Inventory i on p.InvCode = i.cInvCode
 where h.CreateDate = '%s'
	`

	sqlCmd = fmt.Sprintf(sqlCmd, date)
	return loadMatQty(sqlCmd)
}

// 加载一个期间内的成品入库
func loadStockIn(from, to string) (MatQtySlice, error) {
	sqlCmd := `
			select d.cInvCode invCode, i.cInvName invName, sum(d.iQuantity) qty
		from RdRecord10 h
		inner join RdRecords10 d on h.ID = d.ID
		inner join inventory i on d.cInvCode = i.cInvCode
		where 1 = 1
		and (h.dDate > '%s' and h.dDate < '%s')
		group by d.cInvCode, i.cInvName
	`

	sqlCmd = fmt.Sprintf(sqlCmd, from, to)
	return loadMatQty(sqlCmd)
}

// 计算一段期间内的差异
func CalcForecastDiff(from, to string, outFile string) {
	results := MatQtySlice{}
	fmt.Println("获取预测订单...")
	items1, err := loadForecast(from)
	if err != nil {
		fmt.Println("获取预测订单失败")
		fmt.Printf("err: %+v\n", err)
		return
	}

	items2, err := loadForecast(to)
	if err != nil {
		fmt.Println("获取预测订单失败")
		fmt.Printf("err: %+v\n", err)
		return
	}

	fmt.Println("获取成品入库...")
	items3, err := loadStockIn(from, to)
	if err != nil {
		fmt.Println("获取成品入库量失败")
		fmt.Printf("err: %+v\n", err)
		return
	}

	// 三个结果集，按顺序排好
	fmt.Println("计算白坯差异...")
	inputs := []MatQtySlice{items1, items2, items3}
	for s, list := range inputs {
		for _, i := range list {
			qty := i.Quantity

			if s > 0 {
				// + from - to - stockin
				qty = -1 * qty
			}

			isFound := false
			// 从料号名称，提取白坯名称：规则：第一个 -，之后的 两位
			pos := strings.Index(i.InvName, "-")
			if pos > 0 {
				// 有母件
				base := i.InvName[:pos+3]

				for idx, a := range results {
					if base == a.InvName {
						// 相同白坯名，汇总数量
						isFound = true
						results[idx].Quantity += qty
					}
				}

				if !isFound {
					// 白坯第一次出现
					n := MatQty{
						InvCode:  base,
						InvName:  base,
						Quantity: qty,
					}
					results = append(results, n)
				}
			} else {
				// 无母件，仅按料号
				for idx, a := range results {
					if i.InvCode == a.InvCode {
						// 相同料号，汇总数量
						isFound = true
						results[idx].Quantity += qty
					}
				}

				if !isFound {
					// 料号第一次出现
					n := MatQty{
						InvCode:  i.InvCode,
						InvName:  i.InvName,
						Quantity: qty,
					}
					results = append(results, n)
				}
			}
		}
	}

	// 排除掉数量为零的物料
	nonZero := MatQtySlice{}
	for _, i := range results {
		if math.Abs(i.Quantity) >= 0.005 {
			nonZero = append(nonZero, i)
		}
	}

	// 数据保存到文件
	fmt.Println("保存结果...")
	xlsx := excelize.NewFile()
	saveForecastVar(xlsx, from, items1)
	saveForecastVar(xlsx, to, items2)
	saveForecastVar(xlsx, "入库", items3)
	saveForecastVar(xlsx, "白坯差异", nonZero)

	// 删除默认的 tab
	xlsx.DeleteSheet("Sheet1")

	// 保存前，删除上次文件
	os.Remove(outFile)
	err = xlsx.SaveAs(outFile)
	if err != nil {
		fmt.Printf("保存文件失败：%s.  \n%+v\n", outFile, err)
		return
	}

	fmt.Println("完成")

	return
}

// 根据文件，获取 from/to
func GetDuration(inFile string) (string, string, error) {
	xlsx, err := excelize.OpenFile(inFile)
	if err != nil {
		fmt.Println("输入文件错误。")
		fmt.Printf("%+v\n", err)
		return "", "", err
	}

	//  从文件中读取 仓库设置
	sheetName := "Sheet1"
	// 本 worksheet 写入的 helper
	readHelper := sheetReader(xlsx, sheetName)
	// 第一行
	rowData := rowReader(readHelper, 1)
	// 第一行，第二列
	from := rowData(1)

	// 第二行
	rowData = rowReader(readHelper, 2)
	// 第二行，第二列
	to := rowData(1)

	return from, to, nil
}

// saveForecastVar 保存原始数据，和差异数据
func saveForecastVar(workbook *excelize.File, sheetName string, stocks MatQtySlice) {
	// 排序
	sort.Sort(stocks)

	//  保存到文件：现存量
	workbook.NewSheet(sheetName)
	// 本 worksheet 写入的 helper
	writeHelper := sheetWriter(workbook, sheetName)

	// 第一行写入的 helper
	setColumn := rowWriter(writeHelper, 1)

	// 第一行，表头
	setColumn(0, "料号")
	setColumn(1, "名称")
	setColumn(2, "数量")

	// 从第二行开始
	offset := 2
	for idx, s := range stocks {
		if s.Quantity != 0 {
			// 当前行的写入 helper
			setColumn := rowWriter(writeHelper, idx+offset)
			// 写入两列数据：第一列，第二列
			setColumn(0, s.InvCode)
			setColumn(1, s.InvName)
			setColumn(2, s.Quantity)

		}
	}
}
