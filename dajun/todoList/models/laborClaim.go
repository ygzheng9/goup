package models

import (
	"fmt"
	"time"
)

// LaborClaim 记录
type LaborClaim struct {
	ID          int    `db:"ID" json:"id" form:"id"`
	UserName    string `db:"USR_NME" json:"userName" form:"userName"`
	BizDate     string `db:"BIZ_DTE" json:"bizDate" form:"bizDate"`
	ProjCode    string `db:"PROJ_CDE" json:"projCode" form:"projCode"`
	HourCount   int    `db:"HOUR_CNT" json:"hourCount" form:"hourCount"`
	UpdateUser  string `db:"UPD_USR" json:"updateUser" form:"updateUser"`
	UpdateDate  string `db:"UPD_DTE" json:"updateDate" form:"updateDate"`
	Remark      string `db:"RMK" json:"remark" form:"remark"`
	ProjGroup   string `db:"PROJ_GRP" json:"projGroup" form:"projGroup"`
	ProjOwner   string `db:"PROJ_OWNER" json:"projOwner" form:"projOwner"`
	UserDept    string `db:"USR_DEPT" json:"userDept" form:"userDept"`
	UserTitle   string `db:"USR_TITLE" json:"userTitle" form:"userTitle"`
	UserSvcLine string `db:"USR_SVC_LINE" json:"userSvcLine" form:"userSvcLine"`
}

const (
	laborClaimSelect = `select ID,ifnull(USR_NME,'') USR_NME, ifnull(PROJ_CDE,'') PROJ_CDE, ifnull(BIZ_DTE,'') BIZ_DTE, ifnull(HOUR_CNT,'') HOUR_CNT, ifnull(UPD_DTE,'') UPD_DTE, ifnull(RMK,'') RMK, ifnull(PROJ_GRP,'') PROJ_GRP, ifnull(USR_DEPT,'') USR_DEPT,
			ifnull(USR_TITLE,'') USR_TITLE, ifnull(USR_SVC_LINE,'') USR_SVC_LINE, ifnull(PROJ_OWNER,'') PROJ_OWNER
		from t_labor_claim`
)

// LaborClaimFindBy 根据条件查找
func LaborClaimFindBy(cond string) ([]LaborClaim, error) {
	cmd := laborClaimSelect + cond
	// log.Printf("LaborClaim search: %s\n", cmd)

	items := []LaborClaim{}
	err := db.Select(&items, cmd)
	return items, err
}

// LaborClaimFindByWk 查找用户一周内的记录
func LaborClaimFindByWk(userName, monday, sunday string) ([]LaborClaim, error) {
	quote := func(str string) string {
		return "'" + str + "'"
	}
	cond := " where USR_NME=" + quote(userName) + " and BIZ_DTE >= " + quote(monday) + " and BIZ_DTE <= " + quote(sunday)

	return LaborClaimFindBy(cond)
}

// LaborClaimFindProj 查找用户的项目清单
func LaborClaimFindProj(userName, from, to string) ([]string, error) {
	quote := func(str string) string {
		return "'" + str + "'"
	}
	sqlCmd := " select distinct PROJ_CDE from t_labor_claim where USR_NME=" + quote(userName) + " and BIZ_DTE >= " + quote(from) + " and BIZ_DTE <= " + quote(to) + " order by PROJ_CDE"

	items := []string{}
	err := db.Select(&items, sqlCmd)
	return items, err
}

// LaborClaimInsert 把当前对象，作为新记录，插入数据库
func LaborClaimInsert(u LaborClaim) error {
	var err error
	var sqlCmd string

	now := time.Now().Format("2006-01-02 15:04:05")
	u.UpdateDate = now

	// 保存当前项目信息
	projInfo, err := ProjectByCode(u.ProjCode)
	if err != nil {
		fmt.Printf("Proj: %s\n", u.ProjCode)
		fmt.Printf("err: %+v\n", err)
		// return err
	}
	u.ProjGroup = projInfo.ProjectGroup
	u.ProjOwner = projInfo.Owner

	// 保存用户当前时刻的信息
	userInfo := struct {
		Dept    string `db:"DEPT"`
		Title   string `db:"TITLE"`
		SvcLine string `db:"SVC_LINE"`
	}{}

	sqlCmd = "select DEPT, TITLE, SVC_LINE from t_users where NME = '" + u.UserName + "'"
	err = db.Get(&userInfo, sqlCmd)
	if err != nil {
		fmt.Printf("userInfo err: %s\n", sqlCmd)
		fmt.Printf("err: %+v\n", err)
		// return err
	}

	u.UserDept = userInfo.Dept
	u.UserTitle = userInfo.Title
	u.UserSvcLine = userInfo.SvcLine
	fmt.Printf("insert claim: %+v\n", u)

	sqlCmd = `INSERT INTO t_labor_claim (USR_NME,BIZ_DTE,PROJ_CDE,HOUR_CNT,UPD_USR,UPD_DTE,RMK, PROJ_GRP, USR_DEPT, USR_TITLE, USR_SVC_LINE, PROJ_OWNER) VALUES (:USR_NME,:BIZ_DTE,:PROJ_CDE,:HOUR_CNT,:UPD_USR,:UPD_DTE,:RMK, :PROJ_GRP, :USR_DEPT, :USR_TITLE, :USR_SVC_LINE, :PROJ_OWNER)`
	_, err = db.NamedExec(sqlCmd, u)
	return err
}

// LaborClaimDeleteWk 删除该用户，该周的数据
func LaborClaimDeleteWk(userName, monday, sunday string) error {
	quote := func(str string) string {
		return "'" + str + "'"
	}

	sqlCmd := "delete from t_labor_claim where USR_NME=" + quote(userName) + " and BIZ_DTE >= " + quote(monday) + " and BIZ_DTE <= " + quote(sunday)
	_, err := db.Exec(sqlCmd)

	return err
}

// LaborClaimSearchByPeriod 取得期间内记录
func LaborClaimSearchByPeriod(userName, start, end string) ([]LaborClaim, error) {
	quote := func(str string) string {
		return "'" + str + "'"
	}
	cond := " where BIZ_DTE >= " + quote(start) + " and BIZ_DTE <= " + quote(end)

	if len(userName) > 0 {
		cond = cond + " and USR_NME = " + quote(userName)
	}

	cond = cond + " order by BIZ_DTE, USR_NME, PROJ_CDE"

	return LaborClaimFindBy(cond)
}
