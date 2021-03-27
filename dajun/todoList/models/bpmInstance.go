package models

import (
	"bytes"
	"fmt"
	"strconv"
	"time"
)

// BpmInstance 审批实例
type BpmInstance struct {
	ID          int    `db:"ID" json:"id" form:"id"`
	BizType     string `db:"BIZ_TYP" json:"bizType" form:"bizType"`
	StatusField string `db:"STS_FIELD" json:"statusField" form:"statusField"`
	KeyField    string `db:"KEY_FIELD" json:"keyField" form:"keyField"`
	ProcessID   int    `db:"PROCESS_ID" json:"processID" form:"processID"`
	RefID       int    `db:"REF_ID" json:"refID" form:"refID"`
	SubmitUser  string `db:"SUBMIT_USR" json:"submitUser" form:"submitUser"`
	SubmitDate  string `db:"SUBMIT_DTE" json:"submitDate" form:"submitDate"`
	BizDesc     string `db:"BIZ_DESC" json:"bizDesc" form:"bizDesc"`
	CurrentStep int    `db:"CURR_STEP" json:"currentStep" form:"currentStep"`
	IsComplete  string `db:"IS_CMPL" json:"isComplete" form:"isComplete"`
	FinalStatus string `db:"FINAL_STS" json:"finalStatus" form:"finalStatus"`
	RefText1    string `db:"REF_TEXT1" json:"refTxt1" form:"refTxt1"`
	RefText2    string `db:"REF_TEXT2" json:"refTxt2" form:"refTxt2"`
	RefText3    string `db:"REF_TEXT3" json:"refTxt3" form:"refTxt3"`
	RefAmt1     string `db:"REF_AMOUNT1" json:"refAmt1" form:"refAmt1"`
	RefAmt2     string `db:"REF_AMOUNT2" json:"refAmt2" form:"refAmt2"`
	UpdateDate  string `db:"REC_UPD_DTE" json:"updateDate" form:"updateDate"`
}

const (
	bpmInstanceSelect = `select ID, PROCESS_ID, CURR_STEP, ifnull(BIZ_TYP,'') BIZ_TYP, ifnull(STS_FIELD,'') STS_FIELD, ifnull(KEY_FIELD,'') KEY_FIELD, ifnull(REF_ID,'') REF_ID, ifnull(SUBMIT_USR,'') SUBMIT_USR,ifnull(SUBMIT_DTE,'') SUBMIT_DTE, ifnull(IS_CMPL,'') IS_CMPL, ifnull(FINAL_STS,'') FINAL_STS, ifnull(REF_TEXT1,'') REF_TEXT1, ifnull(REF_TEXT2,'') REF_TEXT2, ifnull(REF_TEXT3,'') REF_TEXT3, ifnull(REF_AMOUNT1,'') REF_AMOUNT1, ifnull(REF_AMOUNT2,'') REF_AMOUNT2
		from t_bpm_ins`
)

// BpmInstanceFindBy 根据条件查找
func BpmInstanceFindBy(cond string) ([]BpmInstance, error) {
	sqlCmd := bpmProcessSelect + cond

	// fmt.Printf("BpmInstanceFindBy: %s\n", sqlCmd)

	items := []BpmInstance{}
	err := db.Select(&items, sqlCmd)

	return items, err
}

// BpmInstanceFindAll 返回部门清单
func BpmInstanceFindAll() ([]BpmInstance, error) {
	// 查询清单
	return BpmInstanceFindBy(" ")
}

// BpmInstanceFindByID 按照 id 查询
func BpmInstanceFindByID(id int) (BpmInstance, error) {
	cmd := bpmInstanceSelect + ` where ID=?`
	item := BpmInstance{}
	err := db.Get(&item, cmd, id)
	return item, err
}

