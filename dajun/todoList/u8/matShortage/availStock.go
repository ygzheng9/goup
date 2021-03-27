package matShortage

import (
	"fmt"
	"sort"

	"github.com/360EntSecGroup-Skylar/excelize"
)

// 当前时刻，可用的库存量：现有库存量 + 到货未入库量 + 生产备料仓差异

// WarehouseInfo 仓库基本信息系统
type WarehouseInfo struct {
	WhCode string `db:"whCode" json:"whCode"`
	WhName string `db:"whName" json:"whName"`
}

// 从文件中读取所需的仓库
func getWarehouses(workbook *excelize.File) []WarehouseInfo {
	//  从文件中读取 仓库设置
	sheetName := "Sheet1"
	// 本 worksheet 写入的 helper
	readHelper := sheetReader(workbook, sheetName)

	results := []WarehouseInfo{}
	for r := 2; r < MaxLoop; r++ {
		// 从第二行开始
		rowData := rowReader(readHelper, r)

		// 第一列：仓库代码
		whCode := rowData(0)
		if len(whCode) == 0 {
			break
		}

		// 第二列：仓库名称，第三列：标记
		whName := rowData(1)
		flag := rowData(2)
		if flag == "Y" {
			//  代码是两位，不足两位，左边补零
			t := fmt.Sprintf("00%s", whCode)
			w := WarehouseInfo{
				WhCode: t[len(t)-2:],
				WhName: whName,
			}
			results = append(results, w)
		}
	}

	return results
}

// StockLevel 料号 + 数量
type StockLevel struct {
	InvCode string  `db:"invCode" json:"invCode"`
	Qty     float64 `db:"qty" json:"qty"`
}

// StockLevelSlice 一组物料的（料号，数量）
type StockLevelSlice []StockLevel

// 实现 sort.Interface
func (c StockLevelSlice) Len() int {
	return len(c)
}
func (c StockLevelSlice) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c StockLevelSlice) Less(i, j int) bool {
	return c[i].InvCode < c[j].InvCode
}

// getCurrentStock 取得指定仓库的现存量
func getCurrentStock(warehouses []WarehouseInfo) (StockLevelSlice, error) {
	// 根据 设置的仓库，读取 现存量

	// 拼接 in 条件
	cond := "''"
	for _, w := range warehouses {
		cond = fmt.Sprintf("%s, '%s'", cond, w.WhCode)
	}

	//  现存量查询
	sqlCmd := `
					select a.invCode, a.qty
					  from (
            select c.cInvCode invCode, SUM(iQuantity) qty
            from CurrentStock c
            where 1 = 1
            and c.cWhCode in (%s)
						group by  c.cInvCode) a
					where a.qty <> 0;
		`
	sqlCmd = fmt.Sprintf(sqlCmd, cond)
	// fmt.Println(sqlCmd)

	items := []StockLevel{}
	err := db.Select(&items, sqlCmd)

	return items, err
}

// 取得 到货未入库 的 采购量，委外量
func getOpenStock() (StockLevelSlice, error) {
	sqlCmd := `
			with dtl as (
			select h.dDate, h.cBusType, h.cCode, d.cInvCode, d.iQuantity, d.fValidInQuan, (d.iQuantity - d.fValidInQuan) avlQty
			from PU_ArrivalVouchs d
			inner join PU_ArrivalVouch h on h.ID = d.ID
			where 1 = 1
			and d.iQuantity <> d.fValidInQuan
			and h.dDate >= '2017-06-01'
			)
			select d.cInvCode invCode, SUM(avlQty) qty
			from dtl d
			group by d.cInvCode
			having SUM(avlQty) <> 0;
	`

	return getStockLevel(sqlCmd)
}

// 生产备料仓 和 生产订单子件需求 的差异
func getMiddleStock() (StockLevelSlice, error) {
	sqlCmd := `
        -- 生产订单子件需求 - 备料仓库存
    with po as (
        select a.invCode, sum(a.Qty - a.IssQty) qty
        from v_mom_moallocate a
        inner join v_mom_orderdetail_rpt d on d.ModID = a.ModID
        inner join v_mom_order_rpt h  on h.MoID = d.MoID
        inner join inventory i on i.cInvCode = d.invCode
            where 1 =  1
                -- and h.MoCode = '0000010537'
                and d.Status = '3'  -- 审核
        group by a.invCode
        ),
        inv as (
        select  a.cInvCode invCode, SUM(a.iQuantity) qty
            from currentstock a
        where a.cWhCode = '07'
            and a.iQuantity <> 0
        group by a.cInvCode
        ),
        cmb as (
        select a.invCode, (a.qty) qty
            from po a
        union all
        select b.invCode, (-1 * b.qty) qty
            from inv b
        ),
        cmb2 as (
        select invCode, SUM(qty) diff
            from cmb
        group by InvCode
            having SUM(qty) <> 0
        )
        select a.InvCode invCode, -1 * a.diff qty
            from cmb2 a
        inner join Inventory i on i.cInvCode = a.InvCode
    `

	return getStockLevel(sqlCmd)
}

// getStockLevel 辅助函数，执行 sql
func getStockLevel(sqlCmd string) (StockLevelSlice, error) {
	// fmt.Printf("getStockLevel: %s\n", sqlCmd)

	items := []StockLevel{}
	err := db.Select(&items, sqlCmd)

	return items, err
}

// saveStockLevel 保存库存量
func saveStockLevel(workbook *excelize.File, sheetName string, stocks StockLevelSlice) {
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
	setColumn(1, sheetName)

	offset := 2
	for idx, s := range stocks {
		if s.Qty != 0 {
			// 当前行的写入 helper
			setColumn := rowWriter(writeHelper, idx+offset)
			// 写入两列数据：第一列，第二列
			setColumn(0, s.InvCode)
			setColumn(1, s.Qty)
		}
	}
}

func mergeStocks(groups []StockLevelSlice) StockLevelSlice {
	// 先把各种库存的料号，放到一个 list 中
	allInvs := []string{}
	for _, g := range groups {
		for _, i := range g {
			allInvs = append(allInvs, i.InvCode)
		}
	}

	// 取得不重复料号
	distInvs := distString(allInvs)
	sort.Strings(distInvs)

	results := StockLevelSlice{}
	for _, i := range distInvs {
		// 每颗料
		s := StockLevel{
			InvCode: i,
			Qty:     0,
		}

		// 累加：现存量，未入库量，备料仓量
		for _, g := range groups {
			for _, k := range g {
				if i == k.InvCode {
					s.Qty += k.Qty
				}
			}
		}

		results = append(results, s)
	}

	return results
}

// 按照替代料，汇总
func replaceStock(inter InterchangeMatSlice, items StockLevelSlice) StockLevelSlice {
	results := StockLevelSlice{}

	for _, raw := range items {
		invCode := raw.InvCode
		// 查找替代料
		for _, a := range inter {
			if a.SrcMat == invCode {
				// 找到了替代料
				invCode = a.DestMat
				break
			}
		}

		isFound := false
		for idx, i := range results {
			if i.InvCode == invCode {
				isFound = true
				results[idx].Qty += raw.Qty
				break
			}
		}

		if !isFound {
			results = append(results, raw)
		}
	}

	return results
}
