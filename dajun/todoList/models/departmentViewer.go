package models

// DeptViewer 部门查看权限的用户
type DeptViewer struct {
	ID         int    `db:"ID" form:"id" json:"id" `
	Department string `db:"DEPT" form:"department" json:"department" `
	User       string `db:"USR" form:"user" json:"user" `
}

// DeptViewerFindAll 取得所有查看权限
func DeptViewerFindAll() ([]DeptViewer, error) {
	sqlCmd := "select ID, DEPT, USR from t_dept_viewer"

	items := []DeptViewer{}
	err := db.Select(&items, sqlCmd)

	return items, err
}

// DeptViewerBatchAdd 批量建立关系
func DeptViewerBatchAdd(dept string, users []string) error {
	var sqlCmd string
	var err error
	for _, user := range users {
		sqlCmd = "insert into t_dept_viewer (DEPT, USR) values ('" + dept + "', '" + user + "')"
		_, err = db.Exec(sqlCmd)
		if err != nil {
			return err
		}
	}

	return nil
}

// DeptViewerRemove 移除关联关系
func DeptViewerRemove(dept, user string) error {
	var sqlCmd string
	var err error
	sqlCmd = "delete from t_dept_viewer where DEPT = '" + dept + "' and USR = '" + user + "'"
	_, err = db.Exec(sqlCmd)

	return err
}
