package models

import (
	"strconv"
)

// SoftInst 软件安装申请记录
type SoftInst struct {
	ID            int    `db:"ID" form:"id" json:"id" `
	ApplyNum      string `db:"APPLY_NUM" form:"apply_num" json:"apply_num"`
	ApplyUser     string `db:"APPLY_USR" form:"apply_user" json:"apply_user" `
	Department    string `db:"DEPT" form:"department" json:"department" `
	Position      string `db:"PSTN" form:"position" json:"position" `
	ApplyDate     string `db:"APPLY_DTE" form:"apply_date" json:"apply_date" `
	SoftName      string `db:"SOFT_NME" form:"soft_name" json:"soft_name" `
	SoftVer       string `db:"SOFT_VER" form:"soft_ver" json:"soft_ver" `
	SoftType      string `db:"SOFT_TYP" form:"soft_type" json:"soft_type" `
	ApplyReason   string `db:"APPLY_RSON" form:"apply_reason" json:"apply_reason" `
	DepartmentMgr string `db:"DEPT_MGR" form:"dept_manager" json:"dept_manager" `
	Remark        string `db:"RMK" form:"remark" json:"remark" `
	CreateUser    string `db:"CRE_USR" form:"create_user" json:"create_user" `
	CreateDate    string `db:"CRE_DTE" form:"create_date" json:"create_date" `
}

const (
	softInstSelect = `select ID, ifnull(APPLY_NUM, '') APPLY_NUM, ifnull(APPLY_USR,'') APPLY_USR, ifnull(DEPT,'') DEPT, ifnull(PSTN,'') PSTN,
													ifnull(APPLY_DTE,'') APPLY_DTE, ifnull(SOFT_NME,'') SOFT_NME, ifnull(SOFT_VER,'') SOFT_VER, ifnull(SOFT_TYP,'') SOFT_TYP, ifnull(APPLY_RSON,'') APPLY_RSON,
													ifnull(DEPT_MGR,'') DEPT_MGR, ifnull(RMK,'') RMK, ifnull(CRE_USR,'') CRE_USR, ifnull(CRE_DTE,'') CRE_DTE
										from t_soft_inst`
)

// CreateSoftInst 返回一个空对象
func CreateSoftInst() SoftInst {
	return SoftInst{}
}

// FindBy 根据条件查询
func (t SoftInst) FindBy(cond string) ([]SoftInst, error) {
	cmd := softInstSelect + " " + cond
	items := []SoftInst{}
	err := db.Select(&items, cmd)
	return items, err
}

// FindAll 返回全部
func (t SoftInst) FindAll() ([]SoftInst, error) {
	return t.FindBy("")
}

// FindByID 根据ID查找
func (t SoftInst) FindByID(id int) (SoftInst, error) {
	// 按照 id 查询
	cmd := softInstSelect +
		` where ID=?`

	item := SoftInst{}
	err := db.Get(&item, cmd, id)

	return SoftInst{}, err
}

// SoftInstFindByID 根据 ID 加载
func SoftInstFindByID(id int) (SoftInst, error) {
	// 按照 id 查询
	cmd := softInstSelect + " where ID=" + strconv.Itoa(id)

	item := SoftInst{}
	err := db.Get(&item, cmd)

	return item, err
}

// Insert 把当前对象插入到数据库
func (t SoftInst) Insert() error {
	// 根据 struct 中的 DB tag 进行自动 named parameter
	insertCmd := `INSERT INTO t_soft_inst
									(APPLY_NUM, APPLY_USR, DEPT, DEPT_MGR, PSTN, APPLY_DTE, SOFT_NME, SOFT_VER, SOFT_TYP, APPLY_RSON, RMK, CRE_USR, CRE_DTE)
								VALUES
									(:APPLY_NUM, :APPLY_USR, :DEPT, :DEPT_MGR, :PSTN, :APPLY_DTE, :SOFT_NME, :SOFT_VER, :SOFT_TYP, :APPLY_RSON, :RMK, :CRE_USR, :CRE_DTE)
							`
	_, err := db.NamedExec(insertCmd, t)
	return err
}

// Update 更新数据库，当前对象
func (t SoftInst) Update() error {
	// 根据 struct 中的 tag 进行自动 named parameter
	sqlCmd := `update t_soft_inst set
								APPLY_NUM=:APPLY_NUM,
								APPLY_USR=:APPLY_USR,
								DEPT=:DEPT,
								DEPT_MGR=:DEPT_MGR,
								PSTN=:PSTN,
								APPLY_DTE=:APPLY_DTE,
								SOFT_NME=:SOFT_NME,
								SOFT_VER=:SOFT_VER,
								SOFT_TYP=:SOFT_TYP,
								APPLY_RSON=:APPLY_RSON,
								RMK=:RMK,
								CRE_USR=:CRE_USR,
								CRE_DTE=:CRE_DTE
							where ID=:ID`
	_, err := db.NamedExec(sqlCmd, t)

	return err
}

// Delete 从数据库中删除当前对象
func (t SoftInst) Delete() error {
	// 按照 id 删除
	cmd := "delete from t_soft_inst where ID=?"
	_, err := db.Exec(cmd, t.ID)

	return err
}

// StartBPM 启动工作流
func (t SoftInst) StartBPM() bool {
	// 业务类型，业务ID
	bpmType := "t_soft_inst"
	bizID := t.ID
	submitUser := t.ApplyUser
	bizDesc := "申请安装 " + t.SoftType + " " + t.SoftName + " 版本 " + t.SoftVer + " " + t.Remark

	return BpmInstanceStart(bizID, bpmType, submitUser, bizDesc)
}
