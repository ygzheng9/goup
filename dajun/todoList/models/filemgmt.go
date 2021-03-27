package models

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/copier"
)

// FileMgmt 对应于每一个文件/目录
type FileMgmt struct {
	ID           int    `db:"ID" json:"id" form:"id"`
	FileName     string `db:"FILE_NME" json:"file_name" form:"file_name"`
	IsDir        string `db:"IS_DIR" json:"is_dir" form:"is_dir"`
	Depth        int    `db:"DEPTH" json:"depth" form:"depth"`
	ParentID     int    `db:"PARENT_ID" json:"parent_id" form:"parent_id"`
	UpdateUser   string `db:"UPD_USR" json:"upd_user" form:"upd_user"`
	UpdateDate   string `db:"UPD_DTE" json:"upd_dte" form:"upd_dte"`
	RefID        int    `db:"REF_ID" json:"ref_id" form:"ref_id"`
	Content      string `db:"CTENT" json:"content" form:"content"`
	Status       string `db:"STS" json:"status" form:"status"`
	AbsolutePath string `db:"ABS_PATH" json:"abs_path" form:"abs_path"`
}

// CreateFileMgmt 返回一个空对象
func CreateFileMgmt() FileMgmt {
	return FileMgmt{}
}

const (
	fileMgmtSelect = `select ID, ifnull(FILE_NME,'') FILE_NME, ifnull(IS_DIR,'') IS_DIR, ifnull(DEPTH,'') DEPTH,
											ifnull(PARENT_ID,'') PARENT_ID, ifnull(UPD_USR,'') UPD_USR, ifnull(UPD_DTE,'') UPD_DTE,
											ifnull(REF_ID,'') REF_ID, ifnull(CTENT,'') CTENT, ifnull(STS,'') STS,
											ifnull(ABS_PATH,'') ABS_PATH
										from t_file_mgmt`
)

// FindBy 返回符合条件的记录
func (f *FileMgmt) FindBy(cond string, args ...interface{}) ([]FileMgmt, error) {
	cmd := fileMgmtSelect

	if len(strings.TrimSpace(cond)) > 0 {
		cmd = cmd + " " + cond
	}

	items := []FileMgmt{}
	err := db.Select(&items, cmd, args...)
	return items, err
}

// nextLevel 根据 parentID，返回文件列表，有状态过滤
func (f *FileMgmt) nextLevel() ([]FileMgmt, error) {
	cond := ` where STS = 'A' and PARENT_ID = ` + strconv.Itoa(f.ID)
	return f.FindBy(cond)
}

// NextLevelRaw 根据 parentID 返回文件列表，无状态过滤
func (f *FileMgmt) NextLevelRaw() ([]FileMgmt, error) {
	cond := ` where PARENT_ID = ?;`
	return f.FindBy(cond, f.ID)
}

// FindHistory 查找当前文件所有的变更记录
func (f *FileMgmt) FindHistory() ([]FileMgmt, error) {
	cond := ` where REF_ID = ?
						order by FILE_NME desc;`
	return f.FindBy(cond, f.ID)
}

// LoadByID 根据 id 加载一条记录, 更新当前对象
func (f *FileMgmt) LoadByID(id int) error {
	cmd := fileMgmtSelect +
		` where ID = ?; `

	err := db.Get(f, cmd, id)
	return err
}

// FindParent 向上找父节点
func (f *FileMgmt) FindParent() (FileMgmt, error) {
	var parent FileMgmt

	err := parent.LoadByID(f.ParentID)
	return parent, err
}

// RenameFile 重命名文件, 有副作用：使用新名字更新当前对象
func (f *FileMgmt) RenameFile(newName string) error {
	// 更新文件系统
	// 取得老文件的信息
	f.LoadByID(f.ID)
	fullPath, err := f.GetFullName()
	if err != nil {
		return err
	}
	dir := path.Dir(fullPath)
	newFullName := dir + "/" + newName
	// 重命名文件
	err = os.Rename(fullPath, newFullName)

	// 更新数据库
	f.FileName = newName
	f.UpdateDate = time.Now().Format("2006-01-02 15:04:05")

	sqlCmd := `update t_file_mgmt
	             set FILE_NME=:FILE_NME, UPD_USR=:UPD_USR, UPD_DTE=:UPD_DTE
						 where id=:ID`

	_, err = db.NamedExec(sqlCmd, f)

	return err
}

// RemoveFile 为文件打删除标记；
// 对文件系统，不做操作；
func (f *FileMgmt) RemoveFile() error {
	f.UpdateDate = time.Now().Format("2006-01-02 15:04:05")

	sqlCmd := `update t_file_mgmt
							 set STS='D',
							 		 UPD_USR=:UPD_USR,
							     UPD_DTE=:UPD_DTE
						 where id=:ID`

	_, err := db.NamedExec(sqlCmd, f)

	return err
}

// Insert 保存一个新对象
func (f *FileMgmt) Insert() error {
	// 根据 struct 中的 DB tag 进行自动 named parameter
	cmd := `INSERT INTO t_file_mgmt
					(FILE_NME, IS_DIR, DEPTH, PARENT_ID, UPD_USR, UPD_DTE, REF_ID, CTENT, STS, ABS_PATH)
					VALUES
					(:FILE_NME, :IS_DIR, :DEPTH, :PARENT_ID, :UPD_USR, :UPD_DTE, :REF_ID, :CTENT, :STS, :ABS_PATH)`
	_, err := db.NamedExec(cmd, f)
	return err
}

