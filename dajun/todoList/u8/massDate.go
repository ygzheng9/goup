package u8

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	// "github.com/xuri/excelize"
	"github.com/360EntSecGroup-Skylar/excelize"
)

// 上载的数据
type massDate struct {
	InvCode  string
	DayCount int
}

// MassDateUpload 上载销售发货单的实际批次号
func MassDateUpload(c *gin.Context) {
	// 获取上载的文件
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("FormFile error: %+v\n", err),
		})
		return
	}

	// 解析上载的文件
	items, err := parseBatchMassDate(file)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err2,
			"message": fmt.Sprintf("parseBatchMassDate error: %+v\n", err),
		})
		return
	}
	// fmt.Printf("orderNum: %s, upload count: %d\n", orderNum, len(items))

	updatedCnt, err := updateMassDate(items)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err3,
			"message": fmt.Sprintf("updateMassDate error: %+v\n", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rtnCode":    ok,
		"uploadCnt":  len(items),
		"matchedCnt": updatedCnt,
	})
}

// 解析上载的文件
func parseBatchMassDate(file io.Reader) ([]massDate, error) {
	items := []massDate{}

	// 从文件中读取
	xlsx, err := excelize.OpenReader(file)
	if err != nil {
		return items, err
	}

	// sheet 的名字是固定的，只能是 Sheet1
	index := xlsx.GetSheetIndex("Sheet1")

	// Get all the rows in a sheet.
	rows := xlsx.GetRows("sheet" + strconv.Itoa(index))
	// fmt.Printf("total rows: %d\n", len(rows))

	// 物料号，保质期天数
	for idx, row := range rows {
		// 第一行，表头，跳过
		if idx == 0 {
			continue
		}

		// 第二行：料号、保质期天数为空，表示是最后一行
		if len(row[0]) == 0 || len(row[1]) == 0 {
			break
		}

		i := massDate{}
		i.InvCode = strings.TrimSpace(row[0])
		i.DayCount, err = strconv.Atoi(strings.TrimSpace(row[1]))

		// 保质期天数输入错误
		if err != nil {
			continue
		}

		items = append(items, i)
	}

	return items, nil
}

// 批量更新数据库
func updateMassDate(items []massDate) (int, error) {
	var err error
	var sqlCmd string

	// -- 有效期推算方式：日
	sqlCmd = `update Inventory_sub set iExpiratDateCalcu = 2 where 1 = 1`
	_, err = db.Exec(sqlCmd)
	if err != nil {
		return 0, err
	}

	// 有效期单位：天
	sqlCmd = `update Inventory set cMassUnit = 3 where 1 = 1 `
	_, err = db.Exec(sqlCmd)
	if err != nil {
		return 0, err
	}

	// 入库表，现存量表
	tables := []string{"RdRecords01", "RdRecords08", "RdRecords09", "RdRecords10", "RdRecords11", "RdRecords32", "RdRecords34", "currentstock"}

	// 物料主数据：有效期
	const invTempl = ` update Inventory set iMassDate = {{.MassDate}} where 1 = 1 and cInvCode = '{{.InvCode}}' `
	t1 := template.Must(template.New("invTempl").Parse(invTempl))

	buf := &bytes.Buffer{}
	done := 0
	for _, item := range items {
		massDay := strconv.Itoa(item.DayCount)
		cond := ` where 1 = 1 and cInvCode = '` + item.InvCode + `'; `

		// 物料主数据：有效期
		// sqlCmd = `update Inventory set iMassDate = ` + massDay + cond
		param := map[string]interface{}{
			"MassDate": item.DayCount,
			"InvCode":  item.InvCode,
		}
		if err := t1.Execute(buf, param); err != nil {
			panic(err)
		}
		sqlCmd := buf.String()

		_, err = db.Exec(sqlCmd)
		if err != nil {
			return 0, err
		}

		// 所有入库记录
		for _, tbl := range tables {
			// 更新保质期
			sqlCmd = " update " + tbl + ` set	iExpiratDateCalcu = 2, iMassDate = ` + massDay + cond
			_, err = db.Exec(sqlCmd)
			if err != nil {
				return 0, err
			}

			// 重新计算失效天数
			sqlCmd = " update " + tbl + ` set	dVDate =   dateadd(DAY, iMassDate + 1, dMadeDate ),
			  cExpirationDate =  CONVERT(varchar(100), dateadd(DAY, iMassDate, dMadeDate ), 23) ,
		  	dExpirationDate = dateadd(DAY, iMassDate, dMadeDate ) ` + cond

			if tbl == "currentstock" {
				sqlCmd = " update " + tbl + ` set	dVDate =   dateadd(DAY, iMassDate + 1, dMDate ),
					cExpirationDate =  CONVERT(varchar(100), dateadd(DAY, iMassDate, dMDate ), 23) ,
					dExpirationDate = dateadd(DAY, iMassDate, dMDate ) ` + cond
			}

			_, err = db.Exec(sqlCmd)
			if err != nil {
				fmt.Printf("%s-%s", tbl, item.InvCode)
				return 0, err
			}
		}

		fmt.Printf("%d: %s\n", done, item.InvCode)
		done = done + 1
	}

	return done, nil
}
