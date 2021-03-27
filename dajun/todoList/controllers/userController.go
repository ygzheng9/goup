package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"pickup/dajun/todoList/models"
)

// UserController user 相关的操作
type UserController struct{}

// CreateUserController 返回一个空对象
func CreateUserController() UserController {
	return UserController{}
}

// FindAll 取得用户清单
func (u UserController) FindAll(c *gin.Context) {
	item := models.User{}
	items, err := item.FindAll()

	if err != nil {
		c.JSON(200, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("FindAll error: %+v\n", err),
		})
		return
	}

	c.JSON(200, gin.H{
		"rtnCode": ok,
		"items":   items,
	})
}

// Search 取得用户清单
func (u UserController) Search(c *gin.Context) {
	// fetch post 过来的是 json，不是 FormData(), 所以 BindJSON
	type searchParam struct {
		User string `json:"user"`
		Dept string `json:"dept"`
	}

	p := searchParam{}
	err := c.BindJSON(&p)
	if err != nil {
		c.JSON(200, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("BindJSON error: %+v\n", err),
		})
		return
	}

	// 拼接查询条件
	cond := " where 1 = 1 "
	if len(p.User) > 0 {
		cond = cond + ` and ( (NME like '%` + p.User + `%') or (CDE like '%` + p.User + `%') or (EMAIL like '%` + p.User + `%') )`
	}

	if len(p.Dept) > 0 {
		// cond = cond + ` and ( (DEPT like '%` + p.Dept + `%') )`

		// %% 是 literal value
		cond = fmt.Sprintf("%s and ( (DEPT like '%%%s%%') ) ", cond, p.Dept)
	}

	item := models.User{}
	items, err := item.FindBy(cond)

	if err != nil {
		c.JSON(200, gin.H{
			"rtnCode": err2,
			"message": fmt.Sprintf("FindAll error: %+v\n", err),
		})
		return
	}

	c.JSON(200, gin.H{
		"rtnCode": ok,
		"items":   items,
	})
}

// FindByID 根据用户 ID 查找
func (u UserController) FindByID(c *gin.Context) {
	// 从 url path 中提取 id
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("path param error: %+v\n", err),
		})
		return
	}

	item := models.User{}
	item, err = item.FindByID(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err2,
			"message": fmt.Sprintf("FindByID error: %+v\n", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rtnCode": ok,
		"item":    item,
	})
}

// FindByName 根据 name 模糊查找
func (u UserController) FindByName(c *gin.Context) {
	name := c.Query("name")

	// 没有输入检索条件，返回空列表
	if len(name) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": "未输入参数",
		})
		return
	}

	item := models.User{}
	items, err := item.FindByName(name)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err2,
			"message": fmt.Sprintf("FindByName error: %+v\n", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rtnCode": ok,
		"items":   items,
	})
}

// Insert 插入新对象到数据库
func (u UserController) Insert(c *gin.Context) {
	// 创建一个对象，用以绑定 post 过来的 json
	param := models.CreateUser()

	// post 过来的 json，绑定到 对象上；
	err := c.BindJSON(&param)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("Update BindJson error: %+v\n", err),
		})
		return
	}

	// 模拟：特定用户，不能保存
	if strings.Compare(param.Name, "god") == 0 {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err2,
			"message": "Can not save god",
		})
		return
	}

	err = param.Insert()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err3,
			"message": fmt.Sprintf("Update error: %+v\n", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rtnCode": ok,
	})
}

// Update 更新当前对象
func (u UserController) Update(c *gin.Context) {
	// 创建一个对象，用以绑定 post 过来的 json
	param := models.CreateUser()

	// post 过来的 json，绑定到 对象上；
	err := c.BindJSON(&param)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("Update BindJson error: %+v\n", err),
		})
		return
	}

	err = param.Update()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err2,
			"message": fmt.Sprintf("Update error: %+v\n", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rtnCode": ok,
	})
}

// Delete 删除用户
func (u UserController) Delete(c *gin.Context) {
	// 按照 id 查询用户
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("path param error: %+v\n", err),
		})
		return
	}

	item := models.CreateUser()
	item.ID = id
	err = item.Delete()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("Delete error: %+v\n", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rtnCode": ok,
	})
}

// ChangePassword 修改用户密码
func ChangePassword(c *gin.Context) {
	type paramT struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var param paramT

	// 按照 id 查询用户
	err := c.BindJSON(&param)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("BindJSON error: %+v\n", err),
		})
		return
	}

	err = models.ChangePassword(param.Email, param.Password)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err2,
			"message": fmt.Sprintf("Delete error: %+v\n", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rtnCode": ok,
	})
}
