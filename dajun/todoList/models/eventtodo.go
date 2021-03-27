package models

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	// "github.com/xuri/excelize"
	"github.com/360EntSecGroup-Skylar/excelize"
)

// EventTodo 每个事项(Event)下的任务列表(todoList)
type EventTodo struct {
	ID            int    `db:"EVT_TODO_ID" json:"id" form:"id"`
	RefID         int    `db:"REF_ID" json:"ref_id" form:"ref_id"`
	RefType       string `db:"REF_TYP" json:"ref_type" form:"ref_type"`
	RefTitle      string `db:"REF_TITLE" json:"ref_title" form:"ref_title"`
	Seq           int    `db:"TODO_SEQ" json:"seq" form:"seq"`
	Content       string `db:"TODO_CTENT" json:"content" form:"content"`
	OwnerID       int    `db:"OWNER_ID" json:"owner_id" form:"owner_id"`
	OwnerName     string `db:"OWNER_NME" json:"owner_name" form:"owner_name"`
	DueDate       string `db:"DUE_DTE" json:"due_date" form:"due_date"`
	Status        string `db:"TODO_STS" json:"status" form:"status"`
	ActualCmpDate string `db:"ACT_CMPL_DTE" json:"actual_cmp_date" form:"actual_cmp_date"`
	Memo          string `db:"MEMO" json:"memo" form:"memo"`
	CreateUser    string `db:"CRE_USR" json:"create_user" form:"create_user"`
	CreateDate    string `db:"CRE_DT" json:"create_date" form:"create_date"`
	UpdateUser    string `db:"UPD_USR" json:"update_user" form:"update_user"`
	UpdateDate    string `db:"UPD_DT" json:"update_date" form:"update_date"`

	// 物料数量校验逻辑
	MatRule string `db:"MAT_RULE" json:"mat_rule" form:"mat_rule"`
	MatCode string `db:"MAT_CODE" json:"mat_code" form:"mat_code"`
	MatCond string `db:"MAT_COND" json:"mat_cond" form:"mat_cond"`
}

const (
	todoSelect = `select	EVT_TODO_ID, ifnull(REF_ID,'') REF_ID, ifnull(REF_TYP,'') REF_TYP, ifnull(REF_TITLE,'') REF_TITLE, ifnull(TODO_SEQ,'') TODO_SEQ, ifnull(TODO_CTENT,'') TODO_CTENT, ifnull(OWNER_ID,'') OWNER_ID, ifnull(OWNER_NME,'') OWNER_NME, ifnull(DUE_DTE,'') DUE_DTE, ifnull(TODO_STS,'') TODO_STS, ifnull(ACT_CMPL_DTE,'') ACT_CMPL_DTE, ifnull(MEMO,'') MEMO, ifnull(CRE_USR,'') CRE_USR, ifnull(CRE_DT,'') CRE_DT, ifnull(UPD_USR,'') UPD_USR, ifnull(UPD_DT,'') UPD_DT
									from t_event_todo `
)

// CreateEventTodo 返回一个空对象
func CreateEventTodo() EventTodo {
	return EventTodo{}
}

// FindBy 返回符合条件的记录
func (t EventTodo) FindBy(cond string) ([]EventTodo, error) {
	cmd := todoSelect

	if len(strings.TrimSpace(cond)) > 0 {
		cmd = cmd + " " + cond
	}

	// fmt.Printf("sql: %s\n", cmd)

	items := []EventTodo{}
	err := db.Select(&items, cmd)
	return items, err
}

// FindAll 返回所有记录
func (t EventTodo) FindAll() ([]EventTodo, error) {
	return t.FindBy("")
}

// FindByParam 根据 责任人 待办内容 模糊查找
func (t EventTodo) FindByParam(param string) ([]EventTodo, error) {
	// 未完成，或者显示近三个月内已完成
	// 按 责任人，待办事项的内容，模糊匹配
	cond := `'%` + param + `%'`
	cmd := ` where (TODO_CTENT like ` + cond + ` or OWNER_NME like ` + cond + `)
				and ( (TODO_STS is null or length(TODO_STS) = 0) or (CRE_DT >= date_format(date_add(now(), interval -90 day), '%Y-%m-%d') ))
			order by TODO_STS asc, DUE_DTE asc 	`

	return t.FindBy(cmd)
}

// FindByRef 根据 ref 查找，参数是当前对象的属性
func (t EventTodo) FindByRef() ([]EventTodo, error) {
	// 查询的返回结果，默认是空列表
	items := []EventTodo{}
	// 查询: 状态为空的显示在上面（未完成），按照到期日，从早到晚显示
	cmd := todoSelect +
		` where REF_TYP=:REF_TYP and REF_ID = :REF_ID
		order by TODO_STS asc, DUE_DTE asc
		`
	nstmt, err := db.PrepareNamed(cmd)
	if err != nil {
		return items, err
	}

	// 通过当前对象的属性，传递查询条件
	err = nstmt.Select(&items, t)
	return items, err

}

// FindByID 根据 ID 查找，返回一个新对象
func (t EventTodo) FindByID(id int) (EventTodo, error) {
	cmd := eventSelect +
		` where EVT_TODO_ID=?`

	item := EventTodo{}
	err := db.Get(&item, cmd, id)

	return item, err
}

