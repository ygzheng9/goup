package u8

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// POInvDiff 生产订单子件需求 - 生产备料仓库存
type POInvDiff struct {
	InvCode string  `db:"invCode" json:"invCode"`
	InvName string  `db:"invName" json:"invName"`
	Diff    float64 `db:"diff" json:"diff"`
}

// FetchPoInvDiff 生产订单子件需求 - 生产备料仓库存
func FetchPoInvDiff(c *gin.Context) {
	// 从文件系统加载 sql
	sqlCmd := loadFile("./templates/poinvdiff.sql")
	fmt.Println(sqlCmd)

	items := []POInvDiff{}
	err := db.Select(&items, sqlCmd)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("%+v", err),
		})

		fmt.Println(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rtnCode": ok,
		"items":   items,
		// "message": fmt.Sprintf("%s", sqlCmd),
	})
}

// loadFile 从文件系统，加载文件
func loadFile(fileName string) string {
	var buffer bytes.Buffer

	// 打开文件
	f, err := os.Open(fileName)
	if err != nil {
		return ""
	}
	// 使用 buf
	buf := bufio.NewReader(f)
	for {
		// 逐行读取
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		if err != nil {
			if err == io.EOF {
				// 读取完毕后的结果
				return buffer.String()
			}
			// 这里是多余的一句
			return buffer.String()
		}

		// 保留原文件的样子，一行一行
		buffer.WriteString(line + "\n")
	}
}
