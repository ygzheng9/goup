package models

import "strconv"

// ActivityLog 每天的工作记录
type ActivityLog struct {
	ID      int    `db:"ID" json:"id" form:"id"`
	RefType string `db:"REF_TYP" json:"ref_type" form:"ref_type"`
	RefID   int    `db:"REF_ID" json:"ref_id" form:"ref_id"`
	Msg1    string `db:"MSG1" json:"msg1" form:"msg1"`
	Msg2    string `db:"MSG2" json:"msg2" form:"msg2"`
	Msg3    string `db:"MSG3" json:"msg3" form:"msg3"`
	OpUser  string `db:"OPR_USR" json:"op_user" form:"op_user"`
	OpDate  string `db:"OPR_DTE" json:"op_date" form:"op_date"`
}

const (
	activityLogSelect = `select ID, ifnull(REF_TYP,'') REF_TYP, ifnull(REF_ID,'') REF_ID, ifnull(MSG1,'') MSG1, ifnull(MSG2,'') MSG2, ifnull(MSG3,'') MSG3, ifnull(OPR_USR,'') OPR_USR, ifnull(OPR_DTE,'') OPR_DTE	from t_actvy_log`
)

// ActivityLogFindBy 根据条件查找
func ActivityLogFindBy(cond string) ([]ActivityLog, error) {
	cmd := activityLogSelect + cond
	// log.Printf("ActivityLog search: %s\n", cmd)

	items := []ActivityLog{}
	err := db.Select(&items, cmd)
	return items, err
}

// ActivityLogFindByRefID 根据 refID 查找审批记录
func ActivityLogFindByRefID(refID int, refType string) ([]ActivityLog, error) {
	cond := " where REF_ID = " + strconv.Itoa(refID) + " and REF_TYP = '" + refType + "' order by OPR_DTE desc"
	return ActivityLogFindBy(cond)
}

// ActivityLogInsert 把当前对象，作为新记录，插入数据库
func ActivityLogInsert(u ActivityLog) error {
	sqlCmd := `INSERT INTO t_actvy_log (REF_TYP,REF_ID,MSG1,MSG2,MSG3,OPR_USR,OPR_DTE) VALUES (:REF_TYP,:REF_ID,:MSG1,:MSG2,:MSG3,:OPR_USR,:OPR_DTE)`
	_, err := db.NamedExec(sqlCmd, u)

	return err
}
