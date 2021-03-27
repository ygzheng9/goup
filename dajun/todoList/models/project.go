package models

import (
	"strconv"
)

// Project 项目清单
type Project struct {
	ID           int    `db:"ID" json:"id" form:"id"`
	ProjectGroup string `db:"PROJ_GRP" json:"projectGroup" form:"projectGroup"`
	ProjectCode  string `db:"PROJ_CDE" json:"projectCode" form:"projectCode"`
	Remark       string `db:"RMK" json:"remark" form:"remark"`
	Owner        string `db:"OWNER" json:"owner" form:"owner"`
}

const (
	projectSelect = `select ID, ifnull(PROJ_GRP,'') PROJ_GRP, ifnull(PROJ_CDE,'') PROJ_CDE, ifnull(RMK,'') RMK,
			ifnull(OWNER,'') OWNER
		from t_proj`
)

// ProjectByCode 根据 code 条件查找 一条记录
func ProjectByCode(projCode string) (Project, error) {
	cmd := projectSelect + " where PROJ_CDE = '" + projCode + "'"
	// log.Printf("Project search: %s\n", cmd)

	item := Project{}
	err := db.Get(&item, cmd)
	return item, err
}

// ProjectFindBy 根据条件查找
func ProjectFindBy(cond string) ([]Project, error) {
	cmd := projectSelect + cond
	// log.Printf("Project search: %s\n", cmd)

	items := []Project{}
	err := db.Select(&items, cmd)
	return items, err
}

// ProjectInsert 把当前对象，作为新记录，插入数据库
func ProjectInsert(u Project) error {
	sqlCmd := `INSERT INTO t_proj (PROJ_GRP,PROJ_CDE,RMK, OWNER) VALUES (:PROJ_GRP,:PROJ_CDE,:RMK, :OWNER)`
	_, err := db.NamedExec(sqlCmd, u)
	return err
}

// ProjectDelete 当前对象，按照 ID 从数据库删除
func ProjectDelete(id int) error {
	// 按照 id 删除
	cmd := "delete from t_proj where ID=" + strconv.Itoa(id)
	_, err := db.Exec(cmd)

	return err
}

// ProjectUpdate 当前对象，更新到数据库
func ProjectUpdate(c Project) (Project, error) {
	cmd := `update t_proj set
							PROJ_GRP = :PROJ_GRP,
							PROJ_CDE = :PROJ_CDE,
							OWNER = :OWNER,
							RMK = :RMK
						where ID=:ID`
	_, err := db.NamedExec(cmd, c)
	if err != nil {
		return c, err
	}

	return ProjectFindByID(c.ID)
}

// ProjectFindByID 按照 id 查询
func ProjectFindByID(id int) (Project, error) {
	cmd := projectSelect + ` where ID=?`
	item := Project{}
	err := db.Get(&item, cmd, id)
	return item, err
}
