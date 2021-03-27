package models

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/robfig/cron"

	gomail "gopkg.in/gomail.v2"
)

// EmailNotice 每一个邮件通知
type EmailNotice struct {
	ID         int    `db:"ID" json:"id" form:"id"`
	CreateDate string `db:"CRE_DTE" json:"create_date" form:"create_date"`
	SendTo     string `db:"SEND_TO" json:"send_to" form:"send_to"`
	UserName   string `db:"USR_NME" json:"user_name" form:"user_name"`
	Subject    string `db:"SUBJ" json:"subject" form:"subject"`
	Content    string `db:"CTENT" json:"content" form:"content"`
	Status     string `db:"STS" json:"status" form:"status"`
	SentDate   string `db:"SENT_DTE" json:"sent_date" form:"sent_date"`
}

const (
	// WAITING 邮件未发送
	WAITING = "WAITING"

	// SENT 邮件已发送
	SENT = "SENT"
)

const (
	emailNoticeSelect = `select
				ID, ifnull(CRE_DTE,'') CRE_DTE, ifnull(SEND_TO,'') SEND_TO, ifnull(USR_NME,'') USR_NME, ifnull(SUBJ,'') SUBJ, ifnull(CTENT,'') CTENT, ifnull(STS,'') STS, ifnull(SENT_DTE,'') SENT_DTE
		from t_email `
)

// CreateEmailNotice 返回一个空对象
func CreateEmailNotice() EmailNotice {
	return EmailNotice{}
}

// FindBy 返回符合条件的记录
func (t EmailNotice) FindBy(cond string) ([]EmailNotice, error) {
	cmd := emailNoticeSelect

	if len(strings.TrimSpace(cond)) > 0 {
		cmd = cmd + " " + cond
	}

	items := []EmailNotice{}
	err := db.Select(&items, cmd)
	return items, err
}

// FindAll 返回所有记录
func (t EmailNotice) FindAll() ([]EmailNotice, error) {
	return t.FindBy("")
}

// FindByID 根据 ID 查找，返回一个新对象
func (t EmailNotice) FindByID(id int) (EmailNotice, error) {
	cmd := emailNoticeSelect +
		` where ID=?`

	item := EmailNotice{}
	err := db.Get(&item, cmd, id)

	return item, err
}

// Insert 当前对象，插入到数据库中
func (t EmailNotice) Insert() error {
	// 创建时间为当前时刻
	t.CreateDate = time.Now().Format("2006-01-02 15:04:05")
	t.Status = WAITING

	// 根据 struct 中的 DB tag 进行自动 named parameter
	insertCmd := `INSERT INTO t_email
									(CRE_DTE, SEND_TO, USR_NME, SUBJ, CTENT, STS, SENT_DTE)
								VALUES
									(:CRE_DTE, :SEND_TO, :USR_NME, :SUBJ, :CTENT, :STS, :SENT_DTE)`
	_, err := db.NamedExec(insertCmd, t)

	return err
}

// Update 当前对象更新回数据库
func (t EmailNotice) Update() error {
	// 根据 struct 中的 tag 进行自动 named parameter
	sqlCmd := `update t_email set
							CRE_DTE = :CRE_DTE,
							SEND_TO = :SEND_TO,
							USR_NME = :USR_NME,
							SUBJ = :SUBJ,
							CTENT = :CTENT,
							STS = :STS,
							SENT_DTE = :SENT_DTE
							where ID=:ID`
	_, err := db.NamedExec(sqlCmd, t)

	return err
}

// Delete 从数据库删除当前对象
func (t EmailNotice) Delete() error {
	// 按照 id 删除
	cmd := "delete from t_email where ID=?"
	_, err := db.Exec(cmd, t.ID)
	return err
}

// SendAllEmails 发送所有邮件
func SendAllEmails() error {
	t := CreateEmailNotice()

	// 取得所有未发送的邮件
	cond := " where sts = '" + WAITING + "' "
	items, err := t.FindBy(cond)
	if err != nil {
		return err
	}

	// 发送邮件
	for _, i := range items {
		err = i.Send()

		if err != nil {
			// 发送失败后，打印错误信息，然后继续发下一个，而不是终止
			log.Printf("%d send failed: %s.\n", i.ID, err.Error())
		}
	}

	return nil
}

// Send 发送通知，并且标记已发送
func (t EmailNotice) Send() error {
	// 发送邮件
	m := gomail.NewMessage()
	// 发件人
	m.SetAddressHeader("From", "bpit@dajuntech.com", "流程与IT")

	// 收件人
	m.SetHeader("To", m.FormatAddress(t.SendTo, t.UserName))

	// 主题
	m.SetHeader("Subject", t.Subject)

	// 正文
	m.SetBody("text/html", t.Content)

	d := gomail.NewDialer("smtp.qiye.163.com", 25, "bpit@dajuntech.com", "zaq12wsX") // 发送邮件服务器、端口、发件人账号、发件人密码
	err := d.DialAndSend(m)
	if err != nil {
		return err
	}

	// 标记已发送
	err = t.MarkSent()
	return err
}

// MarkSent 标记已发送
func (t EmailNotice) MarkSent() error {
	// 设置已发送
	t.Status = SENT
	t.SentDate = time.Now().Format("2006-01-02 15:04:05")

	// 根据 struct 中的 tag 进行自动 named parameter
	sqlCmd := `update t_email set
						STS = :STS,
						SENT_DTE = :SENT_DTE
						where ID=:ID`
	_, err := db.NamedExec(sqlCmd, t)

	return err
}

// StartNoticeSvc 启动定时任务
func StartNoticeSvc() {

	c := cron.New()

	// 定时发邮件，每 15 分钟发一次邮件
	spec1 := "30 */15 8-20 * * *"
	c.AddFunc(spec1, func() {
		SendAllEmails()
		fmt.Printf("send mails: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	})

	// 定时扫描未完成待办 每天 8 / 12  点
	spec2 := "0 1 8,13 * * *"
	c.AddFunc(spec2, func() {
		ScanAllTodos()
		ScanMatRule(0)
		DiskSpaceMonitor()
		fmt.Printf("scan: %s", time.Now().Format("2006-01-02 15:04:05"))
	})

	c.Start()
}
