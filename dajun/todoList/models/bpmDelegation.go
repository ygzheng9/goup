package models

import (
	"strconv"
	"time"
)

// BpmDelegation 授权，授权不允许传递
type BpmDelegation struct {
	ID         int    `db:"ID" json:"id" form:"id"`
	UserName   string `db:"USR_NME" json:"userName" form:"userName"`
	DelegateTo string `db:"DELEGATE_TO" json:"delegateTo" form:"delegateTo"`
	StartDate  string `db:"START_DTE" json:"startDate" form:"startDate"`
	EndDate    string `db:"END_DTE" json:"endDate" form:"endDate"`
	ValidInd   string `db:"VALID_IND" json:"validInd" form:"validInd"`
}

const (
	bpmDelegationSelect = `select ID, ifnull(USR_NME,'') USR_NME, ifnull(DELEGATE_TO,'') DELEGATE_TO, ifnull(START_DTE,'') START_DTE,ifnull(END_DTE,'') END_DTE,  ifnull(VALID_IND,'') VALID_IND
		from t_bpm_delegation`
)

// BpmDelegationFindBy 根据条件查找
func BpmDelegationFindBy(cond string) ([]BpmDelegation, error) {
	sqlCmd := bpmDelegationSelect + cond

	items := []BpmDelegation{}
	err := db.Select(&items, sqlCmd)

	return items, err
}

// BpmDelegationFindAll 返回部门清单
func BpmDelegationFindAll() ([]BpmDelegation, error) {
	// 查询清单
	return BpmDelegationFindBy(" ")
}

// BpmDelegationFindByID 按照 id 查询
func BpmDelegationFindByID(id int) (BpmDelegation, error) {
	cmd := bpmDelegationSelect + ` where ID=?`
	item := BpmDelegation{}
	err := db.Get(&item, cmd, id)
	return item, err
}

// BpmDelegationInsert 当前对象，插入到数据库
func BpmDelegationInsert(c BpmDelegation) (BpmDelegation, error) {
	// 默认都有效
	c.ValidInd = "Y"
	c.StartDate = time.Now().Format("2006-01-02 15:04:05")

	// 根据 struct 中的 DB tag 进行自动 named parameter
	cmd := `INSERT INTO t_bpm_delegation (USR_NME,DELEGATE_TO,START_DTE,END_DTE,VALID_IND) VALUES
						(:USR_NME,:DELEGATE_TO,:START_DTE,:END_DTE,:VALID_IND)`
	res, err := db.NamedExec(cmd, c)
	if err != nil {
		return BpmDelegation{}, err
	}

	// 取回插入的记录
	id, err := res.LastInsertId()
	if err != nil {
		return BpmDelegation{}, err
	}

	return BpmDelegationFindByID(int(id))
}

// BpmDelegationComplete 当前对象，更新到数据库
func BpmDelegationComplete(id int) error {
	now := time.Now().Format("2006-01-02 15:04:05")
	cmd := " update t_bpm_delegation set VALID_IND = 'N', END_DTE='" + now + "' where ID=" + strconv.Itoa(id)
	_, err := db.Exec(cmd)
	return err
}

// BpmDelegationDelete 当前对象，按照 ID 从数据库删除
func BpmDelegationDelete(id int) error {
	// 按照 id 删除
	strID := strconv.Itoa(id)

	cmd := "delete from t_bpm_delegation where ID=" + strID
	_, err := db.Exec(cmd)
	return err
}
