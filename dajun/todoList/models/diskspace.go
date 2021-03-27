package models

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/rpc"
	"pickup/dajun/todoList/server"
	"strconv"
	"time"
)

// DiskSpace 磁盘空间
type DiskSpace struct {
	ID           int     `db:"ID" json:"id" form:"id"`
	ServerName   string  `db:"SRVR_NME" json:"server_name" form:"server_name"`
	ServerIP     string  `db:"SRVR_IP" json:"server_ip" form:"server_ip"`
	DiskName     string  `db:"DISK_NME" json:"disk_name" form:"disk_name"`
	WarningPoint int     `db:"WARNING_PT" json:"warning_point" form:"warning_point"`
	NoticeUser   string  `db:"NOTICE_USR" json:"notice_user" form:"notice_user"`
	Status       string  `db:"STS_IND" json:"status_ind" form:"status_ind"`
	LastTick     string  `db:"L_TICK" json:"last_tick" form:"last_tick"`
	LastTotalAmt float64 `db:"L_TTL_AMT" json:"last_totalAmt" form:"last_totalAmt"`
	LastFreeAmt  float64 `db:"L_FREE_AMT" json:"last_freeAmt" form:"last_freeAmt"`
}

const (
	diskSelect = `select ID, ifnull(SRVR_NME,'') SRVR_NME, ifnull(SRVR_IP,'') SRVR_IP, ifnull(DISK_NME,'') DISK_NME, ifnull(WARNING_PT,'') WARNING_PT, ifnull(NOTICE_USR,'') NOTICE_USR, ifnull(STS_IND,'') STS_IND, ifnull(L_TICK,'') L_TICK, ifnull(L_TTL_AMT,'') L_TTL_AMT, ifnull(L_FREE_AMT,'') L_FREE_AMT
			from t_diskspace `
)

// CreateDiskSpace 返回一个空的 DiskSpace 对象
func CreateDiskSpace() DiskSpace {
	return DiskSpace{}
}

// DiskSpaceFindBy 根据条件查找
func DiskSpaceFindBy(cond string) ([]DiskSpace, error) {
	sqlCmd := diskSelect + cond
	// fmt.Printf("%s\n", sqlCmd)

	items := []DiskSpace{}
	err := db.Select(&items, sqlCmd)

	return items, err
}

// DiskSpaceFindAll 返回部门清单
func DiskSpaceFindAll() ([]DiskSpace, error) {
	// 查询清单
	return DiskSpaceFindBy(" order by SRVR_IP ")
}

// DiskSpaceFindByID 按照 id 查询
func DiskSpaceFindByID(id int) (DiskSpace, error) {
	cmd := diskSelect + ` where ID=?`

	item := DiskSpace{}
	err := db.Get(&item, cmd, id)
	return item, err
}

// Insert 当前对象，插入到数据库
func (d DiskSpace) Insert() error {
	if len(d.Status) == 0 {
		d.Status = "ACTIVE"
	}
	// 根据 struct 中的 DB tag 进行自动 named parameter
	cmd := `INSERT INTO t_diskspace (SRVR_NME,SRVR_IP,DISK_NME,WARNING_PT,NOTICE_USR,STS_IND,L_TICK,L_TTL_AMT,L_FREE_AMT) VALUES
						(:SRVR_NME,:SRVR_IP,:DISK_NME,:WARNING_PT,:NOTICE_USR,:STS_IND,:L_TICK,:L_TTL_AMT,:L_FREE_AMT)`
	_, err := db.NamedExec(cmd, d)
	return err
}

// Update 当前对象，更新到数据库
func (d DiskSpace) Update() error {
	cmd := `update t_diskspace
					  set SRVR_NME = :SRVR_NME,
								SRVR_IP = :SRVR_IP,
								DISK_NME = :DISK_NME,
								WARNING_PT = :WARNING_PT,
								NOTICE_USR = :NOTICE_USR,
								STS_IND = :STS_IND
						where ID=:ID`
	_, err := db.NamedExec(cmd, d)
	return err
}

// UpdateSpace 更新当前的容量信息
func (d DiskSpace) UpdateSpace() error {
	cmd := `update t_diskspace
					  set L_TICK = :L_TICK,
								L_TTL_AMT = :L_TTL_AMT,
								L_FREE_AMT = :L_FREE_AMT
						where ID=:ID`
	_, err := db.NamedExec(cmd, d)
	return err
}

// UpdateStatus 更新当前的状态信息
func (d DiskSpace) UpdateStatus() error {
	cmd := `update t_diskspace
					  set STS_IND = :STS_IND
						where ID=:ID`
	_, err := db.NamedExec(cmd, d)
	return err
}

// Delete 当前对象，按照 ID 从数据库删除
func (d DiskSpace) Delete() error {
	// 按照 id 删除
	cmd := "delete from t_diskspace where ID=:ID"
	_, err := db.NamedExec(cmd, d)
	return err
}

