package models

import (
	"strconv"
	"time"
)

// Contract 合同
type Contract struct {
	ID            int    `db:"ID" json:"id" form:"id"`
	ContractNo    string `db:"CTRCT_NUM" json:"contract_no" form:"contract_no"`
	ContractTitle string `db:"CTRCT_TITLE" json:"contract_title" form:"contract_title"`
	VendorName    string `db:"VNDR_NME" json:"vendor_name" form:"vendor_name"`
	FromDate      string `db:"FM_DTE" json:"from_date" form:"from_date"`
	ToDate        string `db:"TO_DTE" json:"to_date" form:"to_date"`
	Remark        string `db:"RMK" json:"remark" form:"remark"`
	UpdateUser    string `db:"UPD_USR" json:"update_user" form:"update_user"`
	UpdateDate    string `db:"UPD_DTE" json:"update_date" form:"update_date"`
}

const (
	contractSelect = `select ID, ifnull(CTRCT_NUM,'') CTRCT_NUM, ifnull(CTRCT_TITLE,'') CTRCT_TITLE, ifnull(VNDR_NME,'') VNDR_NME, ifnull(FM_DTE,'') FM_DTE, ifnull(TO_DTE,'') TO_DTE, ifnull(RMK,'') RMK, ifnull(UPD_USR,'') UPD_USR, ifnull(UPD_DTE,'') UPD_DTE from t_ctrct`
)

// ContractFindBy 根据条件查找
func ContractFindBy(cond string) ([]Contract, error) {
	sqlCmd := contractSelect + cond

	items := []Contract{}
	err := db.Select(&items, sqlCmd)

	return items, err
}

// ContractFindAll 返回部门清单
func ContractFindAll() ([]Contract, error) {
	// 查询清单
	return ContractFindBy(" order by FM_DTE desc  ")
}

// ContractFindByID 按照 id 查询
func ContractFindByID(id int) (Contract, error) {
	cmd := contractSelect + ` where ID=?`
	item := Contract{}
	err := db.Get(&item, cmd, id)
	return item, err
}

// ContractInsert 当前对象，插入到数据库
func ContractInsert(c Contract) (Contract, error) {
	c.UpdateDate = time.Now().Format("2006-01-02 15:04:05")

	// 根据 struct 中的 DB tag 进行自动 named parameter
	cmd := `INSERT INTO t_ctrct (CTRCT_NUM,CTRCT_TITLE,VNDR_NME,FM_DTE,TO_DTE,RMK,UPD_USR,UPD_DTE) VALUES
						(:CTRCT_NUM,:CTRCT_TITLE,:VNDR_NME,:FM_DTE,:TO_DTE,:RMK,:UPD_USR,:UPD_DTE)`
	res, err := db.NamedExec(cmd, c)
	if err != nil {
		return Contract{}, err
	}

	// 取回插入的记录
	id, err := res.LastInsertId()
	if err != nil {
		return Contract{}, err
	}

	return ContractFindByID(int(id))
}

// ContractUpdate 当前对象，更新到数据库
func ContractUpdate(c Contract) (Contract, error) {
	c.UpdateDate = time.Now().Format("2006-01-02 15:04:05")

	cmd := `update t_ctrct set
						CTRCT_NUM = :CTRCT_NUM,
						CTRCT_TITLE = :CTRCT_TITLE,
						VNDR_NME = :VNDR_NME,
						FM_DTE = :FM_DTE,
						TO_DTE = :TO_DTE,
						RMK = :RMK,
						UPD_USR = :UPD_USR,
						UPD_DTE = :UPD_DTE
						where ID=:ID`
	_, err := db.NamedExec(cmd, c)
	if err != nil {
		return c, err
	}

	return ContractFindByID(c.ID)
}

// ContractDelete 当前对象，按照 ID 从数据库删除
func ContractDelete(c Contract) error {
	// 按照 id 删除
	cmd := "delete from t_ctrct_plan where CTRCT_ID = " + strconv.Itoa(c.ID)
	_, err := db.Exec(cmd)

	cmd = "delete from t_ctrct_inv where CTRCT_ID = " + strconv.Itoa(c.ID)
	_, err = db.Exec(cmd)

	cmd = "delete from t_ctrct where ID=:ID"
	_, err = db.NamedExec(cmd, c)
	return err
}
