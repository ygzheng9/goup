package models

import (
	"strconv"
)

// ReviewLog 一次的审批
type ReviewLog struct {
	ID         int    `db:"ID" json:"id" form:"id"`
	RefID      int    `db:"REF_ID" json:"ref_id" form:"ref_id"`
	RefType    string `db:"REF_TYP" json:"ref_type" form:"ref_type"`
	FromStatus string `db:"FM_STS" json:"from_status" form:"from_status"`
	ToStatus   string `db:"TO_STS" json:"to_status" form:"to_status"`
	OpUser     string `db:"OPR_USR" json:"op_user" form:"op_user"`
	OpDate     string `db:"OPR_DTE" json:"op_date" form:"op_date"`
}

const (
	reviewLogSelect = `select ID, ifnull(REF_ID,'') REF_ID, ifnull(REF_TYP,'') REF_TYP, ifnull(FM_STS,'') FM_STS, ifnull(TO_STS,'') TO_STS, ifnull(OPR_USR,'') OPR_USR, ifnull(OPR_DTE,'') OPR_DTE	from t_review_log`
)

// CreateReviewLog 创建一个审批的Form
func CreateReviewLog() ReviewLog {
	return ReviewLog{ID: -1,
		RefID: -1}
}

// ReviewLogFindBy 根据条件查找
func ReviewLogFindBy(cond string) ([]ReviewLog, error) {
	cmd := reviewLogSelect + cond
	// log.Printf("ReviewLog search: %s\n", cmd)

	items := []ReviewLog{}
	err := db.Select(&items, cmd)
	return items, err
}

// ReviewLogFindAll 返回所有
func ReviewLogFindAll() ([]ReviewLog, error) {
	return ReviewLogFindBy("")
}

// ReviewLogFindByID 根据 ID 加载记录
func ReviewLogFindByID(id int) (ReviewLog, error) {
	item := ReviewLog{}
	cmd := reviewLogSelect + " where ID=? "
	err := db.Get(&item, cmd, id)
	return item, err
}

// ReviewLogFindByRefID 根据 refID 查找审批记录
func ReviewLogFindByRefID(refID int) ([]ReviewLog, error) {
	cond := " where REF_ID = " + strconv.Itoa(refID) + " order by OPR_DTE desc"
	return ReviewLogFindBy(cond)
}

// ReviewLogInsert 把当前对象，作为新记录，插入数据库
func ReviewLogInsert(u ReviewLog) error {
	sqlCmd := `INSERT INTO t_review_log (REF_ID,REF_TYP,FM_STS,TO_STS,OPR_USR,OPR_DTE) VALUES (:REF_ID,:REF_TYP,:FM_STS,:TO_STS,:OPR_USR,:OPR_DTE)`
	_, err := db.NamedExec(sqlCmd, u)

	return err
}

// ReviewLogUpdate 将当前对象保存到数据库
func ReviewLogUpdate(u ReviewLog) error {
	sqlCmd := `update t_review_log set REF_ID = :REF_ID,
								REF_TYP = :REF_TYP,
								FM_STS = :FM_STS,
								TO_STS = :TO_STS,
								OPR_USR = :OPR_USR,
								OPR_DTE = :OPR_DTE
								where id=:ID`
	_, err := db.NamedExec(sqlCmd, u)

	return err
}

// ReviewLogDelete 根据当前对象的 ID 删除
func ReviewLogDelete(u ReviewLog) error {
	sqlCmd := `delete from t_review_log where ID=:ID`
	_, err := db.NamedExec(sqlCmd, u)

	return err
}
