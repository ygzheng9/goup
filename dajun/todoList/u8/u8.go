package u8

import (
	"bytes"
	"strconv"

	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

// SetupDB 把数据库连接传入
func SetupDB(dbConn *sqlx.DB) {
	db = dbConn
}

// GetDB 返回 U8 数据库连接
func GetDB() *sqlx.DB {
	return db
}

const (
	ok   = 0
	err1 = 1
	err2 = 2
	err3 = 3
	err4 = 4
	err5 = 5
	err6 = 6
	err7 = 7
	err8 = 8
	err9 = 9
)

// CurrentStock 物料的当前库存
type CurrentStock struct {
	InvCode string  `db:"cInvCode" json:"inv_code" form:"inv_code"`
	CurrQty float64 `db:"currQty" json:"curr_qty" form:"curr_qty"`
}

// u8 对应 u8 中的 decimal 类型
type decimal []uint8

// parseDecimal 把数组（存的是ascii值），转换成真正的数字
func parseDecimal(val decimal) (float64, error) {
	var buffer bytes.Buffer
	for _, v := range val {
		r := rune(v)
		buffer.WriteString(string(r))
	}

	result, err := strconv.ParseFloat(buffer.String(), 64)
	if err != nil {
		return -1, err
	}

	return result, err
}
