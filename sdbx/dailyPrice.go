package sdbx

import (
	"encoding/csv"
	"io"
	"log"
	"os"
)

// 根据名称，返回每一天，指定时点的数据
// 格式：日期，时间，价格
func getDailyPrice(idx string) [][]string {
	// 读取文件
	// 对行进行过滤，保留下匹配的
	// 输出结果：日期，价格

	var err error

	// 设置数据源，过滤条件
	fileName := "c:/tmp/rawData/" + idx + ".csv"
	filter := "10:15"

	// 打开文件
	reader, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}

	// 返回的结果
	var results [][]string

	// 读取每一行数据
	r := csv.NewReader(reader)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		if len(record) < 2 {
			log.Fatal("格式错误！")
		}

		// 过滤每一行
		if record[1] == filter {
			// fmt.Printf("%s %s %s\n", record[0], record[1], record[2])
			results = append(results, []string{record[0], record[1], record[2]})
		}
	}
	return results
}

// 把多个数据，合并到同一个表格中
// 参数格式：多支的价格表（每支）；每个价格表有多行（每天）；每行：日期，时间，价格
func mergePrice(priceList [][][]string) [][]string {
	// 如果只传入一枝
	if len(priceList) == 1 {
		return priceList[0]
	}

	first := priceList[0]
	others := priceList[1:]

	log.Printf("first: %d, others: %d\n", len(first), len(others))

	// 以第一个数据为基础，取后面数据中，对应键值的数据；
	// 键值：起始两列
	var results [][]string
	for _, v := range first {
		// 第一个数据集中的每一行
		tmp := v

		// 后续每一个结果集
		for _, o := range others {
			// 结果集中，找到键值相等的
			found := false
			for _, i := range o {
				if i[0] == v[0] && i[1] == v[1] {
					tmp = append(tmp, i[2])
					found = true
					break
				}
			}
			if !found {
				// 如果某个后续结果集中，没有键值相等的，那么插入一个特殊值，否则会错列
				tmp = append(tmp, "NA")
			}
		}

		// 合并后的行
		results = append(results, tmp)
	}

	return results
}

// 多个输入，按列，返回记录
// 参数：列表；
// 返回：每天，时间，按列每支的价格
func getPriceTable(idxList []string) [][]string {
	// 生成表头
	header := []string{"日期", "时间"}
	for _, i := range idxList {
		header = append(header, i)
	}

	// 一个列表，每个值是一支单独的价格表，价格表是每天的数据
	var priceList [][][]string
	for _, idx := range idxList {
		items := getDailyPrice(idx)
		priceList = append(priceList, items)
	}

	// 把每支单独的价格，按列合并成一张表
	results := mergePrice(priceList)

	// 插入表头
	results = append([][]string{header}, results...)

	return results
}
