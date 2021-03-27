package models

import (
	"bytes"
	"fmt"
	"strconv"
	"time"
)

// BpmInstanceNode 审批实例
type BpmInstanceNode struct {
	ID            int    `db:"ID" json:"id" form:"id"`
	InstanceID    int    `db:"INS_ID" json:"instanceID" form:"instanceID"`
	StepNumber    int    `db:"STEP_NUM" json:"stepNum" form:"stepNum"`
	CalcRule      string `db:"CALC_RULE" json:"calcRule" form:"calcRule"`
	CurrentUser   string `db:"CURR_USR" json:"currentUser" form:"currentUser"`
	CurrentResult string `db:"CURR_RESULT" json:"currentResult" form:"currentResult"`
	CurrentRemark string `db:"CURR_RMK" json:"currentRemark" form:"currentRemark"`
	ArrivalDate   string `db:"ARR_DTE" json:"arrivalDate" form:"arrivalDate"`
	CompleteDate  string `db:"CMPL_DTE" json:"completeDate" form:"completeDate"`
}

const (
	bpmInstanceNodeSelect = `select ID, INS_ID, ifnull(STEP_NUM,'') STEP_NUM,  ifnull(CALC_RULE,'') CALC_RULE, ifnull(CURR_RESULT,'') CURR_RESULT, ifnull(CURR_USR,'') CURR_USR, ifnull(CMPL_DTE,'') CMPL_DTE,ifnull(ARR_DTE,'') ARR_DTE,ifnull(CURR_RMK,'') CURR_RMK from t_bpm_ins_node `
)

// BpmInstanceNodeFindBy 根据条件查找
func BpmInstanceNodeFindBy(cond string) ([]BpmInstanceNode, error) {
	sqlCmd := bpmInstanceNodeSelect + cond

	// fmt.Printf("BpmInstanceNodeFindBy: %s\n", sqlCmd)

	items := []BpmInstanceNode{}
	err := db.Select(&items, sqlCmd)

	return items, err
}

// BpmInstanceNodeFindAll 返回部门清单
func BpmInstanceNodeFindAll() ([]BpmInstanceNode, error) {
	// 查询清单
	return BpmInstanceNodeFindBy(" ")
}

// BpmInstanceNodeFindByID 按照 id 查询
func BpmInstanceNodeFindByID(id int) (BpmInstanceNode, error) {
	cmd := bpmInstanceNodeSelect + ` where ID=?`
	item := BpmInstanceNode{}
	err := db.Get(&item, cmd, id)
	return item, err
}

// BpmInstanceNodeInsert 当前对象，插入到数据库
func BpmInstanceNodeInsert(c BpmInstanceNode) (BpmInstanceNode, error) {
	// 默认都有效
	if len(c.CurrentResult) == 0 {
		c.CurrentResult = BpmWaiting
	}

	if len(c.CalcRule) == 0 {
		c.CalcRule = "FIRST_ONE"
	}

	c.ArrivalDate = time.Now().Format("2006-01-02 15:04:05")

	// 根据 struct 中的 DB tag 进行自动 named parameter
	cmd := `INSERT INTO t_bpm_ins_node (INS_ID,STEP_NUM,CALC_RULE,CURR_USR,CURR_RESULT,CMPL_DTE,ARR_DTE) VALUES
						(:INS_ID,:STEP_NUM,:CALC_RULE, :CURR_USR,:CURR_RESULT,:CMPL_DTE,:ARR_DTE)`
	res, err := db.NamedExec(cmd, c)
	if err != nil {
		return BpmInstanceNode{}, err
	}

	// 取回插入的记录
	id, err := res.LastInsertId()
	if err != nil {
		return BpmInstanceNode{}, err
	}

	return BpmInstanceNodeFindByID(int(id))
}

// BpmInstanceNodeDelegate 授权
func BpmInstanceNodeDelegate(c BpmInstanceNode, delegation string) (BpmInstanceNode, error) {
	// 当前审批人，标记 已代理
	now := time.Now().Format("2006-01-02 15:04:05")
	c.CurrentResult = BpmDelegated
	c.CompleteDate = now
	c.CurrentRemark = "授权至：" + delegation

	updCmd := ` UPDATE t_bpm_ins_node SET
								CURR_RESULT=:CURR_RESULT,
								CURR_RMK=:CURR_RMK,
								CMPL_DTE=:CMPL_DTE
							WHERE ID=:ID  `
	_, err := db.NamedExec(updCmd, c)
	if err != nil {
		return BpmInstanceNode{}, err
	}

	// 对代理人，增加记录
	d := BpmInstanceNode{
		InstanceID:  c.InstanceID,
		StepNumber:  c.StepNumber,
		CurrentUser: delegation,
		CalcRule:    c.CalcRule,
	}
	return BpmInstanceNodeInsert(d)
}

