package models

import (
	"strconv"
	"time"
)

// Milestone 合同付款计划
type Milestone struct {
	ID         int    `db:"ID" json:"id" form:"id"`
	ContractID int    `db:"CTRCT_ID" json:"contract_id" form:"contract_id"`
	ItemName   string `db:"ITM_NME" json:"item_name" form:"item_name"`
	PlanDate   string `db:"PLN_DTE" json:"planned_date" form:"planned_date"`
	IsCleared  string `db:"IS_CLR" json:"is_cleared" form:"is_cleared"`
	Remark     string `db:"RMK" json:"remark" form:"remark"`
	UpdateUser string `db:"UPD_USR" json:"update_user" form:"update_user"`
	UpdateDate string `db:"UPD_DTE" json:"update_date" form:"update_date"`
}

const (
	milestoneSelect = `select ID, ifnull(CTRCT_ID,'') CTRCT_ID, ifnull(ITM_NME,'') ITM_NME, ifnull(PLN_DTE,'') PLN_DTE, ifnull(IS_CLR,'') IS_CLR, ifnull(RMK,'') RMK, ifnull(UPD_USR,'') UPD_USR, ifnull(UPD_DTE,'') UPD_DTE	from t_ctrct_plan`
)

// MilestoneFindBy 根据条件查找
func MilestoneFindBy(cond string) ([]Milestone, error) {
	sqlCmd := milestoneSelect + cond

	items := []Milestone{}
	err := db.Select(&items, sqlCmd)

	return items, err
}

// MilestoneFindAll 返回部门清单
func MilestoneFindAll() ([]Milestone, error) {
	// 查询清单
	return MilestoneFindBy(" order by PLN_DTE   ")
}

// MilestoneFindByContract 根据合同id查询
func MilestoneFindByContract(contractID int) ([]Milestone, error) {
	return MilestoneFindBy(" where CTRCT_ID = " + strconv.Itoa(contractID) + " order by PLN_DTE   ")
}

// MilestoneFindByID 按照 id 查询
func MilestoneFindByID(id int) (Milestone, error) {
	cmd := milestoneSelect + ` where ID=?`
	item := Milestone{}
	err := db.Get(&item, cmd, id)
	return item, err
}

// MilestoneInsert 当前对象，插入到数据库
func MilestoneInsert(c Milestone) (Milestone, error) {
	c.UpdateDate = time.Now().Format("2006-01-02 15:04:05")

	// 根据 struct 中的 DB tag 进行自动 named parameter
	cmd := `INSERT INTO t_ctrct_plan (CTRCT_ID,ITM_NME,PLN_DTE,IS_CLR,RMK,UPD_USR,UPD_DTE) VALUES
						(:CTRCT_ID,:ITM_NME,:PLN_DTE,:IS_CLR,:RMK,:UPD_USR,:UPD_DTE)`
	res, err := db.NamedExec(cmd, c)
	if err != nil {
		return Milestone{}, err
	}

	// 取回插入的记录
	id, err := res.LastInsertId()
	if err != nil {
		return Milestone{}, err
	}

	return MilestoneFindByID(int(id))
}

// MilestoneUpdate 当前对象，更新到数据库
func MilestoneUpdate(c Milestone) (Milestone, error) {
	c.UpdateDate = time.Now().Format("2006-01-02 15:04:05")

	cmd := `update t_ctrct_plan set
							CTRCT_ID = :CTRCT_ID,
							ITM_NME = :ITM_NME,
							PLN_DTE = :PLN_DTE,
							IS_CLR = :IS_CLR,
							RMK = :RMK,
							UPD_USR = :UPD_USR,
							UPD_DTE = :UPD_DTE
						where ID=:ID`
	_, err := db.NamedExec(cmd, c)
	if err != nil {
		return c, err
	}

	return MilestoneFindByID(c.ID)
}

// MilestoneDelete 当前对象，按照 ID 从数据库删除
func MilestoneDelete(c Milestone) error {
	// 按照 id 删除
	cmd := "delete from t_ctrct_inv where CTRCT_ID = " + strconv.Itoa(c.ID)
	_, err := db.Exec(cmd)

	cmd = "delete from t_ctrct_plan where ID=:ID"
	_, err = db.NamedExec(cmd, c)
	return err
}
