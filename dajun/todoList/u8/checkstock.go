package u8

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CheckStock 检查当前库存量
func CheckStock(c *gin.Context) {
	type paramT struct {
		MatCode string `json:"mat_code" form:"mat_code"`
		MatCond string `json:"mat_cond" form:"mat_cond"`
	}

	param := paramT{}
	err := c.BindJSON(&param)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("CheckStock BindJson error: %+v\n", err),
		})
		return
	}

	// 查找一次
	qty, err := DoCheckStock(param.MatCode, param.MatCond)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err2,
			"message": fmt.Sprintf("doCheckStock error: %+v\n", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rtnCode": ok,
		"qty":     qty,
	})
}

// DoCheckStock 检查现存量是否满足规则
func DoCheckStock(mat, cond string) (float64, error) {
	// 拼接查询字符串
	sqlCmd := ` select sum(iQuantity) qty from currentStock where cInvCode = '` + mat +
		`' and cWhCode in ('01', '02', '05') having sum(iQuantity) ` + cond
	// log.Printf("sqlCmd: %s\n", sqlCmd)

	// 情况1：拼接的 sql 有错误
	// 情况2：拼接的 sql 无错误，但是条件不满足
	// 情况3：拼接的 sql 无错误，但是条件满足

	rows, err := db.Query(sqlCmd)
	if err != nil {
		// sql 语句有错误
		return -1, err
	}

	// 一次只查一个物料，所以只需要取第一条记录
	for rows.Next() {
		var d decimal
		// mssql 中的 decimal 返回的数据是 []uint8，而不是 数字
		// 需要转换
		err = rows.Scan(&d)
		if err != nil {
			// 从数据库中读取数据
			return -1, err
		}

		result, err := parseDecimal(d)
		if err != nil {
			// 转换成数字时出错，应该不会发生
			return -1, err
		}
		return result, nil
	}

	// sql 语句没错，但是没有满足条件的记录，返回 -1
	return -1, nil
}

// GetCurrentStock 返回物料的当前库存
func GetCurrentStock(invCodes []string) ([]CurrentStock, error) {
	items := []CurrentStock{}
	if len(invCodes) == 0 {
		return items, nil
	}

	// 拼接 in 条件；
	var cond = " and cInvCode in ( "
	for _, code := range invCodes {
		cond = cond + "'" + code + "',"
	}
	cond = cond + "'') group by cInvCode order by cInvCode "

	sqlCmd := ` select cInvCode, sum(iQuantity) currQty from currentStock where cWhCode in ('01', '02', '05') ` + cond

	err := db.Select(&items, sqlCmd)

	return items, err
}
