package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"pickup/dajun/todoList/models"
)

// SoftInstController Feedback 相关的操作
type SoftInstController struct{}

// CreateSoftInstController 返回一个空对象
func CreateSoftInstController() SoftInstController {
	return SoftInstController{}
}

// FindAll 取得清单
func (d SoftInstController) FindAll(c *gin.Context) {
	item := models.CreateSoftInst()
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

// FindByID 根据 ID 查找
func (d SoftInstController) FindByID(c *gin.Context) {
	// 从 url path 中提取 id
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("path param error: %+v\n", err),
		})
		return
	}

	item := models.CreateSoftInst()
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

// FindBy 按照查询条件查找
func (d SoftInstController) FindBy(c *gin.Context) {
	// 定义数据结构，接收 post 过来的参数
	type param struct {
		Cond string `json:"cond" form:"cond"`
	}

	var p param

	// post 过来的 json，绑定到 对象上；
	err := c.BindJSON(&p)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("FindBy BindJson error: %+v\n", err),
		})
		return
	}

	// 生成查询条件
	cmd := " where SOFT_NME like '%" + p.Cond + "%'"

	item := models.CreateSoftInst()
	items, err := item.FindBy(cmd)

	if err != nil {
		c.JSON(200, gin.H{
			"rtnCode": err2,
			"message": fmt.Sprintf("FindBy error: %+v\n", err),
		})
		return
	}

	c.JSON(200, gin.H{
		"rtnCode": ok,
		"items":   items,
	})
}

// Insert 当前对象插入到数据库
func (d SoftInstController) Insert(c *gin.Context) {
	// 创建一个对象，用以绑定 post 过来的 json
	param := models.CreateSoftInst()

	// post 过来的 json，绑定到 对象上；
	err := c.BindJSON(&param)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("Update BindJson error: %+v\n", err),
		})
		return
	}

	err = param.Insert()
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

// Update 更新当前对象
func (d SoftInstController) Update(c *gin.Context) {
	// 创建一个对象，用以绑定 post 过来的 json
	param := models.CreateSoftInst()

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

// Delete 从数据库删除当前对象
func (d SoftInstController) Delete(c *gin.Context) {
	// 按照 id 查询
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("path param error: %+v\n", err),
		})
		return
	}

	item := models.CreateSoftInst()
	item.ID = id
	err = item.Delete()
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
