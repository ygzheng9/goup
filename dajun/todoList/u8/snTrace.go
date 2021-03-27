package u8

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SNInfo 序列号查询返回结果
type SNInfo struct {
	SN         string `db:"sn" json:"sn"`
	SoWarranty string `db:"soWarranty" json:"soWarranty"`

	BatchNo    string `db:"batchNo" json:"batchNo"`
	InvCode    string `db:"invCode" json:"invCode"`
	InvName    string `db:"invName" json:"invName"`
	CusInvName string `db:"cusInvName" json:"cusInvName"`
	SoCode     string `db:"soCode" json:"soCode"`
	SoDate     string `db:"soDate" json:"soDate"`
	CustCode   string `db:"custCode" json:"custCode"`
	CustName   string `db:"custName" json:"custName"`
	SoMemo     string `db:"soMemo" json:"soMemo"`
	CustAddr   string `db:"custAddr" json:"custAddr"`

	DeptName     string `db:"deptName" json:"deptName"`
	SoMaker      string `db:"soMaker" json:"soMaker"`
	SoPersonCode string `db:"soPersonCode" json:"soPersonCode"`
	SoPersonName string `db:"soPersonName" json:"soPersonName"`
	SoItemNo     string `db:"soItemNo" json:"soItemNo"`
	ItemCode     string `db:"itemCode" json:"itemCode"`
	ItemName     string `db:"itemName" json:"itemName"`

	StockOutCode    string `db:"stockOutCode" json:"stockOutCode"`
	StockOutDt      string `db:"stockOutDt" json:"stockOutDt"`
	StockOutHandler string `db:"stockOutHandler" json:"stockOutHandler"`
	StockOutMemo    string `db:"stockOutMemo" json:"stockOutMemo"`
	StockOutMaker   string `db:"stockOutMaker" json:"stockOutMaker"`
}

// TraceSNParam 所需的参数
type TraceSNParam struct {
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	CustCode  string `json:"custCode"`
	SoCode    string `json:"soCode"`
	SN        string `json:"sn"`
}

// TraceSN 查找质保期
func TraceSN(c *gin.Context) {
	// 查询参数
	param := TraceSNParam{}

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

	items, err := findSNInfo(param)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err2,
			"message": fmt.Sprintf("findSNInfo: %+v", err),
		})
		return
	}
	// fmt.Printf("matched count: %d\n", count)

	c.JSON(http.StatusOK, gin.H{
		"rtnCode": ok,
		"items":   items,
	})
}

func findSNInfo(param TraceSNParam) ([]SNInfo, error) {
	// 拼接 sql
	base := `
		select
			isnull(b.cDefine28,'') sn, isnull(e.cDefine13, '') 'soWarranty',
			isnull(d.cInvCode,'') invCode, isnull(d.cInvName,'') invName, isnull(b.cBatch,'') batchNo,  isnull(b.ccusinvname, '') cusInvName,
			isnull(e.cSoCode,'') soCode, isnull(e.cCusCode,'') custCode, isnull(e.cCusName,'') custName, isnull(e.cCusOAddress, '') custAddr, isnull(f.cDepName,'') deptName, isnull(e.cMaker,'') soMaker, isnull(e.cPersonCode,'') soPersonCode, isnull(g.cPersonName,'') soPersonName,
			isnull(e.dDate,'') soDate, isnull(e.cMemo,'') soMemo,
			isnull(d.iRowNo,'') soItemNo, isnull(d.cItemCode,'') itemCode, isnull(d.cItemName,'') itemName,
			isnull(a.cCode,'') stockOutCode, isnull(a.dDate,'') stockOutDt, isnull(a.cHandler,'') stockOutHandler, isnull(a.cMemo,'') stockOutMemo, isnull(a.cMaker,'') stockOutMaker
		from RdRecord32 a
		inner join RdRecords32 b on a.id = b.id
		inner join DispatchLists c on  c.iDLsID =  b.iDLsID
		inner join SO_SODetails d on d.iSOsID =  c.iSOsID
		inner join SO_SOMain e on e.ID = d.ID
		inner join Department f on f.cDepCode = e.cDepCode
		inner join Person g on  g.cPersonCode = e.cPersonCode
		where 1 = 1
	`

	cond := fmt.Sprintf(" and e.dDate >= '%s' and e.dDate <= '%s'", param.StartDate, param.EndDate)

	if len(param.CustCode) > 0 {
		cond = cond + fmt.Sprintf(" and e.cCusCode = '%s'", param.CustCode)
	}

	if len(param.SoCode) > 0 {
		// cond = cond + fmt.Sprintf(" and e.cSoCode = '%s'", param.SoCode)
		// 如果输入了订单号，那么忽略掉 时间，客户 的约束条件
		cond = fmt.Sprintf(" and e.cSoCode = '%s'", param.SoCode)

	}

	if len(param.SN) > 0 {
		// 如果输入了 钢印号，忽略掉 时间、客户、单号
		cond = fmt.Sprintf(" and b.cDefine28 = '%s'", param.SN)
	}

	sqlCmd := base + cond + " order by b.cDefine28;"
	// fmt.Printf("sql: %s\n", sqlCmd)

	items := []SNInfo{}
	err := db.Select(&items, sqlCmd)
	return items, err
}