// Insert 当前对象，插入到数据库中
func (t EventTodo) Insert() error {
	// 创建时间为当前时刻
	t.CreateDate = time.Now().Format("2006-01-02 15:04:05")

	// 根据 struct 中的 DB tag 进行自动 named parameter
	insertCmd := `INSERT INTO t_event_todo
						 	(REF_ID, REF_TYP, REF_TITLE, TODO_SEQ, TODO_CTENT, OWNER_ID, OWNER_NME, DUE_DTE, TODO_STS, ACT_CMPL_DTE, MEMO, CRE_USR, CRE_DT, UPD_USR, UPD_DT)
						VALUES
					 		(:REF_ID, :REF_TYP, :REF_TITLE, :TODO_SEQ, :TODO_CTENT, :OWNER_ID, :OWNER_NME, :DUE_DTE, :TODO_STS, :ACT_CMPL_DTE, :MEMO, :CRE_USR, :CRE_DT, :UPD_USR, :UPD_DT)`
	_, err := db.NamedExec(insertCmd, t)
	if err != nil {
		return err
	}

	// 在事件下添加了新的待办，将该事件状态变为 OPEN
	err = CloseEvent(t.RefID, false)
	return err
}

// Update 当前对象更新回数据库
func (t EventTodo) Update() error {
	// 设置更新时间为当前时刻
	t.UpdateDate = time.Now().Format("2006-01-02 15:04:05")

	// 根据 struct 中的 tag 进行自动 named parameter
	sqlCmd := `update t_event_todo set
										REF_ID = :REF_ID,
										REF_TYP = :REF_TYP,
										REF_TITLE = :REF_TITLE,
										TODO_SEQ = :TODO_SEQ,
										TODO_CTENT = :TODO_CTENT,
										OWNER_ID = :OWNER_ID,
										OWNER_NME = :OWNER_NME,
										DUE_DTE = :DUE_DTE,
										TODO_STS = :TODO_STS,
										ACT_CMPL_DTE = :ACT_CMPL_DTE,
										MEMO = :MEMO,
										CRE_USR = :CRE_USR,
										CRE_DT = :CRE_DT,
										UPD_USR = :UPD_USR,
										UPD_DT = :UPD_DT
								where EVT_TODO_ID=:EVT_TODO_ID`
	_, err := db.NamedExec(sqlCmd, t)

	return err
}

// Delete 从数据库删除当前对象
func (t EventTodo) Delete() error {
	// 按照 id 删除
	cmd := "delete from t_event_todo where EVT_TODO_ID=?"
	_, err := db.Exec(cmd, t.ID)
	return err
}

// UpdateStatus 保存当前对象的 status 到数据库
func (t EventTodo) UpdateStatus() error {
	t.UpdateDate = time.Now().Format("2006-01-02 15:04:05")

	// 根据 struct 中的 tag 进行自动 named parameter
	sqlCmd := `update t_event_todo set
				TODO_STS = :TODO_STS,
				ACT_CMPL_DTE = :ACT_CMPL_DTE,
				MEMO = :MEMO
				where EVT_TODO_ID=:EVT_TODO_ID`
	_, err := db.NamedExec(sqlCmd, t)

	return err
}

// LoadTodos 从 reader 中，读取数据，返回对象列表
// Reader 的格式是固定的
// tab 名字必须是 Sheet1
// 第一行是标题；
func (t EventTodo) LoadTodos(file io.Reader) ([]EventTodo, error) {
	// 每一列也是固定的；
	items := []EventTodo{}

	// 从文件中读取
	xlsx, err := excelize.OpenReader(file)
	if err != nil {
		return nil, err
	}

	// sheet 的名字是固定的，只能是 Sheet1
	index := xlsx.GetSheetIndex("Sheet1")

	// Get all the rows in a sheet.
	rows := xlsx.GetRows("sheet" + strconv.Itoa(index))
	// if err != nil {
	// 	println(err.Error())
	// 	return items, err
	// }

	fmt.Printf("total rows: %d\n", len(rows))

	for idx, row := range rows {
		// 表头，跳过
		if idx == 0 {
			continue
		}

		// 第一列，序号为空，表示是最后一行
		if len(row[0]) == 0 {
			break
		}

		// 任务 或 责任人 为空，跳过该行
		if row[1] == "" || row[2] == "" {
			continue
		}

		i := EventTodo{}
		i.Seq, err = strconv.Atoi(row[0])
		if err != nil {
			return nil, err
		}

		i.Content = row[1]
		i.OwnerName = row[2]

		// 检查日期格式
		_, err := time.Parse("2006-01-02", row[3])
		if err != nil {
			fmt.Printf("bad date: %+v\n", err)
			return nil, err
		}

		i.DueDate = row[3]
		i.Memo = row[4]

		items = append(items, i)
	}

	return items, nil
}

// SetMatRule 设置物料数量的提醒逻辑
func SetMatRule(todoid int, matRule, matCode, matCond string) error {
	sqlCmd := `update t_event_todo set
			MAT_RULE = :MAT_RULE,
			MAT_CODE = :MAT_CODE,
			MAT_COND = :MAT_COND
		where EVT_TODO_ID=:ID`

	// 一样可以使用 map 作为绑定的参数
	p := map[string]interface{}{"ID": todoid, "MAT_RULE": matRule, "MAT_CODE": matCode, "MAT_COND": matCond}
	_, err := db.NamedExec(sqlCmd, p)

	return err
}
