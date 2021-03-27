package models

import (
	"strings"
	"time"
)

// Feedback 每一笔互动信息
type Feedback struct {
	ID         int    `db:"ID" json:"id" form:"id"`
	RefType    string `db:"REF_TYP" json:"ref_type" form:"ref_type"`
	RefID      int    `db:"REF_ID" json:"ref_id" form:"ref_id"`
	RefTitle   string `db:"REF_TITLE" json:"ref_title" form:"ref_title"`
	Comment    string `db:"CMMT" json:"comment" form:"comment"`
	CreateUser string `db:"CRE_USR" json:"create_user" form:"create_user"`
	CreateDate string `db:"CRE_DTE" json:"create_date" form:"create_date"`
}

const (
	feedbackSelect = `select
											ID, ifnull(REF_TYP,'') REF_TYP, ifnull(REF_ID,'') REF_ID, ifnull(REF_TITLE,'') REF_TITLE, ifnull(CMMT,'') CMMT, ifnull(CRE_USR,'') CRE_USR, ifnull(CRE_DTE,'') CRE_DTE
										from t_feedback `
)

// CreateFeedback 返回一个空对象
func CreateFeedback() Feedback {
	return Feedback{}
}

// FindBy 返回符合条件的记录
func (t Feedback) FindBy(cond string) ([]Feedback, error) {
	cmd := feedbackSelect

	if len(strings.TrimSpace(cond)) > 0 {
		cmd = cmd + " " + cond
	}

	items := []Feedback{}
	err := db.Select(&items, cmd)
	return items, err
}

// FindAll 返回所有记录
func (t Feedback) FindAll() ([]Feedback, error) {
	return t.FindBy("")
}

// FindByRef 根据 ref 查找，参数是当前对象的属性
func (t Feedback) FindByRef() ([]Feedback, error) {
	// 查询的返回结果，默认是空列表
	items := []Feedback{}
	// 查询
	cmd := feedbackSelect +
		` where REF_TYP = :REF_TYP and REF_ID = :REF_ID `

	nstmt, err := db.PrepareNamed(cmd)
	if err != nil {
		return items, err
	}

	// 通过当前对象的属性，传递查询条件
	err = nstmt.Select(&items, t)
	return items, err

}

// FindByID 根据 ID 查找，返回一个新对象
func (t Feedback) FindByID(id int) (Feedback, error) {
	cmd := feedbackSelect +
		` where ID=?`

	item := Feedback{}
	err := db.Get(&item, cmd, id)

	return item, err
}

// Insert 当前对象，插入到数据库中
func (t Feedback) Insert() error {
	// 创建时间为当前时刻
	t.CreateDate = time.Now().Format("2006-01-02 15:04:05")

	// 根据 struct 中的 DB tag 进行自动 named parameter
	insertCmd := `INSERT INTO t_feedback
									(REF_TYP, REF_ID, REF_TITLE, CMMT, CRE_USR, CRE_DTE)
								VALUES
									(:REF_TYP, :REF_ID, :REF_TITLE, :CMMT, :CRE_USR, :CRE_DTE)`
	_, err := db.NamedExec(insertCmd, t)

	return err
}

// Update 当前对象更新回数据库
func (t Feedback) Update() error {
	// 根据 struct 中的 tag 进行自动 named parameter
	sqlCmd := `update t_feedback set
								REF_TYP = :REF_TYP,
								REF_ID = :REF_ID,
								REF_TITLE = :REF_TITLE,
								CMMT = :CMMT,
								CRE_USR = :CRE_USR,
								CRE_DTE = :CRE_DTE
							where ID=:ID`
	_, err := db.NamedExec(sqlCmd, t)

	return err
}

// Delete 从数据库删除当前对象
func (t Feedback) Delete() error {
	// 按照 id 删除
	cmd := "delete from t_feedback where ID=?"
	_, err := db.Exec(cmd, t.ID)
	return err
}
