package rawInput

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

// AmtBase 金额的转换基数
const AmtBase = 3.14

// toNumber 把字符串转成 float
func toNumber(str string) float64 {
	number, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0.0
	}
	return number
}

// POHead 采购订单头信息
type POHead struct {
	POID       string `json:"poID"`
	POCode     string `json:"poCode"`
	PODate     string `json:"poDate"`
	VendorCode string `json:"vendorCode"`
	PersonCode string `json:"personCode"`
}

// ReadPOHead 读取采购订单头
func ReadPOHead(fileName string) ([]POHead, error) {
	items := []POHead{}

	// 打开文件
	f, err := os.Open(fileName)
	if err != nil {
		return items, err
	}
	// 使用 buf
	buf := bufio.NewReader(f)

	log.Println("begin reading....")

	// 第一行是标题，只读取一次
	isHead := true
	for {
		// 逐行读取
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		if err != nil {
			if err == io.EOF {
				// 读取完毕后的结果
				return items, nil
			}
			return items, err
		}

		// 第一行是标题，跳过
		if isHead {
			isHead = false
			continue
		}

		// 每个栏位是以 tab 分割；如果是 逗号，直接使用 csv package；
		record := strings.Split(line, "\t")

		// 每一行的列数
		if len(record) != 5 {
			continue
		}

		// 日期只保留到天，格式为：2012-01-23
		poDate := record[2]
		if len(poDate) > 10 {
			poDate = poDate[0:10]
		}

		item := POHead{
			POID:       record[0],
			POCode:     record[1],
			PODate:     poDate,
			VendorCode: record[3],
			PersonCode: record[4],
		}

		items = append(items, item)
	}
}

// POItem 采购订单头信息
type POItem struct {
	ID          string  `json:"itemID"`
	POID        string  `json:"poID"`
	InvCode     string  `json:"invCode"`
	Quantity    float64 `json:"quantity"`
	UnitPrice   float64 `json:"unitPrice"`
	NetAmt      float64 `json:"netAmt"`
	TaxAmt      float64 `json:"taxAmt"`
	TotalAmt    float64 `json:"totalAmt"`
	ProjectCode string  `json:"projCode"`
	ProjectName string  `json:"projName"`
	POCode      string  `json:"poCode"`
	PODate      string  `json:"poDate"`
	VendorCode  string  `json:"vendorCode"`
	PersonCode  string  `json:"personCode"`
}

// ReadPOItem 读取采购订单行项目
func ReadPOItem(fileName string) ([]POItem, error) {
	items := []POItem{}

	// 打开文件
	f, err := os.Open(fileName)
	if err != nil {
		return items, err
	}
	// 使用 buf
	buf := bufio.NewReader(f)

	log.Println("begin reading....")

	// 第一行是标题，只读取一次
	isHead := true
	for {
		// 逐行读取
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		if err != nil {
			if err == io.EOF {
				// 读取完毕后的结果
				return items, nil
			}
			return items, err
		}

		// 第一行是标题，跳过
		if isHead {
			isHead = false
			continue
		}

		// 每个栏位是以 tab 分割；如果是 逗号，直接使用 csv package；
		record := strings.Split(line, "\t")

		// 每一行的列数
		if len(record) != 10 {
			continue
		}

		item := POItem{
			ID:          record[0],
			POID:        record[1],
			InvCode:     record[2],
			Quantity:    toNumber(record[3]),
			UnitPrice:   toNumber(record[4]),
			NetAmt:      toNumber(record[5]),
			TaxAmt:      toNumber(record[6]),
			TotalAmt:    toNumber(record[7]),
			ProjectCode: record[8],
			ProjectName: record[9],
		}

		items = append(items, item)
	}
}

// LoadPOItems 文件中已经是 头+行
func LoadPOItems(fileName string) ([]POItem, error) {
	items := []POItem{}

	// 打开文件
	f, err := os.Open(fileName)
	if err != nil {
		return items, err
	}
	// 使用 buf
	buf := bufio.NewReader(f)

	log.Println("begin reading....")

	// 第一行是标题，只读取一次
	isHead := true
	for {
		// 逐行读取
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		if err != nil {
			if err == io.EOF {
				// 读取完毕后的结果
				return items, nil
			}
			return items, err
		}

		// 第一行是标题，跳过
		if isHead {
			isHead = false
			continue
		}

		// 每个栏位是以 tab 分割；如果是 逗号，直接使用 csv package；
		record := strings.Split(line, "\t")

		// 每一行的列数
		if len(record) != 14 {
			continue
		}

		// 日期只保留到天，格式为：2012-01-23
		poDate := record[11]
		if len(poDate) > 10 {
			poDate = poDate[0:10]
		}

		item := POItem{
			ID:          record[0],
			POID:        record[1],
			InvCode:     record[2],
			Quantity:    toNumber(record[3]),
			UnitPrice:   toNumber(record[4]),
			NetAmt:      toNumber(record[5]),
			TaxAmt:      toNumber(record[6]),
			TotalAmt:    toNumber(record[7]),
			ProjectCode: record[8],
			ProjectName: record[9],
			POCode:      record[10],
			PODate:      poDate,
			VendorCode:  record[12],
			PersonCode:  record[13],
		}

		items = append(items, item)
	}
}

