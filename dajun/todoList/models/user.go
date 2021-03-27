package models

import (
	"pickup/dajun/todoList/config"
)

// User 用户基本信息
type User struct {
	ID          int    `db:"ID" form:"id" json:"id" `
	Code        string `db:"CDE" form:"code" json:"code" `
	Name        string `db:"NME" form:"name" json:"name"`
	EMail       string `db:"EMAIL" form:"email" json:"email"`
	Department  string `db:"DEPT" form:"department" json:"department"`
	Password    string `db:"PASSWORD" form:"password" json:"password"`
	Permissions string `db:"PERMISSIONS" form:"permissions" json:"permissions"`
	Title       string `db:"TITLE" form:"title" json:"title"`
	SvcLine     string `db:"SVC_LINE" form:"svcLine" json:"svcLine"`
}

// allUsers 缓存使用，全部的用户清单，供 FindAll 返回；
// 如果为空，则从数据库查询；否则直接返回；
// 如果增删改过，重新更新这个 list；
var allUsers []User

const (
	userSelect = `select ID, ifnull(CDE,'') CDE, ifnull(NME,'') NME, ifnull(EMAIL,'') EMAIL, ifnull(DEPT,'') DEPT, ifnull(PASSWORD,'') PASSWORD, ifnull(PERMISSIONS,'') PERMISSIONS, ifnull(TITLE,'') TITLE, ifnull(SVC_LINE,'') SVC_LINE from t_users`
)

// CreateUser 返回一个空的 User 对象
func CreateUser() User {
	return User{}
}

// FindBy 根据条件查找
func (u User) FindBy(cond string) ([]User, error) {
	cmd := userSelect + cond
	// log.Printf("user search: %s\n", cmd)

	items := []User{}
	err := db.Select(&items, cmd)
	return items, err
}

// FindAll 返回所有用户
func (u User) FindAll() ([]User, error) {
	if len(allUsers) == 0 {
		err := reloadAllUsers()
		return allUsers, err
	}
	return allUsers, nil
}

// 重新加载用户清单
func reloadAllUsers() error {
	u := CreateUser()
	var err error
	// 对 allUsers 进行赋值
	allUsers, err = u.FindBy("")
	return err
}

// FindByID 根据 ID 加载记录
func (u User) FindByID(id int) (User, error) {
	item := User{}
	cmd := userSelect + " where ID=? "
	err := db.Get(&item, cmd, id)
	return item, err
}

// FindByName 根据 Name 进行模糊查询
func (u User) FindByName(name string) ([]User, error) {
	cond := " where NME like '%" + name + "%'"
	return u.FindBy(cond)

}

// Insert 把当前对象，作为新记录，插入数据库
func (u User) Insert() error {
	// 默认密码
	u.Password = "1234"
	sqlCmd := `INSERT INTO t_users (CDE,NME,EMAIL,DEPT,PASSWORD,PERMISSIONS, TITLE, SVC_LINE) VALUES (:CDE,:NME,:EMAIL,:DEPT,:PASSWORD,:PERMISSIONS, :TITLE, :SVC_LINE)`
	_, err := db.NamedExec(sqlCmd, u)
	if err != nil {
		return err
	}

	// 重新更新用户清单
	err = reloadAllUsers()
	return err
}

// Update 将当前对象保存到数据库
func (u User) Update() error {
	sqlCmd := `update t_users set CDE=:CDE, NME=:NME,EMAIL=:EMAIL,DEPT=:DEPT, PERMISSIONS=:PERMISSIONS, TITLE = :TITLE, SVC_LINE = :SVC_LINE where id=:ID`
	_, err := db.NamedExec(sqlCmd, u)
	if err != nil {
		return err
	}

	// 重新更新用户清单
	err = reloadAllUsers()
	return err
}

// Delete 根据当前对象的 ID 删除
func (u User) Delete() error {
	// db.MustExec(userDelete, strconv.Itoa(u.ID))

	sqlCmd := `delete from t_users where ID=:ID`
	_, err := db.NamedExec(sqlCmd, u)
	if err != nil {
		return err
	}

	// 重新更新用户清单
	err = reloadAllUsers()
	return err
}

// ChangePassword 设置用户的密码
func ChangePassword(email, password string) error {
	param := make(map[string]interface{})
	param["EMAIL"] = email

	// 如果设置的密码为 1234，那么就是用最简便的方法
	if password == "1234" {
		param["PASSWORD"] = password
	} else {
		var err error
		param["PASSWORD"], err = config.HashPassword(password)
		if err != nil {
			return err
		}
	}

	sqlCmd := `update t_users set PASSWORD=:PASSWORD where EMAIL=:EMAIL`
	_, err := db.NamedExec(sqlCmd, param)
	return err
}

// ValidateLogin 根据 邮箱地址，密码，判断用户是否能登录
func ValidateLogin(email, password string) (User, error) {
	validUser := User{}
	// -1 表示用户名、密码不匹配，不能登录
	validUser.ID = -1

	// 根据 email，找到用户
	cond := ` where EMAIL = :EMAIL`
	cmd := userSelect + " " + cond
	nstmt, err := db.PrepareNamed(cmd)
	if err != nil {
		return validUser, err
	}

	param := make(map[string]interface{})
	param["EMAIL"] = email
	err = nstmt.Get(&validUser, param)
	if err != nil {
		return validUser, err
	}

	// 如果是默认密码，直接返回
	if validUser.Password == "1234" && password == "1234" {
		// 清空返回结果中的密码
		validUser.Password = "******"
		return validUser, err
	}

	// 如果用户修改过密码，则校验匹配密码
	result := config.CheckPasswordHash(password, validUser.Password)

	// 如果密码不对，设置 ID 为 -1
	if !result {
		validUser.ID = -1
	}
	// 清空返回结果中的密码
	validUser.Password = "******"
	return validUser, err
}