// BpmInstanceInsert 当前对象，插入到数据库
func BpmInstanceInsert(c BpmInstance) (BpmInstance, error) {
	// 默认都有效
	c.IsComplete = "N"
	c.SubmitDate = time.Now().Format("2006-01-02 15:04:05")

	// 根据 struct 中的 DB tag 进行自动 named parameter
	cmd := `INSERT INTO t_bpm_ins (PROCESS_ID,BIZ_TYP,STS_FIELD,KEY_FIELD, REF_ID,SUBMIT_USR,SUBMIT_DTE,BIZ_DESC,CURR_STEP,IS_CMPL,FINAL_STS,REF_TEXT1,REF_TEXT2,REF_TEXT3,REF_AMOUNT1,REF_AMOUNT2) VALUES
						(:PROCESS_ID,:BIZ_TYP,:STS_FIELD,:KEY_FIELD,:REF_ID,:SUBMIT_USR,:SUBMIT_DTE,:BIZ_DESC,:CURR_STEP,:IS_CMPL,:FINAL_STS,:REF_TEXT1,:REF_TEXT2,:REF_TEXT3,:REF_AMOUNT1,:REF_AMOUNT2)`
	res, err := db.NamedExec(cmd, c)
	if err != nil {
		return BpmInstance{}, err
	}
	fmt.Println("complete: BpmInstanceInsert.")

	// 取回插入的记录
	id, err := res.LastInsertId()
	if err != nil {
		return BpmInstance{}, err
	}

	return BpmInstanceFindByID(int(id))
}

// BpmInstanceUpdate 当前对象，更新到数据库
func BpmInstanceUpdate(c BpmInstance) (BpmInstance, error) {
	c.UpdateDate = time.Now().Format("2006-01-02 15:04:05")

	cmd := `update t_bpm_ins set
							PROCESS_ID = :PROCESS_ID,
							BIZ_TYP = :BIZ_TYP,
							STS_FIELD = :STS_FIELD,
							KEY_FIELD = :KEY_FIELD,
							REF_ID = :REF_ID,
							SUBMIT_USR = :SUBMIT_USR,
							SUBMIT_DTE = :SUBMIT_DTE,
							BIZ_DESC = :BIZ_DESC,
							CURR_STEP = :CURR_STEP,
							IS_CMPL = :IS_CMPL,
							FINAL_STS = :FINAL_STS,
							REF_TEXT1 = :REF_TEXT1,
							REF_TEXT2 = :REF_TEXT2,
							REF_TEXT3 = :REF_TEXT3,
							REF_AMOUNT1 = :REF_AMOUNT1,
							REF_AMOUNT2 = :REF_AMOUNT2,
							REC_UPD_DTE = :REC_UPD_DTE
						where ID=:ID`
	_, err := db.NamedExec(cmd, c)
	if err != nil {
		return c, err
	}

	return BpmInstanceFindByID(c.ID)
}

// BpmInstanceSetStep 设置当前步骤
func BpmInstanceSetStep(id, step int) error {
	now := time.Now().Format("2006-01-02 15:04:05")

	updCmd := " update t_bpm_ins set CURR_STEP=" + strconv.Itoa(step) + ", REC_UPD_DTE='" + now + "' where id= " + strconv.Itoa(id)

	_, err := db.Exec(updCmd)
	return err
}

// BpmInstanceDelete 当前对象，按照 ID 从数据库删除
func BpmInstanceDelete(id int) error {
	// 按照 id 删除
	strID := strconv.Itoa(id)

	cmd := "delete from t_bpm_ins_node where INS_ID = " + strID
	_, err := db.Exec(cmd)

	cmd = "delete from t_bpm_ins where ID=" + strID
	_, err = db.Exec(cmd)
	return err
}

