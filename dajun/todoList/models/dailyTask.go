package models

import "time"

// DailyTask 每天的工作记录
type DailyTask struct {
	ID         int    `db:"ID" json:"id" form:"id"`
	UserName   string `db:"USR_NME" json:"user_name" form:"user_name"`
	BizDate    string `db:"BIZ_DTE" json:"biz_date" form:"biz_date"`
	WorkRemark string `db:"WORK_RMK" json:"work_remark" form:"work_remark"`
	UpdateDate string `db:"UPD_DTE" json:"update_date" form:"update_date"`
}

const (
	dailyTaskSelect = `select ID, ifnull(USR_NME,'') USR_NME, ifnull(BIZ_DTE,'') BIZ_DTE, ifnull(WORK_RMK,'') WORK_RMK, ifnull(UPD_DTE,'') UPD_DTE	from t_dly_task`
)

// DailyTaskFindBy 根据条件查找
func DailyTaskFindBy(cond string) ([]DailyTask, error) {
	cmd := dailyTaskSelect + cond
	// log.Printf("DailyTask search: %s\n", cmd)

	items := []DailyTask{}
	err := db.Select(&items, cmd)
	return items, err
}

// DailyTaskFindAll 返回所有
func DailyTaskFindAll() ([]DailyTask, error) {
	return DailyTaskFindBy("")
}

// DailyTaskFindByID 根据 ID 加载记录
func DailyTaskFindByID(id int) (DailyTask, error) {
	item := DailyTask{}
	cmd := dailyTaskSelect + " where ID=? "
	err := db.Get(&item, cmd, id)
	return item, err
}

// DailyTaskInsert 把当前对象，作为新记录，插入数据库
func DailyTaskInsert(u DailyTask) (DailyTask, error) {
	item := DailyTask{}
	now := time.Now().Format("2006-01-02 15:04:05")
	u.UpdateDate = now

	sqlCmd := `INSERT INTO t_dly_task (USR_NME,BIZ_DTE,WORK_RMK,UPD_DTE) VALUES (:USR_NME,:BIZ_DTE,:WORK_RMK,:UPD_DTE)`
	res, err := db.NamedExec(sqlCmd, u)
	if err != nil {
		return item, err
	}

	// 取回插入的记录
	id, err := res.LastInsertId()
	if err != nil {
		return item, err
	}

	alog := ActivityLog{}
	alog.RefID = int(id)
	alog.RefType = "DAILYTASK"
	alog.Msg1 = u.WorkRemark
	alog.OpUser = u.UserName
	alog.OpDate = now
	ActivityLogInsert(alog)

	return DailyTaskFindByID(int(id))
}

// DailyTaskUpdate 将当前对象保存到数据库
func DailyTaskUpdate(u DailyTask) (DailyTask, error) {
	now := time.Now().Format("2006-01-02 15:04:05")
	u.UpdateDate = now

	// 只能更新 工作内容  和  修改时间
	sqlCmd := `update t_dly_task set
									WORK_RMK = :WORK_RMK,
									UPD_DTE = :UPD_DTE
								where id=:ID`
	_, err := db.NamedExec(sqlCmd, u)
	if err != nil {
		return DailyTask{}, err
	}

	alog := ActivityLog{}
	alog.RefID = u.ID
	alog.RefType = "DAILYTASK"
	alog.Msg1 = u.WorkRemark
	alog.OpUser = u.UserName
	alog.OpDate = now
	ActivityLogInsert(alog)

	// 重新加载记录
	return DailyTaskFindByID(u.ID)
}

// DailyTaskDelete 根据当前对象的 ID 删除
func DailyTaskDelete(u DailyTask) error {
	sqlCmd := `delete from t_dly_task where ID=:ID`
	_, err := db.NamedExec(sqlCmd, u)

	return err
}
