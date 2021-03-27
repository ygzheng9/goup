package models

// Asset 记录：资产编号，名称，建档时间，保管人，使用人，位置
type Asset struct {
	ID            int    `db:"ID" json:"id" form:"id"`
	AssetNum      string `db:"ASSET_NUM" json:"asset_number" form:"asset_number"`
	AssetName     string `db:"ASSET_NME" json:"asset_name" form:"asset_name"`
	AssetDate     string `db:"ASSET_DTE" json:"asset_date" form:"asset_date"`
	Department    string `db:"OWN_DEPT" json:"asset_department" form:"asset_department"`
	AssetUser     string `db:"OWN_USR" json:"asset_user" form:"asset_user"`
	AssetKeeper   string `db:"ASSET_KEEPER" json:"asset_keeper" form:"asset_keeper"`
	AssetLocation string `db:"ASSET_LOC" json:"asset_location" form:"asset_location"`
	Remark        string `db:"RMK" json:"asset_remark" form:"asset_remark"`
}

const (
	assetSelect = `select ID, ifnull(ASSET_NUM,'') ASSET_NUM, ifnull(ASSET_NME,'') ASSET_NME, ifnull(ASSET_DTE,'') ASSET_DTE, ifnull(OWN_DEPT,'') OWN_DEPT, ifnull(OWN_USR,'') OWN_USR, ifnull(ASSET_KEEPER,'') ASSET_KEEPER, ifnull(ASSET_LOC,'') ASSET_LOC, ifnull(RMK,'') RMK
								from	 t_asset`
)

// CreateAsset 返回一个空对象
func CreateAsset() Asset {
	return Asset{}
}

// FindBy 根据条件查询
func (t Asset) FindBy(cond string) ([]Asset, error) {
	cmd := assetSelect + " " + cond
	items := []Asset{}
	err := db.Select(&items, cmd)
	return items, err
}

// FindAll 返回全部
func (t Asset) FindAll() ([]Asset, error) {
	return t.FindBy("")
}

// FindByID 根据ID查找
func (t Asset) FindByID(id int) (Asset, error) {
	// 按照 id 查询
	cmd := assetSelect +
		` where ID=?`

	item := Asset{}
	err := db.Get(&item, cmd, id)

	return Asset{}, err
}

// Insert 把当前对象插入到数据库
func (t Asset) Insert() error {
	// 根据 struct 中的 DB tag 进行自动 named parameter
	insertCmd := `INSERT INTO t_asset (ASSET_NUM, ASSET_NME, ASSET_DTE, OWN_DEPT, OWN_USR, ASSET_KEEPER, ASSET_LOC, RMK)
					VALUES (:ASSET_NUM, :ASSET_NME, :ASSET_DTE, :OWN_DEPT, :OWN_USR, :ASSET_KEEPER, :ASSET_LOC, :RMK)`

	_, err := db.NamedExec(insertCmd, t)
	return err
}

// Update 更新数据库，当前对象
func (t Asset) Update() error {
	// 根据 struct 中的 tag 进行自动 named parameter
	sqlCmd := `update t_asset
							set ASSET_NUM = :ASSET_NUM,
								ASSET_NME = :ASSET_NME,
								ASSET_DTE = :ASSET_DTE,
								OWN_DEPT = :OWN_DEPT,
								OWN_USR = :OWN_USR,
								ASSET_KEEPER = :ASSET_KEEPER,
								ASSET_LOC = :ASSET_LOC,
								RMK = :RMK
							where ID=:ID`
	_, err := db.NamedExec(sqlCmd, t)

	return err
}

// Delete 从数据库中删除当前对象
func (t Asset) Delete() error {
	// 按照 id 删除
	cmd := "delete from t_asset where ID=?"
	_, err := db.Exec(cmd, t)

	return err
}
