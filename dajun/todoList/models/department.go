package models

import (
	"bytes"
)

// Department 部门
type Department struct {
	ID       int    `db:"DEPT_ID" form:"id" json:"id" `
	Name     string `db:"DEPT_NME" form:"name" json:"name" `
	Code     string `db:"DEPT_CDE" form:"code" json:"code" `
	Seq      string `db:"DISP_SEQ" form:"seq" json:"seq" `
	Head     string `db:"DEPT_HEAD" form:"head" json:"head" `
	ParentID int    `db:"PARENT_ID" form:"parentID" json:"parentID"`
}

const (
	deptSelect = `select DEPT_ID, ifnull(DEPT_CDE, '') DEPT_CDE, ifnull(DEPT_NME,'') DEPT_NME,
		ifnull(DISP_SEQ, 0) DISP_SEQ, ifnull(DEPT_HEAD, '') DEPT_HEAD, ifnull(PARENT_ID, 0) PARENT_ID from t_dept`
)

// CreateDepartment 返回一个空的 Department 对象
func CreateDepartment() Department {
	return Department{}
}

// FindBy 根据条件查找
func (d Department) FindBy(cond string) ([]Department, error) {
	sqlCmd := deptSelect + cond

	items := []Department{}
	err := db.Select(&items, sqlCmd)

	return items, err
}

// FindAll 返回部门清单
func (d Department) FindAll() ([]Department, error) {
	// 查询清单
	return d.FindBy(" order by DISP_SEQ ")
}

// FindByID 按照 id 查询
func (d Department) FindByID(id int) (Department, error) {
	cmd := deptSelect + ` where DEPT_ID=?`
	item := Department{}
	err := db.Get(&item, cmd, id)
	return item, err
}

// Insert 当前对象，插入到数据库
func (d Department) Insert() error {
	// 根据 struct 中的 DB tag 进行自动 named parameter
	cmd := `INSERT INTO t_dept (DEPT_CDE, DEPT_NME, DISP_SEQ,DEPT_HEAD,	PARENT_ID) VALUES
						(:DEPT_CDE, :DEPT_NME, :DISP_SEQ, :DEPT_HEAD,	:PARENT_ID)`
	_, err := db.NamedExec(cmd, d)
	return err
}

// Update 当前对象，更新到数据库
func (d Department) Update() error {
	// 修改时，不修改上级部门；
	cmd := `update t_dept
						set DEPT_CDE=:DEPT_CDE, DEPT_NME=:DEPT_NME, DISP_SEQ=:DISP_SEQ,
								DEPT_HEAD=:DEPT_HEAD
						where DEPT_ID=:DEPT_ID`
	_, err := db.NamedExec(cmd, d)
	return err
}

// Delete 当前对象，按照 ID 从数据库删除
func (d Department) Delete() error {
	// 按照 id 删除
	cmd := "delete from t_dept where DEPT_ID=:DEPT_ID"
	_, err := db.NamedExec(cmd, d)
	return err
}

// DepartmentBatchAddSubs 为部门批量增加员工
func DepartmentBatchAddSubs(deptName string, userNames []string) error {
	var buf bytes.Buffer
	buf.WriteString("(''")
	for _, n := range userNames {
		buf.WriteString(", '" + n + "'")
	}
	buf.WriteString(")")

	sqlCmd := " update t_users set dept = '" + deptName + "' where NME in " + buf.String()
	_, err := db.Exec(sqlCmd)
	if err != nil {
		return err
	}

	// 刷新缓存
	err = reloadAllUsers()
	return err
}

// DepartmentRemoveSub 移除部门下的员工
func DepartmentRemoveSub(userName string) error {
	sqlCmd := " update t_users set dept = '' where NME = '" + userName + "'"
	_, err := db.Exec(sqlCmd)
	if err != nil {
		return err
	}

	// 刷新缓存
	err = reloadAllUsers()
	return err
}