// BpmInstanceStart 启动一个审批流
func BpmInstanceStart(bizID int, bpmType, submitUser, bizDesc string) bool {
	// 根据业务类型，取得流程配置，按优先级排序
	cond := " where BIZ_TYP = '" + bpmType + "' and VALID_IND = 'Y' order by PRIORITY_NUM asc"
	bpmProcessList, err := BpmProcessFindBy(cond)
	if err != nil {
		fmt.Printf("Err: BpmProcessFindBy: %+v\n", err)
		return false
	}

	fmt.Printf("%s\n", cond)

	// 如果没有该业务类型对应的流程配置，那么不能启动审批流
	if len(bpmProcessList) == 0 {
		fmt.Printf("Err: 未配置审批流: %s\n", bpmType)
		return false
	}

	// 一个业务，可以配置多个规则，按优先级，进行匹配；只匹配一个；
	for _, p := range bpmProcessList {
		var buf bytes.Buffer
		buf.WriteString(" select count(1)  from " + p.BizType + " a where 1 = 1 ")
		if len(p.BizRule) != 0 {
			buf.WriteString(" and " + p.BizRule)
		}
		buf.WriteString(" and a." + p.KeyField + "= " + strconv.Itoa(bizID))
		bufStr := buf.String()

		fmt.Printf("%s\n", bufStr)

		var count int
		err = db.Get(&count, bufStr)
		if err != nil {
			fmt.Printf("check rule: %s \n %+v\n", bufStr, err)
			return false
		}

		// 审批流触发规则已匹配
		if count > 0 {
			param := BpmInstance{
				ProcessID:   p.ID,
				BizType:     p.BizType,
				StatusField: p.StatusField,
				KeyField:    p.KeyField,
				RefID:       bizID,
				SubmitUser:  submitUser,
				BizDesc:     bizDesc,
				CurrentStep: -1,
			}

			ins, err := BpmInstanceInsert(param)
			if err != nil {
				fmt.Printf("Err: create instance: %+v\n", err)
				return false
			}

			fmt.Printf("已生成审批实例： %d\n", ins.ID)

			err = BpmInstanceForward(ins.ID)
			if err != nil {
				fmt.Printf("Err: 1st forward: %+v\n", err)
				return false
			}
			return true
		}
	}

	return false
}

// BpmInstanceFinish 审批实例完成
func BpmInstanceFinish(id int, status string) error {
	now := time.Now().Format("2006-01-02 15:04:05")
	// 设置审批实例状态
	var buf bytes.Buffer
	buf.WriteString(" update t_bpm_ins set IS_CMPL='Y' ")
	buf.WriteString(" ,FINAL_STS='" + status + "' ,REC_UPD_DTE='" + now + "' ")
	buf.WriteString(" where ID = " + strconv.Itoa(id))

	fmt.Printf("finished: %s\n", buf.String())

	_, err := db.Exec(buf.String())
	if err != nil {
		fmt.Printf("Err: BpmInstanceFinish, %+v\n", err)
		return err
	}

	// 更新业务单据状态
	ins, err := BpmInstanceFindByID(id)
	if err != nil {
		fmt.Printf("Err: BpmInstanceFindByID %+v", err)
		return err
	}

	buf.Reset()
	buf.WriteString(" update " + ins.BizType + " set " + ins.StatusField + " = '" + status + "' ")
	buf.WriteString(" where " + ins.KeyField + " = " + strconv.Itoa(ins.RefID))
	_, err = db.Exec(buf.String())
	if err != nil {
		fmt.Printf("Err: writeback %+v", err)
		return err
	}

	return nil
}

