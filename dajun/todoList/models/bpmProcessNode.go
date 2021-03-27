package models

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	// BpmApproved 审批通过
	BpmApproved = "APPROVED"
	// BpmRejected 审批拒绝
	BpmRejected = "REJECTED"
	// BpmDelegated 授权给别人
	BpmDelegated = "DELEGATED"
	// BpmSkipped 掉过
	BpmSkipped = "SKIPPED"
	// BpmWaiting 待处理
	BpmWaiting = "WAITING"
	// BpmNonApplicable 不适用
	BpmNonApplicable = "NA"
)

// BpmProcessNode 审批流的类型
type BpmProcessNode struct {
	ID         int    `db:"ID" json:"id" form:"id"`
	ProcessID  int    `db:"PROCESS_ID" json:"processID" form:"processID"`
	NodeName   string `db:"NODE_NME" json:"ndName" form:"ndName"`
	StepNumber int    `db:"STEP_NUM" json:"stepNum" form:"stepNum"`
	NodeDesc   string `db:"NODE_DESC" json:"ndDesc" form:"ndDesc"`
	NodeType   string `db:"NODE_TYP" json:"ndType" form:"ndType"`
	NodeUser   string `db:"NODE_USR" json:"ndUser" form:"ndUser"`
	BizRule    string `db:"BIZ_RULE" json:"bizRule" form:"bizRule"`
	CalcRule   string `db:"CALC_RULE" json:"calcRule" form:"calcRule"`
	BizAction  string `db:"BIZ_ACTION" json:"bizAction" form:"bizAction"`
	UpdateUser string `db:"UPD_USR" json:"updateUser" form:"updateUser"`
	UpdateDate string `db:"UPD_DTE" json:"updateDate" form:"updateDate"`
}

const (
	bpmProcessNodeSelect = `select ID, PROCESS_ID, STEP_NUM, ifnull(NODE_NME,'') NODE_NME, ifnull(NODE_TYP,'') NODE_TYP, ifnull(NODE_USR,'') NODE_USR,
			ifnull(BIZ_RULE,'') BIZ_RULE, ifnull(CALC_RULE,'') CALC_RULE, ifnull(BIZ_ACTION,'') BIZ_ACTION,
			ifnull(NODE_DESC,'') NODE_DESC, ifnull(UPD_USR,'') UPD_USR, ifnull(UPD_DTE,'') UPD_DTE
		from t_bpm_process_node		`
)

// BpmProcessNodeFindBy 根据条件查找
func BpmProcessNodeFindBy(cond string) ([]BpmProcessNode, error) {
	sqlCmd := bpmProcessNodeSelect + cond

	items := []BpmProcessNode{}
	err := db.Select(&items, sqlCmd)

	return items, err
}

// BpmProcessNodeFindAll 返回部门清单
func BpmProcessNodeFindAll() ([]BpmProcessNode, error) {
	// 查询清单
	return BpmProcessNodeFindBy(" order by PROCESS_ID, STEP_NUM  ")
}

// BpmProcessNodeFindByID 按照 id 查询
func BpmProcessNodeFindByID(id int) (BpmProcessNode, error) {
	cmd := bpmProcessNodeSelect + ` where ID=?`
	item := BpmProcessNode{}
	err := db.Get(&item, cmd, id)
	return item, err
}

// BpmProcessNodeFindByProcessID 按照 id 查询
func BpmProcessNodeFindByProcessID(processID int) ([]BpmProcessNode, error) {
	return BpmProcessNodeFindBy(" where PROCESS_ID = " + strconv.Itoa(processID) + " order by PROCESS_ID, STEP_NUM  ")
}

// BpmProcessNodeInsert 当前对象，插入到数据库
func BpmProcessNodeInsert(c BpmProcessNode) (BpmProcessNode, error) {
	c.UpdateDate = time.Now().Format("2006-01-02 15:04:05")

	// 根据 struct 中的 DB tag 进行自动 named parameter
	cmd := `INSERT INTO t_bpm_process_node (PROCESS_ID,NODE_NME,STEP_NUM,NODE_DESC,NODE_TYP, NODE_USR, BIZ_RULE, CALC_RULE, BIZ_ACTION,  UPD_USR,UPD_DTE) VALUES
						(:PROCESS_ID,:NODE_NME,:STEP_NUM,:NODE_DESC,:NODE_TYP, :NODE_USR, :BIZ_RULE, :CALC_RULE, :BIZ_ACTION, :UPD_USR,:UPD_DTE)`
	res, err := db.NamedExec(cmd, c)
	if err != nil {
		return BpmProcessNode{}, err
	}

	// 取回插入的记录
	id, err := res.LastInsertId()
	if err != nil {
		return BpmProcessNode{}, err
	}

	return BpmProcessNodeFindByID(int(id))
}

// BpmProcessNodeUpdate 当前对象，更新到数据库
func BpmProcessNodeUpdate(c BpmProcessNode) (BpmProcessNode, error) {
	c.UpdateDate = time.Now().Format("2006-01-02 15:04:05")

	cmd := `update t_bpm_process_node set
							PROCESS_ID = :PROCESS_ID,
							NODE_NME = :NODE_NME,
							STEP_NUM = :STEP_NUM,
							NODE_DESC = :NODE_DESC,
							NODE_TYP = :NODE_TYP,
							NODE_USR = :NODE_USR,
							BIZ_RULE = :BIZ_RULE,
							CALC_RULE = :CALC_RULE,
							BIZ_ACTION = :BIZ_ACTION,
							UPD_USR = :UPD_USR,
							UPD_DTE = :UPD_DTE
						where ID=:ID`
	_, err := db.NamedExec(cmd, c)
	if err != nil {
		return c, err
	}

	return BpmProcessNodeFindByID(c.ID)
}

// BpmProcessNodeDelete 当前对象，按照 ID 从数据库删除
func BpmProcessNodeDelete(id int) error {
	// 按照 id 删除
	cmd := "delete from t_bpm_process_node where ID= " + strconv.Itoa(id)
	_, err := db.Exec(cmd)
	return err
}

// BpmProcessNodeCheckRule 检查审批节点规则是否有效
func BpmProcessNodeCheckRule(rule string) bool {
	if len(rule) == 0 {
		return true
	}

	if isDangerous(rule) {
		return false
	}

	type userT struct {
		NdUser string `db:"ndUser" json:"ndUser"`
	}
	results := []userT{}

	err := db.Select(&results, rule)
	if err != nil {
		fmt.Printf("err: %+v\n", err)
		return false
	}
	return true
}

// isDangerCmd sql语句是否包含特定子
func isDangerous(cmd string) bool {
	forbidden := []string{"insert", "update", "delete", "drop", "create"}
	sqlCmd := strings.ToLower(cmd)
	words := strings.Fields(sqlCmd)

	for _, i := range forbidden {
		for _, w := range words {
			if i == w {
				return true
			}
		}

	}

	return false
}
