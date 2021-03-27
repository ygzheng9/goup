package matShortage

import (
// "fmt"
)

// MatInfo 物料基本信息
type MatInfo struct {
	InvCode     string  `db:"invCode" json:"invCode"`
	InvName     string  `db:"invName" json:"invName"`
	InvStd      string  `db:"invStd" json:"invStd"`
	Purchase    bool    `db:"purchase" json:"purchcase"`
	SelfMade    bool    `db:"selfMade" json:"selfMade"`
	Outsourcing bool    `db:"outsourcing" json:"outsourcing`
	Moq         float64 `db:"moq" json:"moq"`
	Leadtime    float64 `db:"leadtime" json:"leadtime"`
	McCode      string  `db:"mcCode" json:"mcCode"`
}

// getInvInfo 获取物料基本信息
func getInvInfo() ([]MatInfo, error) {
	sqlCmd := `
        select i.cInvCode invCode, isnull(i.cInvName,'') invName, isnull(i.cInvStd,'') invStd,
               isnull(i.bPurchase,0) purchase , isnull(i.bSelf,0) selfMade, isnull(i.bProxyForeign,0) outsourcing,
               isnull(i.fMinSupply,0) moq, isnull(i.iInvAdvance,0) leadtime, isnull(c.MC_CDE, '') mcCode
          from inventory i
        left join  DJPlan..T_INV_CLASS c on c.INV_CLASS = i.cInvCCode
		`
	// fmt.Printf("sqlCmd: %s\n", sqlCmd)

	items := []MatInfo{}
	err := db.Select(&items, sqlCmd)

	return items, err
}
