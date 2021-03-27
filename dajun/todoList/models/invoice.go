package models

import (
	"strconv"
	"time"
)

// Invoice 合同付款计划
type Invoice struct {
	ID           int    `db:"ID" json:"id" form:"id"`
	ContractID   int    `db:"CTRCT_ID" json:"contract_id" form:"contract_id"`
	MilestoneID  int    `db:"MLST_ID" json:"milestone_id" form:"milestone_id"`
	InvNum       string `db:"INV_NUM" json:"inv_num" form:"inv_num"`
	InvDate      string `db:"INV_DTE" json:"inv_date" form:"inv_date"`
	VenderName   string `db:"VENDER_NME" json:"vender_name" form:"vender_name"`
	InvType      string `db:"INV_TYP" json:"inv_type" form:"inv_type"`
	Currency     string `db:"CCY_CDE" json:"currency_code" form:"currency_code"`
	TotalAmount  string `db:"TTL_AMT" json:"total_amt" form:"total_amt"`
	TaxAmount    string `db:"TAX_AMT" json:"tax_amt" form:"tax_amt"`
	TaxRate      string `db:"TAX_RATE" json:"tax_rate" form:"tax_rate"`
	Remark       string `db:"RMK" json:"remark" form:"remark"`
	UpdateUser   string `db:"UPD_USR" json:"update_user" form:"update_user"`
	UpdateDate   string `db:"UPD_DTE" json:"update_date" form:"update_date"`
	FinUser      string `db:"FIN_USR" json:"fin_user" form:"fin_user"`
	FinDate      string `db:"FIN_DTE" json:"fin_date" form:"fin_date"`
	FinRemark    string `db:"FIN_RMK" json:"fin_remark" form:"fin_remark"`
	IsPaid       string `db:"IS_PAID" json:"is_paid" form:"is_paid"`
	IsPrePaid    string `db:"IS_PPAID" json:"is_prepaid" form:"is_prepaid"`
	PaidDate     string `db:"PAID_DTE" json:"paid_date" form:"paid_date"`
	PaidRemark   string `db:"PAID_RMK" json:"paid_remark" form:"paid_remark"`
	RequestUser  string `db:"REQ_USR" json:"request_user" form:"request_user"`
	HandoverUser string `db:"HDOV_USR" json:"handover_user" form:"handover_user"`
}

const (
	invoiceSelect = `select ID, ifnull(CTRCT_ID,0) CTRCT_ID, ifnull(MLST_ID,0) MLST_ID,ifnull(INV_NUM,'') INV_NUM, ifnull(INV_DTE,'') INV_DTE, ifnull(VENDER_NME,'') VENDER_NME, ifnull(INV_TYP,'') INV_TYP, ifnull(CCY_CDE,'') CCY_CDE, ifnull(TTL_AMT,'') TTL_AMT, ifnull(TAX_AMT,'') TAX_AMT, ifnull(TAX_RATE,'') TAX_RATE, ifnull(RMK,'') RMK, ifnull(UPD_USR,'') UPD_USR, ifnull(UPD_DTE,'') UPD_DTE, ifnull(FIN_USR,'') FIN_USR, ifnull(FIN_DTE,'') FIN_DTE, ifnull(IS_PAID,'') IS_PAID, ifnull(IS_PPAID,'') IS_PPAID, ifnull(FIN_RMK,'') FIN_RMK, ifnull(PAID_RMK,'') PAID_RMK, ifnull(PAID_DTE,'') PAID_DTE, ifnull(REQ_USR,'') REQ_USR, ifnull(HDOV_USR,'') HDOV_USR
	from t_ctrct_inv`
)

// InvoiceFindBy 根据条件查找
func InvoiceFindBy(cond string) ([]Invoice, error) {
	sqlCmd := invoiceSelect + cond

	items := []Invoice{}
	err := db.Select(&items, sqlCmd)

	return items, err
}

// InvoiceFindAll 返回部门清单
func InvoiceFindAll() ([]Invoice, error) {
	// 查询清单
	return InvoiceFindBy(" order by UPD_DTE desc   ")
}