// Update 更新现有对象，把内存数据写入db
func (f *FileMgmt) Update() error {
	f.UpdateDate = time.Now().Format("2006-01-02 15:04:05")

	sqlCmd := `update t_file_mgmt
								 set FILE_NME = :FILE_NME,
										 IS_DIR = :IS_DIR,
										 DEPTH = :DEPTH,
										 PARENT_ID = :PARENT_ID,
										 UPD_USR = :UPD_USR,
										 UPD_DTE = :UPD_DTE,
										 REF_ID = :REF_ID,
										 CTENT = :CTENT,
										 STS = :STS,
										 ABS_PATH = :ABS_PATH
							 where id=:ID`

	_, err := db.NamedExec(sqlCmd, f)

	return err
}

// BackupFile 对已存在的文件做备份
// 如果原文件正常，那么做一个版本，status=M；
// 如果原文件是删除状态，那么也做一个新版本，status=D
func (f *FileMgmt) BackupFile() error {
	err := f.LoadByID(f.ID)
	if err != nil {
		return err
	}

	count, err := f.getHistoryCnt()
	if err != nil {
		return err
	}

	// 计算文件修改的次数，格式为：0021
	next := count + 1
	suffix := "0000" + strconv.Itoa(next)
	suffix = suffix[len(suffix)-4:]

	// 复制原来的对象
	newItem := FileMgmt{}
	copier.Copy(&newItem, f)
	newFileName := f.FileName + "_" + suffix

	// 文件系统重命名，完整路径
	oldName, err := f.GetFullName()
	if err != nil {
		return err
	}
	newName := path.Dir(oldName) + "/" + newFileName
	os.Rename(oldName, newName)

	// 设置新对象引用原来的对象，重命名文件

	// 如果当前记录是 A（正常），那么备份记录就是 M（修改）；
	// 如果当前记录是 D (删除)，备件记录就是 D (删除)；
	status := "M"
	if f.Status == "D" {
		status = "D"
	}

	newItem.RefID = f.ID
	newItem.Status = status
	newItem.FileName = newFileName
	newItem.UpdateDate = time.Now().Format("2006-01-02 15:04:05")

	err = newItem.Insert()
	if err != nil {
		return err
	}

	// 备份后，原纪录状态改为 A
	f.Status = "A"
	err = f.Update()

	return err
}

// 取得当前文件变更的次数，包括：修改，删除
func (f *FileMgmt) getHistoryCnt() (int, error) {
	cmd := `select count(1) from t_file_mgmt where REF_ID = ?`
	count := 0
	err := db.QueryRowx(cmd, f.ID).Scan(&count)

	return count, err
}

// listFS 获取文件系统的文件清单
func (f *FileMgmt) listFS() []FileMgmt {
	log.Printf("path: %s", f.AbsolutePath)

	dirpath := f.AbsolutePath

	// dirpath := f.GetFullName()

	files, err := ioutil.ReadDir(dirpath)
	if err != nil {
		log.Printf("listFS error")
		return nil
	}

	// 查找 fs 中当前目录下文件
	items := make([]FileMgmt, 0, 50)
	for _, file := range files {
		var i FileMgmt
		i.ParentID = f.ID                         // 父节点
		i.Depth = f.Depth + 1                     // 深度 + 1
		i.FileName = file.Name()                  // 文件名
		i.IsDir = fmt.Sprintf("%t", file.IsDir()) // 是否是目录
		i.Status = "A"                            // 正常状态

		items = append(items, i)
	}

	return items
}

func (f *FileMgmt) loadFromFS() ([]FileMgmt, error) {
	// 数据库中没有子元素
	// 直接从 fs 中找
	f.LoadByID(f.ID)

	fullPath, err := f.GetFullName()
	f.AbsolutePath = fullPath

	items := f.listFS()

	// fs 的结果存入 db
	for _, i := range items {
		i.Insert()
	}

	// 再从 db 中查询出来
	results, err := f.nextLevel()

	return results, err
}

// FindSubs 先从 db 中找下一层，如果没有，再从 fs 中重新加载下一层
func (f *FileMgmt) FindSubs() ([]FileMgmt, error) {
	// 先从数据库中找子元素
	items, err := f.nextLevel()

	if err != nil {
		log.Printf("FindSubs error 1: %#v\n", err)
		return nil, err
	}

	if len(items) == 0 {
		// db 中没有，再从 fs 中找
		items, err = f.loadFromFS()

		if err != nil {
			log.Printf("FindSubs error 2: %#v\n", err)
			return nil, err
		}
	}

	return items, nil
}

// GetFullName 从本级开始，逐层往上找，一直到根节点，拼出完整路径
func (f *FileMgmt) GetFullName() (string, error) {
	if f.Depth == 0 {
		return f.AbsolutePath, nil
	}

	ps := make([]string, 0, f.Depth)
	ps = append(ps, f.FileName)

	curr := f
	i := 0
	for {
		p, err := curr.FindParent()
		if err != nil {
			return "", err
		}

		if p.Depth == 0 {
			ps = append(ps, p.AbsolutePath)
			break
		}

		ps = append(ps, p.FileName)

		curr = &p

		i = i + 1
		if i > f.Depth {
			break
		}
	}

	// 拼成完整路径
	result := ps[len(ps)-1]
	for j := len(ps) - 2; j >= 0; j-- {
		result = result + "/" + ps[j]
	}

	return result, nil
}
