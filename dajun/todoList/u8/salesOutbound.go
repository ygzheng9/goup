package u8

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	// "github.com/xuri/excelize"
	"github.com/360EntSecGroup-Skylar/excelize"
)

// OutboundT 上载的数据
type OutboundT struct {
	InvCode string
	SeqNo   string
	Matched int
}

// OutboundUpload 上载销售发货单的实际批次号
func OutboundUpload(c *gin.Context) {
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
	orderNum, items, err := ParseUploadOutbound(file)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err2,
			"message": fmt.Sprintf("parseUploadOutbound error: %+v\n", err),
		})
		return
	}
	// fmt.Printf("orderNum: %s, upload count: %d\n", orderNum, len(items))

	// 根据单号，取得行项目
	oldItems, err := fetchOutboundItems(orderNum)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err3,
			"message": fmt.Sprintf("fetchOutboundItems: %+v", err),
		})
		return
	}
	// fmt.Printf("item count: %d\n", len(oldItems))

	// 行项目，和上载的内容，按照 料号 进行匹配，更新 批次号
	count, err := matchOutboundItems(oldItems, items)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err4,
			"message": fmt.Sprintf("matchOutboundItems: %+v", err),
		})
		return
	}
	// fmt.Printf("matched count: %d\n", count)

	c.JSON(http.StatusOK, gin.H{
		"rtnCode":    ok,
		"orderNum":   orderNum,
		"itemCnt":    len(oldItems),
		"uploadCnt":  len(items),
		"matchedCnt": count,
	})
}

// ParseUploadOutbound 解析上载的文件
func ParseUploadOutbound(file io.Reader) (string, []OutboundT, error) {
	items := []OutboundT{}

	// 从文件中读取
	xlsx, err := excelize.OpenReader(file)
	if err != nil {
		return "", nil, err
	}

	// sheet 的名字是固定的，只能是 Sheet1
	// index := xlsx.GetSheetIndex("Sheet1")
	// fmt.Printf("index: %d\n", index)

	// // Get all the rows in a sheet.
	// rows := xlsx.GetRows("sheet" + strconv.Itoa(index))

	rows := xlsx.GetRows("Sheet1")
	// fmt.Printf("total rows: %d\n", len(rows))

	// 最大行数
	const maxRow = 1000

	// 发货单号
	orderNum := ""
	for idx, row := range rows {
		// 第一行，表头，跳过
		if idx == 0 {
			continue
		}

		// 第二行，第一列，是单号
		if idx == 1 {
			orderNum = row[0]
			fmt.Printf("orderNum: %s\n", orderNum)
		}

		if idx >= maxRow {
			break
		}

		// 料号、新批次号为空，表示是最后一行
		if len(row[1]) == 0 || len(row[2]) == 0 {
			break
		}

		i := OutboundT{}
		i.InvCode = strings.TrimSpace(row[1])
		i.SeqNo = strings.TrimSpace(row[2])
		// 未匹配
		i.Matched = 0

		items = append(items, i)
	}

	return orderNum, items, nil
}

// 销售发货单查询出的数据
type outboundItem struct {
	AutoID  int    `db:"AutoID"`
	InvCode string `db:"cInvCode"`
	Matched int    `db:"Matched"`
}

// 根据 ordernum 取得 items
func fetchOutboundItems(orderNum string) ([]outboundItem, error) {
	// 左边补零，共 10 位
	temp := "0000000000" + strings.TrimSpace(orderNum)
	cCode := temp[len(temp)-10:]

	//  根据 单号 取得 出库单的行项目
	sqlCmd := `select b.AutoID, b.cInvCode, 0 Matched from RdRecord32 a inner join RdRecords32 b on a.id = b.id and a.cCode = '` + cCode + `' order by b.iRowNo `
	fmt.Printf("sqlCmd: %s\n", sqlCmd)

	items := []outboundItem{}

	err := db.Select(&items, sqlCmd)
	return items, err
}

// 根据 发货单号，用上载的新批次号，更新数据库
func matchOutboundItems(items []outboundItem, newSeqs []OutboundT) (int, error) {
	i := len(items)
	j := len(newSeqs)
	// 对 每一个行项目，根据 料号，取出 新的批次号，更新回数据库，并把取出的批次号标记为已使用
	for x := 0; x < i; x++ {
		for y := 0; y < j; y++ {
			// 没有匹配过，并且料号相同，则表示匹配上了
			if newSeqs[y].Matched == 0 && items[x].InvCode == newSeqs[y].InvCode {
				// 匹配上后，更新为 已匹配
				items[x].Matched = 1
				newSeqs[y].Matched = 1

				// 更新数据库
				err := updateOutboundSeq(items[x].AutoID, newSeqs[y].SeqNo)
				if err != nil {
					fmt.Printf("updateOutboundSeq failed: %+v", err)
					return 0, err
				}

				// 匹配上后，匹配下一个
				break
			}
		}
	}

	// 该单子下，匹配了多少项，
	count := 0
	for _, i := range items {
		if i.Matched == 1 {
			count = count + 1
		}
	}

	return count, nil
}

// updateOutboundSeq 根据 autoID，更新 cDefine28 字段
func updateOutboundSeq(autoID int, newSeq string) error {
	param := struct {
		AutoID int    `db:"AutoID"`
		SeqNo  string `db:"SeqNo"`
	}{AutoID: autoID, SeqNo: newSeq}

	// fmt.Printf("param: %+v\n", param)

	// 根据 struct 中的 tag 进行自动 named parameter
	sqlCmd := `update RdRecords32 set cDefine28 = :SeqNo where autoID=:AutoID`

	_, err := db.NamedExec(sqlCmd, param)
	return err
}

//////////////////////////////////

// OutboundUploadSvcParam post 的参数
type OutboundUploadSvcParam struct {
	OrderNum string      `json:"orderNum"`
	Items    []OutboundT `json:"items"`
}

// OutboundUploadSvc 前端解析 excel 后，发送 post 请求
func OutboundUploadSvc(c *gin.Context) {
	// 定义 post 中的数据结构
	param := OutboundUploadSvcParam{}

	// post 过来的 json，绑定到 对象上；
	err := c.BindJSON(&param)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("BindJson error: %+v", err),
		})
		return
	}
	fmt.Printf("param: %+v\n", param)

	orderNum := param.OrderNum
	items := param.Items

	// 根据单号，取得行项目
	oldItems, err := fetchOutboundItems(orderNum)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err3,
			"message": fmt.Sprintf("fetchOutboundItems: %+v", err),
		})
		return
	}
	// fmt.Printf("item count: %d\n", len(oldItems))

	// 行项目，和上载的内容，按照 料号 进行匹配，更新 批次号
	count, err := matchOutboundItems(oldItems, items)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err4,
			"message": fmt.Sprintf("matchOutboundItems: %+v", err),
		})
		return
	}
	// fmt.Printf("matched count: %d\n", count)

	c.JSON(http.StatusOK, gin.H{
		"rtnCode":    ok,
		"orderNum":   orderNum,
		"itemCnt":    len(oldItems),
		"uploadCnt":  len(items),
		"matchedCnt": count,
	})
}