// DiskSpaceScan 扫描磁盘空间
func DiskSpaceScan(id int) (int, error) {
	// 如果 id == -1，表示全表扫描
	cond := ""
	if id != -1 {
		cond = " where id=" + strconv.Itoa(id)
	}

	items, err := DiskSpaceFindBy(cond)
	if err != nil || len(items) == 0 {
		return 0, nil
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	count := 0
	for _, i := range items {
		// 对每一个设定，获取磁盘容量，并更新数据库
		// 初始化 rpc
		client, err := rpc.DialHTTP("tcp", i.ServerIP+":8999")
		if err != nil {
			log.Printf("faild: %s - %+v", i.ServerIP, err)
			continue
		}

		// Synchronous call
		space := &server.DiskSpace{}
		err = client.Call("Real.GetDiskSpace", i.DiskName, &space)
		if err != nil {
			log.Printf("GetDiskSpace error: %s-%s-%+v", i.ServerIP, i.DiskName, err)
			continue
		}

		// 更新回数据库
		i.LastTick = now
		i.LastTotalAmt = space.TotalAmt
		i.LastFreeAmt = space.FreeAMt
		err = i.UpdateSpace()
		if err != nil {
			log.Printf("UpdateSpace error: %+v", err)
			continue
		}
		count = count + 1
	}

	return count, nil
}

// DiskSpaceResultT 邮件通知的数据结构
type DiskSpaceResultT struct {
	ID           int     `db:"ID" json:"id" form:"id"`
	ServerName   string  `db:"SRVR_NME" json:"server_name" form:"server_name"`
	ServerIP     string  `db:"SRVR_IP" json:"server_ip" form:"server_ip"`
	DiskName     string  `db:"DISK_NME" json:"disk_name" form:"disk_name"`
	WarningPoint int     `db:"WARNING_PT" json:"warning_point" form:"warning_point"`
	NoticeUser   string  `db:"NOTICE_USR" json:"notice_user" form:"notice_user"`
	Status       string  `db:"STS_IND" json:"status_ind" form:"status_ind"`
	LastTick     string  `db:"L_TICK" json:"last_tick" form:"last_tick"`
	LastTotalAmt float64 `db:"L_TTL_AMT" json:"last_totalAmt" form:"last_totalAmt"`
	LastFreeAmt  float64 `db:"L_FREE_AMT" json:"last_freeAmt" form:"last_freeAmt"`
	Usage        float64 `db:"USAGE" json:"usage"`
	Email        string  `db:"EMAIL" json:"email"`
}

// DiskSpaceMonitor 对超出预警的磁盘，生成待发送邮件
func DiskSpaceMonitor() error {
	// 先扫描磁盘
	DiskSpaceScan(-1)

	// 再生成待发送邮件
	items := []DiskSpaceResultT{}

	// 仅处理 ACTIVE 的记录，并且是有邮箱的；
	// 获取超过预警值，或者最近一次更新时间是1天前的；
	sqlCmd := `select a.ID, ifnull(a.SRVR_NME,'') SRVR_NME, ifnull(a.SRVR_IP,'') SRVR_IP, ifnull(a.DISK_NME,'') DISK_NME, ifnull(a.WARNING_PT,'') WARNING_PT, ifnull(a.NOTICE_USR,'') NOTICE_USR, ifnull(a.STS_IND,'') STS_IND, ifnull(a.L_TICK,'') L_TICK, ifnull(a.L_TTL_AMT,'') L_TTL_AMT, ifnull(a.L_FREE_AMT,'') L_FREE_AMT,
								    b.email EMAIL, 100 * (1-a.L_FREE_AMT/a.L_TTL_AMT) 'USAGE'
							 from t_diskspace a inner join t_users b on a.NOTICE_USR = b.NME
							where a.STS_IND = 'ACTIVE' and ( (100 * (1-a.L_FREE_AMT/a.L_TTL_AMT) >= a.WARNING_PT) or (a.L_TICK < date_format(date_add(now(), interval -1 day), '%Y-%m-%d') ))
						`

	// fmt.Println(sqlCmd)
	err := db.Select(&items, sqlCmd)

	if err != nil {
		return err
	}

	// 取出 distinct 的 email
	emails := make(map[string]int)
	for _, i := range items {
		emails[i.Email] = 1
	}

	if len(emails) == 0 {
		return nil
	}

	// 加载 template
	add := func(x, y int) int {
		return x + y
	}
	numberFormat := func(x float64) string {
		return fmt.Sprintf("%.2f", x)
	}

	funcs := template.FuncMap{"add": add, "numberFormat": numberFormat}
	emailTemplate := templateDir + "/diskspaceEmail.tmpl"
	tmpl, err := template.New("diskspaceEmail.tmpl").Funcs(funcs).ParseFiles(emailTemplate)
	if err != nil {
		return err
	}

	// 对每个 email 中的信息，生成一封待发送邮件
	for key := range emails {
		oneBatch := []DiskSpaceResultT{}
		for _, i := range items {
			if i.Email == key {
				oneBatch = append(oneBatch, i)
			}
		}

		if len(oneBatch) == 0 {
			continue
		}

		notice := CreateEmailNotice()
		notice.UserName = oneBatch[0].NoticeUser
		notice.SendTo = oneBatch[0].Email
		notice.Subject = fmt.Sprintf("流程通知: 磁盘容量 %d 条设定条件已满足，请尽快处理", len(oneBatch))

		content, err := genDiskSpaceEmailBody(oneBatch, tmpl)
		if err != nil {
			// 出现了错误，继续
			fmt.Printf("genDiskSpaceEmailBody: %+v\n", err)
			continue
		}
		notice.Content = content

		err = notice.Insert()
		if err != nil {
			// 插入失败，继续处理下一条
			fmt.Printf("Insert: %+v\n", err)
			continue
		}
	}

	return nil
}

func genDiskSpaceEmailBody(items []DiskSpaceResultT, tmpl *template.Template) (string, error) {
	if len(items) == 0 {
		return "", nil
	}

	// 准备 template 所需的 params
	param := struct {
		Items      []DiskSpaceResultT
		TotalCount int
	}{
		Items:      items,
		TotalCount: len(items),
	}

	// Stores the parsed template
	var buff bytes.Buffer

	// Send the parsed template to buff
	err := tmpl.Execute(&buff, param)
	if err != nil {
		return "", err
	}

	return buff.String(), nil
}