// POItemsByDate 根据订单日期，过滤行项目
func POItemsByDate(items []POItem, start, end string) []POItem {
	// 根据传入的参数过滤
	var results []POItem
	for _, v := range items {
		// fmt.Println(v.PODate)
		if v.PODate >= start && v.PODate <= end {
			results = append(results, v)
		}
	}

	fmt.Printf("results len: %d, start: %s, end: %s\n", len(results), start, end)

	// 返回 nil 比较麻烦，所以返回一个长度为零array
	if len(results) == 0 {
		results = []POItem{}
	}

	return results
}

// MatByMonth 按月度汇总的物料信息
type MatByMonth struct {
	BizMonth string  `json:"bizMonth"`
	InvCode  string  `json:"invCode"`
	Qty      float64 `json:"qty"`
	Amt      float64 `json:"amt"`
}

// LoadMatByMonth 加载按月度汇总的物料信息
func LoadMatByMonth(fileName string) ([]MatByMonth, error) {
	items := []MatByMonth{}

	// 打开文件
	f, err := os.Open(fileName)
	if err != nil {
		return items, err
	}
	// 使用 buf
	buf := bufio.NewReader(f)

	log.Println("begin reading....")

	// 第一行是标题，只读取一次
	isHead := true
	for {
		// 逐行读取
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		if err != nil {
			if err == io.EOF {
				// 读取完毕后的结果
				return items, nil
			}
			return items, err
		}

		// 第一行是标题，跳过
		if isHead {
			isHead = false
			continue
		}

		// 每个栏位是以 tab 分割；如果是 逗号，直接使用 csv package；
		record := strings.Split(line, "\t")

		// 每一行的列数
		if len(record) != 4 {
			continue
		}

		item := MatByMonth{
			BizMonth: record[0],
			InvCode:  record[1],
			Qty:      toNumber(record[2]),
			Amt:      toNumber(record[3]),
		}

		items = append(items, item)
	}
}

// BOMComponent 单层BOM，母件及子件
type BOMComponent struct {
	ChildInv   string `json:"childInv"`
	ChildName  string `json:"childName"`
	ParentInv  string `json:"parentInv"`
	ParentName string `json:"parentName"`
}

// LoadBOMComponent 加载单层 BOM
func LoadBOMComponent(fileName string) ([]BOMComponent, error) {
	items := []BOMComponent{}

	// 打开文件
	f, err := os.Open(fileName)
	if err != nil {
		return items, err
	}
	// 使用 buf
	buf := bufio.NewReader(f)

	log.Println("begin reading....")
	// 第一行是标题，需要跳过
	isHead := true
	for {
		// 逐行读取
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		if err != nil {
			if err == io.EOF {
				// 读取完毕后的结果
				return items, nil
			}
			return items, err
		}

		// 第一行是标题，跳过
		if isHead {
			isHead = false
			continue
		}

		// 每个栏位是以 tab 分割；如果是 逗号，直接使用 csv package；
		record := strings.Split(line, "\t")

		// 数据文件中只有 4 列，
		if len(record) != 4 {
			continue
		}

		item := BOMComponent{
			ChildInv:   record[0],
			ChildName:  record[1],
			ParentInv:  record[2],
			ParentName: record[3],
		}

		items = append(items, item)
	}
}

// MatInfo 物料基本信息
type MatInfo struct {
	InvCode    string `json:"invCode"`
	InvName    string `json:"invName"`
	InvStd     string `json:"invStd"`
	IsPurchase int64  `json:"isPurchase"`
	IsSelfMade int64  `json:"isSelfMade"`
	IsProxy    int64  `json:"isProxy"`
	MoQ        int64  `json:"moQ"`
	LeadTime   int64  `json:"leadTime"`
	FileName   string `json:"fileName"`
	Version    string `json:"version"`
}

// LoadMatInfo 加载物料基本信息
func LoadMatInfo(fileName string) ([]MatInfo, error) {
	items := []MatInfo{}

	// 打开文件
	f, err := os.Open(fileName)
	if err != nil {
		return items, err
	}
	// 使用 buf
	buf := bufio.NewReader(f)

	// 字符串转成数字
	parseInt := func(str string) int64 {
		if str == "NULL" {
			return 0
		}

		number, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return 0
		}
		return number
	}

	// 过滤掉某些字符
	normalizeStr := func(str string) string {
		if str == "NULL" {
			return ""
		}
		return str
	}

	log.Println("begin reading....")
	// 第一行是标题，需要跳过
	isHead := true
	for {
		// 逐行读取
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		if err != nil {
			if err == io.EOF {
				// 读取完毕后的结果
				return items, nil
			}
			return items, err
		}

		// 第一行是标题，跳过
		if isHead {
			isHead = false
			continue
		}

		// 每个栏位是以 tab 分割；如果是 逗号，直接使用 csv package；
		record := strings.Split(line, "\t")

		// 数据文件中只有 10 列，
		if len(record) != 10 {
			continue
		}

		item := MatInfo{
			InvCode:    normalizeStr(record[0]),
			InvName:    normalizeStr(record[1]),
			InvStd:     normalizeStr(record[2]),
			IsPurchase: parseInt(record[3]),
			IsSelfMade: parseInt(record[4]),
			IsProxy:    parseInt(record[5]),
			MoQ:        parseInt(record[6]),
			LeadTime:   parseInt(record[7]),
			FileName:   normalizeStr(record[8]),
			Version:    normalizeStr(record[9]),
		}

		items = append(items, item)
	}
}