// InvoiceFindByContract 根据合同id查询
func InvoiceFindByContract(contractID int) ([]Invoice, error) {
	return InvoiceFindBy(" where CTRCT_ID = " + strconv.Itoa(contractID) + " order by MLST_ID   ")
}

// InvoiceFindByID 按照 id 查询
func InvoiceFindByID(id int) (Invoice, error) {
	cmd := invoiceSelect + ` where ID=?`
	item := Invoice{}
	err := db.Get(&item, cmd, id)
	return item, err
}

// InvoiceInsert 当前对象，插入到数据库
func InvoiceInsert(c Invoice) (Invoice, error) {
	c.UpdateDate = time.Now().Format("2006-01-02 15:04:05")

	// 根据 struct 中的 DB tag 进行自动 named parameter
	cmd := `INSERT INTO t_ctrct_inv (CTRCT_ID,MLST_ID,INV_NUM,INV_DTE,VENDER_NME,INV_TYP,CCY_CDE,TTL_AMT,TAX_AMT,TAX_RATE,RMK,UPD_USR,UPD_DTE,FIN_USR,FIN_DTE,IS_PAID,IS_PPAID) VALUES
						(:CTRCT_ID,:MLST_ID,:INV_NUM,:INV_DTE,:VENDER_NME,:INV_TYP,:CCY_CDE,:TTL_AMT,:TAX_AMT,:TAX_RATE,:RMK,:UPD_USR,:UPD_DTE,:FIN_USR,:FIN_DTE,:IS_PAID,:IS_PPAID)`
	res, err := db.NamedExec(cmd, c)
	if err != nil {
		return Invoice{}, err
	}

	// 取回插入的记录
	id, err := res.LastInsertId()
	if err != nil {
		return Invoice{}, err
	}

	return InvoiceFindByID(int(id))
}

// InvoiceUpdate 当前对象，更新到数据库
func InvoiceUpdate(c Invoice) (Invoice, error) {
	c.UpdateDate = time.Now().Format("2006-01-02 15:04:05")

	cmd := `update t_ctrct_inv set
								CTRCT_ID = :CTRCT_ID,
								MLST_ID = :MLST_ID,
								INV_NUM = :INV_NUM,
								INV_DTE = :INV_DTE,
								VENDER_NME = :VENDER_NME,
								INV_TYP = :INV_TYP,
								CCY_CDE = :CCY_CDE,
								TTL_AMT = :TTL_AMT,
								TAX_AMT = :TAX_AMT,
								TAX_RATE = :TAX_RATE,
								RMK = :RMK,
								UPD_USR = :UPD_USR,
								UPD_DTE = :UPD_DTE,
								FIN_USR = :FIN_USR,
								FIN_DTE = :FIN_DTE,
								IS_PAID = :IS_PAID,
								IS_PPAID = :IS_PPAID
						where ID=:ID`
	_, err := db.NamedExec(cmd, c)
	if err != nil {
		return c, err
	}

	return InvoiceFindByID(c.ID)
}

// InvoiceDelete 当前对象，按照 ID 从数据库删除
func InvoiceDelete(c Invoice) error {
	// 按照 id 删除
	cmd := "delete from t_ctrct_plan where ID=:ID"
	_, err := db.NamedExec(cmd, c)
	return err
}

// InvoiceHandOver 发票交接
func InvoiceHandOver(c Invoice) error {
	c.UpdateDate = time.Now().Format("2006-01-02 15:04:05")

	cmd := `update t_ctrct_inv set
									HDOV_USR = :HDOV_USR,
									FIN_USR = :FIN_USR,
									FIN_DTE = :FIN_DTE,
									FIN_RMK = :FIN_RMK
							where ID=:ID`
	_, err := db.NamedExec(cmd, c)

	return err
}

// InvoicePaymentRequest 发票付款申请
func InvoicePaymentRequest(c Invoice) error {
	c.UpdateDate = time.Now().Format("2006-01-02 15:04:05")

	cmd := `update t_ctrct_inv set
								REQ_USR = :REQ_USR,
								IS_PAID = :IS_PAID,
								IS_PPAID = :IS_PPAID,
								PAID_DTE = :PAID_DTE,
								PAID_RMK = :PAID_RMK
							where ID=:ID`
	_, err := db.NamedExec(cmd, c)

	return err
}