// BpmInstanceForward 工作流前进一步
func BpmInstanceForward(id int) error {
	// 加载实例
	ins, err := BpmInstanceFindByID(id)
	if err != nil {
		fmt.Printf("ins id: %d\n", id)
		return err
	}
	// 保存之前的步骤，一次 forward，有可能跳过几步，需要记下最初的状态
	oldStep := ins.CurrentStep
	fmt.Printf("ins: %+v\n", ins)

	// 根据实例当前节点，取出后续节点，按顺序排
	var buf bytes.Buffer
	buf.WriteString(" where PROCESS_ID = " + strconv.Itoa(ins.ProcessID))
	buf.WriteString(" and STEP_NUM > " + strconv.Itoa(ins.CurrentStep))
	buf.WriteString(" order by STEP_NUM asc; ")
	nodes, err := BpmProcessNodeFindBy(buf.String())
	fmt.Printf("find node: %s\n", buf.String())

	if err != nil {
		fmt.Printf("Err: %+v\n", err)
		return err
	}

	// 无后续节点，则标记实例已完成
	if len(nodes) == 0 {
		fmt.Println("no nodes")
		err = BpmInstanceFinish(id, "APPROVED")
		if err != nil {
			return err
		}
	}

	// 有后续节点，每个节点可能对应多个用户（比如：一票否决）
	type approverT struct {
		NodeUser string `db:"ndUser" json:"ndUser"`
	}
	nodeUsers := []approverT{}

	// 第一个即为下一个，查询时已经排序
	// 是否找到了一个 审批人
	gotApprover := false
	for _, nt := range nodes {
		nextStep := nt
		fmt.Printf("next node: %+v\n", nextStep)
		// 设置审批示例的当前步骤，为该节点步骤
		err = BpmInstanceSetStep(ins.ID, nextStep.StepNumber)
		if err != nil {
			fmt.Printf("Err: setstep %+v\n", err)
			return err
		}

		// 根据节点规则，找到该环节的所有审批人
		// 节点规则中的 a 必须是业务主表，其中必须有 ID 字段，和审批实例中的 RefID 对应
		buf.Reset()
		buf.WriteString(nextStep.BizRule)
		buf.WriteString(" where a.ID = " + strconv.Itoa(ins.RefID))
		err = db.Select(&nodeUsers, buf.String())
		if err != nil {
			fmt.Printf("Err: node rule: %s\n", buf.String())
			return err
		}

		if len(nodeUsers) == 0 {
			// 流程中有节点，但是根据节点的规则，没有找到对应的审批人
			// 该节点插入审批记录，然后进入下一个节点
			fmt.Printf("not applicable: %+v\n", nextStep)
			user := BpmInstanceNode{
				InstanceID:    ins.ID,
				StepNumber:    nextStep.StepNumber,
				CurrentUser:   "不适用",
				CalcRule:      nextStep.CalcRule,
				CurrentResult: BpmNonApplicable,
			}

			_, err = BpmInstanceNodeInsert(user)
			if err != nil {
				return err
			}
			// 该审批环节，没有满足条件的审批人，自动跳过
			continue
		}

		// 找到了审批人
		gotApprover = true
		// 该节点每一个审批人，都要在审批历史表中，创建一条记录
		for _, u := range nodeUsers {
			user := BpmInstanceNode{
				InstanceID:  ins.ID,
				StepNumber:  nextStep.StepNumber,
				CurrentUser: u.NodeUser,
				CalcRule:    nextStep.CalcRule,
			}

			nd, err := BpmInstanceNodeInsert(user)
			if err != nil {
				return err
			}

			// 如果有代理人，那么把代理人插入审批历史
			dsql := " where USR_NME='" + u.NodeUser + "' and VALID_IND ='Y'"
			delegations, err := BpmDelegationFindBy(dsql)
			fmt.Printf("Delegation: %s\n", dsql)
			if err != nil {
				fmt.Printf("Err: find delegation: %+v\n", err)
				return err
			}

			if len(delegations) > 0 {
				// 只能代理给一个人，并且代理授权不能传递
				d := delegations[0]
				fmt.Printf("got delegation: %+v\n", d)

				_, err = BpmInstanceNodeDelegate(nd, d.DelegateTo)
				if err != nil {
					return err
				}
			}
		}

		// 已经找到了审批人，一次只前进一步
		if gotApprover {
			break
		}
	}

	// 没有找到任何一个审批人，审批流实例标记完成
	if !gotApprover {
		if oldStep > 0 {
			fmt.Println("完成审批")
			err = BpmInstanceFinish(id, "APPROVED")
		} else {
			// 初始步骤小于零，表示之前没有任何一个人批过，要驳回
			fmt.Println("根据配置找不到任何审批人")
			err = BpmInstanceFinish(id, "REJECTED")
		}

		if err != nil {
			return err
		}
	}

	return nil
}
