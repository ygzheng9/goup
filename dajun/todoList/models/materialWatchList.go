package models

import (
	"time"

	"pickup/dajun/todoList/u8"
)

// 查看 BOM 是否已经存在

// MaterialWatchItem 关注的物料号
type MaterialWatchItem struct {
	ID         int    `db:"ID" json:"id" form:"id"`
	CreateDate string `db:"CRE_DTE" json:"create_date" form:"create_date"`
	CreateUser string `db:"CRE_USR" json:"create_user" form:"create_user"`
	InvCode    string `db:"INV_CDE" json:"inv_code" form:"inv_code"`
	Remark     string `db:"RMK" json:"remark" form:"remark"`
}

const (
	materailWLSelect = `select ID, ifnull(CRE_DTE,'') CRE_DTE, ifnull(CRE_USR,'') CRE_USR, ifnull(INV_CDE,'') INV_CDE,  ifnull(RMK,'') RMK from t_mat_watchlist`
)

// MatWLFindBy 根据条件查找
func MatWLFindBy(cond string) ([]MaterialWatchItem, error) {
	cmd := materailWLSelect + cond
	// log.Printf("DailyTask search: %s\n", cmd)

	items := []MaterialWatchItem{}
	err := db.Select(&items, cmd)
	return items, err
}

// MatWLFindAll 返回所有
func MatWLFindAll() ([]MaterialWatchItem, error) {
	return MatWLFindBy("")
}

// MatWLFindByID 根据 ID 加载记录
func MatWLFindByID(id int) (MaterialWatchItem, error) {
	item := MaterialWatchItem{}
	cmd := materailWLSelect + " where ID=? "
	err := db.Get(&item, cmd, id)
	return item, err
}

// MatWLInsert 把当前对象，作为新记录，插入数据库
func MatWLInsert(u MaterialWatchItem) (MaterialWatchItem, error) {
	item := MaterialWatchItem{}
	now := time.Now().Format("2006-01-02 15:04:05")
	u.CreateDate = now

	sqlCmd := `INSERT INTO t_mat_watchlist (CRE_DTE,CRE_USR,INV_CDE, RMK) VALUES (:CRE_DTE, :CRE_USR, :INV_CDE, :RMK)`
	res, err := db.NamedExec(sqlCmd, u)
	if err != nil {
		return item, err
	}

	// 取回插入的记录
	id, err := res.LastInsertId()
	if err != nil {
		return item, err
	}

	return MatWLFindByID(int(id))
}

// MatWLDelete 根据当前对象的 ID 删除
func MatWLDelete(u MaterialWatchItem) error {
	sqlCmd := `delete from t_mat_watchlist where ID=:ID`
	_, err := db.NamedExec(sqlCmd, u)

	return err
}

// MatBOM 物料对应的 BOM
type MatBOM struct {
	InvCode  string `db:"cInvCode" json:"inv_code"`
	InvName  string `db:"cInvName" json:"inv_name"`
	InvStd   string `db:"cInvStd" json:"inv_std"`
	DesignNo string `db:"cDesignNo" json:"design_no"`
	VerDesc  string `db:"VersionDesc" json:"ver_desc"`
	VerEff   string `db:"VersionEffDate" json:"ver_eff"`
}

// MatWLLoadBOM 物料对应的BOM信息
func MatWLLoadBOM(input []MaterialWatchItem) ([]MatBOM, error) {
	items := []MatBOM{}

	totalMat := len(input)
	if totalMat == 0 {
		return items, nil
	}

	temp := " select '" + input[0].InvCode + "' id "
	for idx := 1; idx < totalMat; idx++ {
		temp = temp + "  union all select '" + input[idx].InvCode + "' id"
	}

	cmd := `
			with a_list as (
				-- 构建清单
					select id from (
		` + temp + `
		) a
		),
		matched as (
		-- 这些料号中，有 BOM 的
		select l.id, a.bomid, l.id cInvCode, d.cInvName, d.cInvAddCode, d.cInvStd,
			  d.cInvDefine7,
				b.ParentScrap,
				d.cInvDefine2 VersionDesc, a.VersionEffDate, a.Status
			from a_list l
			left join inventory d on d.cInvCode = l.id
			left join bas_part c on c.InvCode = d.cInvCode
			left join bom_parent b on b.parentId = c.partid
			left join bom_bom a on a.bomid = b.bomid
		)
		select m.cInvCode, isnull(m.cInvName, '') cInvName, isnull(m.cInvStd, '') cInvStd, isnull(m.cInvDefine7, '') cDesignNo,
					 isnull(m.VersionDesc, '') VersionDesc, isnull(m.VersionEffDate, '') VersionEffDate
			from   matched m ;	`

	err := u8.GetDB().Select(&items, cmd)
	return items, err
}
