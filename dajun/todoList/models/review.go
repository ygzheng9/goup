package models

import "time"

// Review 一次的审批
type Review struct {
	ID            int    `db:"ID" json:"id" form:"id"`
	RefID         int    `db:"REF_ID" json:"ref_id" form:"ref_id"`
	BizType       string `db:"BIZ_TYP" json:"biz_type" form:"biz_type"`
	TxName        string `db:"TX_NME" json:"tx_name" form:"tx_name"`
	FormDate      string `db:"FORM_DTE" json:"form_date" form:"form_date"`
	FormUser      string `db:"FORM_USR" json:"form_user" form:"form_user"`
	FormContent   string `db:"FORM_CTENT" json:"form_content" form:"form_content"`
	ReviewUser    string `db:"REVIEW_USR" json:"review_user" form:"review_user"`
	ReviewDate    string `db:"REVIEW_DTE" json:"review_date" form:"review_date"`
	ReviewContent string `db:"REVIEW_CTENT" json:"review_content" form:"review_content"`
	CurrentStatus string `db:"CURR_STS" json:"current_status" form:"current_status"`
}

const (
	reviewSelect = `select ID, ifnull(REF_ID,'') REF_ID, ifnull(BIZ_TYP,'') BIZ_TYP, ifnull(TX_NME,'') TX_NME, ifnull(FORM_DTE,'') FORM_DTE, ifnull(FORM_USR,'') FORM_USR, ifnull(FORM_CTENT,'') FORM_CTENT, ifnull(REVIEW_USR,'') REVIEW_USR, ifnull(REVIEW_DTE,'') REVIEW_DTE, ifnull(REVIEW_CTENT,'') REVIEW_CTENT, ifnull(CURR_STS,'') CURR_STS	from t_review`
)

// CreateReview 创建一个审批的Form
func CreateReview() Review {
	return Review{ID: -1,
		RefID: -1}
}

// ReviewFindBy 根据条件查找
func ReviewFindBy(cond string) ([]Review, error) {
	cmd := reviewSelect + cond
	// log.Printf("Review search: %s\n", cmd)

	items := []Review{}
	err := db.Select(&items, cmd)
	return items, err
}

// ReviewFindAll 返回所有
func ReviewFindAll() ([]Review, error) {
	return ReviewFindBy("")
}

// ReviewFindByID 根据 ID 加载记录
func ReviewFindByID(id int) (Review, error) {
	item := Review{}
	cmd := reviewSelect + " where ID=? "
	err := db.Get(&item, cmd, id)
	return item, err
}

// ReviewInsert 把当前对象，作为新记录，插入数据库
func ReviewInsert(u Review) (Review, error) {
	item := Review{}
	now := time.Now().Format("2006-01-02 15:04:05")
	u.FormDate = now

	sqlCmd := `INSERT INTO t_review (REF_ID,BIZ_TYP,TX_NME,FORM_DTE,FORM_USR,FORM_CTENT,REVIEW_USR,REVIEW_DTE,REVIEW_CTENT,CURR_STS) VALUES (:REF_ID,:BIZ_TYP,:TX_NME,:FORM_DTE,:FORM_USR,:FORM_CTENT,:REVIEW_USR,:REVIEW_DTE,:REVIEW_CTENT,:CURR_STS)`
	res, err := db.NamedExec(sqlCmd, u)
	if err != nil {
		return item, err
	}

	// 取回插入的记录
	id, err := res.LastInsertId()
	if err != nil {
		return item, err
	}

	return ReviewFindByID(int(id))
}

// ReviewUpdate 将当前对象保存到数据库
func ReviewUpdate(u Review) (Review, error) {
	now := time.Now().Format("2006-01-02 15:04:05")
	if u.CurrentStatus == "DRAFT" ||
		u.CurrentStatus == "SUBMITTED" ||
		u.CurrentStatus == "WITHDRAW" {
		u.FormDate = now
	} else {
		u.ReviewDate = now
	}

	sqlCmd := `update t_review set REF_ID = :REF_ID,
								BIZ_TYP = :BIZ_TYP,
								TX_NME = :TX_NME,
								FORM_DTE = :FORM_DTE,
								FORM_USR = :FORM_USR,
								FORM_CTENT = :FORM_CTENT,
								REVIEW_USR = :REVIEW_USR,
								REVIEW_DTE = :REVIEW_DTE,
								REVIEW_CTENT = :REVIEW_CTENT,
								CURR_STS = :CURR_STS
								where id=:ID`
	_, err := db.NamedExec(sqlCmd, u)
	if err != nil {
		return Review{}, err
	}

	// 从新加载记录
	return ReviewFindByID(u.ID)
}

// ReviewDoAction 仅变更状态，生成日志
func ReviewDoAction(id int, target string, username string) error {
	// 先找到原先的记录
	old, err := ReviewFindByID(id)
	if err != nil {
		return err
	}

	now := time.Now().Format("2006-01-02 15:04:05")

	// 保存变更记录
	keep := ReviewLog{}
	keep.RefType = "Review"
	keep.RefID = id
	keep.FromStatus = old.CurrentStatus
	keep.ToStatus = target
	keep.OpUser = username
	keep.OpDate = now
	err = ReviewLogInsert(keep)
	if err != nil {
		return err
	}

	// 更新数据库
	sqlCmd := `update t_review set
		CURR_STS = :CURR_STS
		where id=:ID`

	param := struct {
		ID        int    `db:"ID"`
		Target    string `db:"CURR_STS"`
		FormDate  string `db:"FORM_DTE"`
		RevewDate string `db:"REVIEW_DTE"`
	}{
		ID:     id,
		Target: target,
	}

	// 更新提交时间
	if target == "SUBMITTED" ||
		target == "DRAFT" {
		sqlCmd = `update t_review set
			CURR_STS = :CURR_STS,
			FORM_DTE = :FORM_DTE
			where id=:ID`

		if target == "SUBMITTED" {
			param.FormDate = now
		} else if target == "DRAFT" {
			param.FormDate = ""
		}
	}

	// 更新审批的时间
	if target == "APPROVED" || target == "REJECTED" || target == "REVIEWING" {
		sqlCmd = `update t_review set
			CURR_STS = :CURR_STS,
			REVIEW_DTE = :REVIEW_DTE
			where id=:ID`

		if target == "REVIEWING" {
			param.RevewDate = ""
		} else {
			param.RevewDate = now
		}
	}

	_, err = db.NamedExec(sqlCmd, param)
	return err
}

// ReviewDelete 根据当前对象的 ID 删除
func ReviewDelete(u Review) error {
	sqlCmd := `delete from t_review where ID=:ID`
	_, err := db.NamedExec(sqlCmd, u)

	return err
}
