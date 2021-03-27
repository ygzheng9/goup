package models

import "time"

// Article 公告信息
type Article struct {
	ID         int    `db:"ID" json:"id" form:"id"`
	UpdateDate string `db:"UPD_DTE" json:"update_date" form:"update_date"`
	UpdateUser string `db:"UPD_USR" json:"update_user" form:"update_user"`
	Title      string `db:"ARTICLE_TITLE" json:"title" form:"title"`
	Keywords   string `db:"ARTICLE_KEYWD" json:"keywords" form:"keywords"`
	Content    string `db:"CTENT" json:"content" form:"content"`
}

const (
	articleSelect = `select ID, ifnull(UPD_DTE,'') UPD_DTE, ifnull(UPD_USR,'') UPD_USR, ifnull(ARTICLE_TITLE,'') ARTICLE_TITLE, ifnull(ARTICLE_KEYWD,'') ARTICLE_KEYWD, ifnull(CTENT,'') CTENT	from t_article`
)

// ArticleFindBy 根据条件查找
func ArticleFindBy(cond string) ([]Article, error) {
	cmd := articleSelect + cond
	// log.Printf("DailyTask search: %s\n", cmd)

	items := []Article{}
	err := db.Select(&items, cmd)
	return items, err
}

// ArticleFindAll 返回所有
func ArticleFindAll() ([]Article, error) {
	return ArticleFindBy("")
}

// ArticleFindByID 根据 ID 加载记录
func ArticleFindByID(id int) (Article, error) {
	item := Article{}
	cmd := articleSelect + " where ID=? "
	err := db.Get(&item, cmd, id)
	return item, err
}

// ArticleInsert 把当前对象，作为新记录，插入数据库
func ArticleInsert(u Article) (Article, error) {
	item := Article{}
	now := time.Now().Format("2006-01-02 15:04:05")
	u.UpdateDate = now

	sqlCmd := `INSERT INTO t_article (UPD_DTE,UPD_USR,ARTICLE_TITLE,ARTICLE_KEYWD,CTENT) VALUES (:UPD_DTE,:UPD_USR,:ARTICLE_TITLE,:ARTICLE_KEYWD,:CTENT)`
	res, err := db.NamedExec(sqlCmd, u)
	if err != nil {
		return item, err
	}

	// 取回插入的记录
	id, err := res.LastInsertId()
	if err != nil {
		return item, err
	}

	return ArticleFindByID(int(id))
}

// ArticleUpdate 当前对象更新回数据库
func ArticleUpdate(t Article) (Article, error) {
	// 设置更新时间为当前时刻
	t.UpdateDate = time.Now().Format("2006-01-02 15:04:05")

	// 根据 struct 中的 tag 进行自动 named parameter
	sqlCmd := `update t_article set
								UPD_DTE = :UPD_DTE,
								UPD_USR = :UPD_USR,
								ARTICLE_TITLE = :ARTICLE_TITLE,
								ARTICLE_KEYWD = :ARTICLE_KEYWD,
								CTENT = :CTENT
								where ID=:ID`
	_, err := db.NamedExec(sqlCmd, t)
	if err != nil {
		return t, err
	}

	return ArticleFindByID(t.ID)
}

// ArticleDelete 根据当前对象的 ID 删除
func ArticleDelete(u Article) error {
	sqlCmd := `delete from t_article where ID=:ID`
	_, err := db.NamedExec(sqlCmd, u)

	return err
}
