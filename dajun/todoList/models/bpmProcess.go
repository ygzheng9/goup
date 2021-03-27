package models

import (
	"strconv"
	"time"
)

// BpmProcess 审批流的类型
type BpmProcess struct {
	ID          int    `db:"ID" json:"id" form:"id"`
	BizType     string `db:"BIZ_TYP" json:"bizType" form:"bizType"`
	StatusField string `db:"STS_FIELD" json:"statusField" form:"statusField"`
	KeyField    string `db:"KEY_FIELD" json:"keyField" form:"keyField"`
	ProcessName string `db:"PROCESS_NME" json:"processName" form:"processName"`
	PriorityNum string `db:"PRIORITY_NUM" json:"priorityNum" form:"priorityNum"`
	ValidInd    string `db:"VALID_IND" json:"validInd" form:"validInd"`
	PrcessDesc  string `db:"PROCESS_DESC" json:"processDesc" form:"processDesc"`
	BizRule     string `db:"BIZ_RULE" json:"bizRule" form:"bizRule"`
	UpdateUser  string `db:"UPD_USR" json:"updateUser" form:"updateUser"`
	UpdateDate  string `db:"UPD_DTE" json:"updateDate" form:"updateDate"`
}

const (
	bpmProcessSelect = `select ID, ifnull(BIZ_TYP,'') BIZ_TYP, ifnull(STS_FIELD,'') STS_FIELD, ifnull(KEY_FIELD,'') KEY_FIELD, ifnull(PRIORITY_NUM,'') PRIORITY_NUM,  ifnull(PROCESS_NME,'') PROCESS_NME, ifnull(VALID_IND,'') VALID_IND, ifnull(PROCESS_DESC,'') PROCESS_DESC, ifnull(BIZ_RULE,'') BIZ_RULE, ifnull(UPD_USR,'') UPD_USR, ifnull(UPD_DTE,'') UPD_DTE
		from t_bpm_process`
)

// BpmProcessFindBy 根据条件查找
func BpmProcessFindBy(cond string) ([]BpmProcess, error) {
	sqlCmd := bpmProcessSelect + cond

	// fmt.Printf("BpmProcessFindBy: %s\n", sqlCmd)

	items := []BpmProcess{}
	err := db.Select(&items, sqlCmd)

	return items, err
}

// BpmProcessFindAll 返回部门清单
func BpmProcessFindAll() ([]BpmProcess, error) {
	// 查询清单
	return BpmProcessFindBy(" order by BIZ_TYP, VALID_IND, PRIORITY_NUM ")
}

// BpmProcessFindByID 按照 id 查询
func BpmProcessFindByID(id int) (BpmProcess, error) {
	cmd := bpmProcessSelect + ` where ID=?`
	item := BpmProcess{}
	err := db.Get(&item, cmd, id)
	return item, err
}

// BpmProcessInsert 当前对象，插入到数据库
func BpmProcessInsert(c BpmProcess) (BpmProcess, error) {
	// 默认都有效
	c.ValidInd = "Y"

	c.UpdateDate = time.Now().Format("2006-01-02 15:04:05")

	// 根据 struct 中的 DB tag 进行自动 named parameter
	cmd := `INSERT INTO t_bpm_process (BIZ_TYP,STS_FIELD,KEY_FIELD, PROCESS_NME,PRIORITY_NUM,VALID_IND,PROCESS_DESC,BIZ_RULE,UPD_USR,UPD_DTE) VALUES
						(:BIZ_TYP,:STS_FIELD,:KEY_FIELD,:PROCESS_NME,:PRIORITY_NUM,:VALID_IND,:PROCESS_DESC,:BIZ_RULE,:UPD_USR,:UPD_DTE)`
	res, err := db.NamedExec(cmd, c)
	if err != nil {
		return BpmProcess{}, err
	}

	// 取回插入的记录
	id, err := res.LastInsertId()
	if err != nil {
		return BpmProcess{}, err
	}

	return BpmProcessFindByID(int(id))
}

// BpmProcessUpdate 当前对象，更新到数据库
func BpmProcessUpdate(c BpmProcess) (BpmProcess, error) {
	c.UpdateDate = time.Now().Format("2006-01-02 15:04:05")

	cmd := `update t_bpm_process set
							BIZ_TYP = :BIZ_TYP,
							STS_FIELD = :STS_FIELD,
							KEY_FIELD = :KEY_FIELD,
							PROCESS_NME = :PROCESS_NME,
							PRIORITY_NUM = :PRIORITY_NUM,
							VALID_IND = :VALID_IND,
							PROCESS_DESC = :PROCESS_DESC,
							BIZ_RULE = :BIZ_RULE,
							UPD_USR = :UPD_USR,
							UPD_DTE = :UPD_DTE
						where ID=:ID`
	_, err := db.NamedExec(cmd, c)
	if err != nil {
		return c, err
	}

	return BpmProcessFindByID(c.ID)
}

// BpmProcessDelete 当前对象，按照 ID 从数据库删除
func BpmProcessDelete(id int) error {
	// 按照 id 删除
	strID := strconv.Itoa(id)

	cmd := "delete from t_bpm_process_node where PROCESS_ID = " + strID
	_, err := db.Exec(cmd)

	cmd = "delete from t_bpm_process where ID=" + strID
	_, err = db.Exec(cmd)
	return err
}

// BpmProcessCheckRule 检查流程规则是否有效
func BpmProcessCheckRule(tableName string, rule string) bool {
	if len(rule) == 0 {
		return true
	}

	if isDangerous(rule) {
		return false
	}

	sqlCmd := "select count(1) from " + tableName + " where 1 = 0 and (" + rule + ")"
	_, err := db.Exec(sqlCmd)
	return err == nil
}