// BpmInstanceNodeDelete 当前对象，按照 ID 从数据库删除
func BpmInstanceNodeDelete(id int) error {
	// 按照 id 删除
	strID := strconv.Itoa(id)

	cmd := "delete from t_bpm_ins_node where INS_ID = " + strID
	_, err := db.Exec(cmd)

	return err
}

// BpmInstanceNodeMark 当前审批人审批完成
func BpmInstanceNodeMark(id int, result, remark string) error {
	now := time.Now().Format("2006-01-02 15:04:05")
	var buf bytes.Buffer
	buf.WriteString(" update t_bpm_ins_node set CURR_RESULT='" + result + "' ")
	buf.WriteString(" ,CURR_RMK='" + remark + "' ")
	buf.WriteString(" ,CMPL_DTE='" + now + "'")
	buf.WriteString(" where ID = " + strconv.Itoa(id))

	_, err := db.Exec(buf.String())
	if err != nil {
		fmt.Printf("Err: BpmInstanceNodeMark, %+v\n", err)
	}
	return err
}

// BpmInstanceNodeComplete 当前审批人审批完成
func BpmInstanceNodeComplete(id int, result, remark string) {
	// 设置当前的审批记录
	err := BpmInstanceNodeMark(id, result, remark)
	if err != nil {
		return
	}

	// 加载当前审批记录
	n, err := BpmInstanceNodeFindByID(id)
	if err != nil {
		return
	}

	// 加载同一节点的其他审批记录
	var buf bytes.Buffer
	buf.WriteString(" where INS_ID=" + strconv.Itoa(n.InstanceID))
	buf.WriteString(" and STEP_NUM=" + strconv.Itoa(n.StepNumber))
	buf.WriteString(" and ID <> " + strconv.Itoa(n.ID))
	peers, err := BpmInstanceNodeFindBy(buf.String())
	if err != nil {
		return
	}
	peersCount := len(peers)

	// 给节点，只有当前一个审批人
	if peersCount == 0 {
		if result == BpmApproved {
			err = BpmInstanceForward(n.InstanceID)
			if err != nil {
				fmt.Printf("ERR: no peer, approved, BpmInstanceForward: %+v\n", err)
			}
			return
		}

		if result == BpmRejected {
			err = BpmInstanceFinish(n.InstanceID, BpmRejected)
			if err != nil {
				fmt.Printf("ERR: no peer, rejected, BpmInstanceReject: %+v\n", err)
			}
			return
		}
	}

	// 该节点有多个审批人，再查看那些审批人中，是否还有未处理的
	if peersCount > 0 {
		waitingList := []BpmInstanceNode{}
		for _, v := range peers {
			if v.CurrentResult == BpmWaiting {
				waitingList = append(waitingList, v)
			}
		}
		waitingCount := len(waitingList)

		if n.CalcRule == "ONE_PASS" {
			// 一票通过
			if result == BpmApproved {
				err = BpmInstanceForward(n.InstanceID)
				if err != nil {
					fmt.Printf("ERR: peers, ONE_PASS, APPROVED: %+v\n", err)
				}

				for _, v := range waitingList {
					BpmInstanceNodeMark(v.ID, "SKIPPED", "")
				}
			}

			if result == BpmRejected {
				if waitingCount == 0 {
					// 一票通过，所有人都处理过，但是还没有被通过，含义就是要被拒绝
					err = BpmInstanceFinish(n.InstanceID, BpmRejected)
					if err != nil {
						fmt.Printf("ERR: peers, ONE_PASS, REJECTED: %+v\n", err)
					}
				}

				// 还有未处理的人，等待他们处理
			}
		}

		// 一票否决
		if n.CalcRule == "ONE_REJECT" {
			if result == BpmRejected {
				err = BpmInstanceFinish(n.InstanceID, BpmRejected)
				if err != nil {
					fmt.Printf("ERR: peers, ONE_REJCT, REJECTED: %+v\n", err)
				}
				for _, v := range waitingList {
					BpmInstanceNodeMark(v.ID, "SKIPPED", "")
				}
			}

			if result == BpmApproved {
				if waitingCount == 0 {
					err = BpmInstanceForward(n.InstanceID)
					if err != nil {
						fmt.Printf("ERR: peers, ONE_REJECT, APPROVED: %+v\n", err)
					}
					return
				}
				// 还有未处理的人，等待他们处理
			}
		}

		if n.CalcRule == "FIRST_ONE" {
			if result == BpmApproved {
				err = BpmInstanceForward(n.InstanceID)
				if err != nil {
					fmt.Printf("ERR: peers, FIRST_ONE, APPROVED: %+v\n", err)
				}
			}

			if result == BpmRejected {
				err = BpmInstanceFinish(n.InstanceID, BpmRejected)
				if err != nil {
					fmt.Printf("ERR: peers, FIRST_ONE, REJECTED: %+v\n", err)
				}
			}

			// 第一人处理后，其他人自动跳过
			for _, v := range waitingList {
				BpmInstanceNodeMark(v.ID, "SKIPPED", "")
			}
		}
	}
}
