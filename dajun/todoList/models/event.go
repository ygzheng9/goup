package models

import (
	"fmt"
	"strconv"
)

// Event 事件表
type Event struct {
	ID           int    `db:"EVT_ID" form:"id" json:"id" `
	Type         string `db:"TYP" form:"type" json:"type" `
	Department   string `db:"DEPT" form:"department" json:"department" `
	UserID       int    `db:"USR_ID" form:"user_id" json:"user_id" `
	UserName     string `db:"USR_NME" form:"user_name" json:"user_name" `
	HappenAt     string `db:"EVT_DT" form:"happen_at" json:"happen_at" `
	Place        string `db:"PLACE" form:"place" json:"place" `
	Subject      string `db:"SUBJ" form:"subject" json:"subject" `
	Memo         string `db:"MEMO" form:"memo" json:"memo" `
	CreateUser   string `db:"CRE_USR" form:"create_user" json:"create_user" `
	CreateDate   string `db:"CRE_DT" form:"create_date" json:"create_date" `
	UpdateUser   string `db:"UPD_USR" form:"update_user" json:"update_user" `
	UpdateDate   string `db:"UPD_DT" form:"update_date" json:"update_date" `
	MaterialDesc string `db:"MATRL_DESC" form:"material_desc" json:"material_desc" `
	EventStatus  string `db:"EVENT_STATUS" form:"event_status" json:"event_status" `

	OpenCnt int `db:"OPEN_CNT" form:"open_cnt" json:"open_cnt" `
}

const (
	eventSelect = `select EVT_ID, TYP, DEPT, USR_ID, USR_NME, EVT_DT, PLACE,
										SUBJ, ifnull(MEMO,"") MEMO,
										ifnull(CRE_USR,"") CRE_USR, ifnull(CRE_DT,"") CRE_DT,
										ifnull(UPD_USR,"") UPD_USR, ifnull(UPD_DT,"") UPD_DT,
										ifnull(MATRL_DESC,"") MATRL_DESC,
										ifnull(EVENT_STATUS,"") EVENT_STATUS
										from t_event`
)

// CreateEvent 返回一个空对象
func CreateEvent() Event {
	return Event{}
}

// FindBy 根据条件查询
func (t Event) FindBy(cond string) ([]Event, error) {
	cmd := eventSelect + " " + cond
	// fmt.Printf("event: %s\n", cmd)

	items := []Event{}
	err := db.Select(&items, cmd)
	if err != nil {
		return items, err
	}

	if len(items) == 0 {
		return items, nil
	}

	// 对每一个 event，查找 open todo 的数量
	items, err = updateOpenCnt(items)
	return items, err
}

func updateOpenCnt(items []Event) ([]Event, error) {
	if len(items) == 0 {
		return items, nil
	}

	// 把 items 中的 id ，拼接成 in (1,2,3)
	sqlCmd := ""
	cnt := 0
	for _, v := range items {
		if cnt == 0 {
			sqlCmd = " and a.evt_id in ( " + strconv.Itoa(v.ID)
		} else {
			sqlCmd = sqlCmd + ", " + strconv.Itoa(v.ID)
		}
		cnt++
	}
	sqlCmd = sqlCmd + ") "

	sel := ` select a.EVT_ID, count(1) OPEN_CNT
							from t_event a
						inner join   t_event_todo b on a.evt_id = b.ref_id and b.ref_typ = 'Events'
							where ( (b.todo_sts is null) or (b.todo_sts = '') ) `
	sqlCmd = sel + sqlCmd + " group by a.evt_id "

	// fmt.Println(sqlCmd)

	type resT struct {
		ID      int `db:"EVT_ID"`
		OpenCnt int `db:"OPEN_CNT"`
	}
	opens := []resT{}
	err := db.Select(&opens, sqlCmd)
	if err != nil {
		return items, err
	}

	// 按照 id 匹配
	itemsLen := len(items)
	for _, j := range opens {
		for i := 0; i < itemsLen; i++ {
			if items[i].ID == j.ID {
				items[i].OpenCnt = j.OpenCnt
				break
			}
		}
	}

	return items, nil
}

// FindAll 返回全部
func (t Event) FindAll() ([]Event, error) {
	return t.FindBy("")
}

// FindByID 根据ID查找
func (t Event) FindByID(id int) (Event, error) {
	// 按照 id 查询
	cond := fmt.Sprintf(" where EVT_ID=%d", id)
	items, err := t.FindBy(cond)

	if len(items) > 0 {
		return items[0], err
	}

	return Event{}, err
}

// Insert 把当前对象插入到数据库
func (t Event) Insert() error {
	// 新增事件，状态是 OPEN
	t.EventStatus = "OPEN"
	// 根据 struct 中的 DB tag 进行自动 named parameter
	insertCmd := `INSERT INTO t_event (TYP, DEPT, USR_ID, USR_NME, EVT_DT, PLACE, SUBJ, MEMO, CRE_USR, CRE_DT, MATRL_DESC, EVENT_STATUS) VALUES
						(:TYP, :DEPT, :USR_ID, :USR_NME, :EVT_DT, :PLACE, :SUBJ, :MEMO, :CRE_USR, :CRE_DT, :MATRL_DESC, :EVENT_STATUS)`
	_, err := db.NamedExec(insertCmd, &t)
	return err
}

// Update 更新数据库，当前对象
func (t Event) Update() error {
	// 根据 struct 中的 tag 进行自动 named parameter
	sqlCmd := `update t_event
					set DEPT=:DEPT, USR_ID=:USR_ID, USR_NME=:USR_NME, EVT_DT=:EVT_DT,
						PLACE=:PLACE, SUBJ=:SUBJ, MEMO=:MEMO, UPD_USR=:UPD_USR, UPD_DT=:UPD_DT, MATRL_DESC=:MATRL_DESC
					where EVT_ID=:EVT_ID`
	_, err := db.NamedExec(sqlCmd, t)

	return err
}

// Delete 从数据库中删除当前对象
func (t Event) Delete() error {
	tx, err := db.Begin()

	// 按照 id 删除
	cmd := "delete from t_event where EVT_ID=" + strconv.Itoa(t.ID)
	_, err = tx.Exec(cmd)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 删除子表
	cmd = "delete from t_event_todo where REF_TYP='Events' and REF_ID=" + strconv.Itoa(t.ID)
	_, err = tx.Exec(cmd)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

// CloseEvent 关闭 event
func CloseEvent(eventid int, close bool) error {
	status := "OPEN"
	if close {
		status = "CLOSED"
	}
	sqlCmd := `update t_event set EVENT_STATUS='` + status + `'	where EVT_ID=?`
	_, err := db.Exec(sqlCmd, eventid)

	return err
}

// FindEventByParam 根据条件查找 活动
func FindEventByParam(start, end, status, cond string) ([]Event, error) {
	// 拼接查询的 sql
	where := " where (EVT_DT >= '" + start + "' and EVT_DT <= '" + end + "' ) "
	if len(status) > 0 && status != "ALL" {
		where = where + " and EVENT_STATUS = '" + status + "' "
	}

	if len(cond) > 0 {
		where = where + " and ( (DEPT like '%" + cond + "%') or (USR_NME like '%" + cond + "%') or (SUBJ like '%" + cond + "%') or (MEMO like '%" + cond + "%') )"
	}

	item := CreateEvent()
	return item.FindBy(where)
}
