package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"pickup/dajun/todoList/models"
)

// EventController Event 相关的操作
type EventController struct{}

// CreateEventController 返回一个空对象
func CreateEventController() EventController {
	return EventController{}
}

// FindAll 取得事项清单
func (d EventController) FindAll(c *gin.Context) {
	item := models.CreateEvent()
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
func (d EventController) FindByID(c *gin.Context) {
	// 从 url path 中提取 id
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("path param error: %+v\n", err),
		})
		return
	}

	item := models.CreateEvent()
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

// Insert 当前对象插入到数据库
func (d EventController) Insert(c *gin.Context) {
	// 创建一个对象，用以绑定 post 过来的 json
	param := models.CreateEvent()

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
func (d EventController) Update(c *gin.Context) {
	// 创建一个对象，用以绑定 post 过来的 json
	param := models.CreateEvent()

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
func (d EventController) Delete(c *gin.Context) {
	// 按照 id 查询
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("path param error: %+v\n", err),
		})
		return
	}

	item := models.CreateEvent()
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

// NotifyOpen 对一个事件下的待办，发通知
func (d EventController) NotifyOpen(c *gin.Context) {
	// 按照 id 查询
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("ScanEvent param error: %+v\n", err),
		})
		return
	}

	// 生成待发送邮件
	err = models.ScanEvent(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("ScanEvent error: %+v\n", err),
		})
		return
	}

	// 发送邮件
	err = models.SendAllEmails()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err2,
			"message": fmt.Sprintf("SendAllEmails error: %+v\n", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rtnCode": ok,
	})
}

// CloseEvent 关闭当前事件
func CloseEvent(c *gin.Context) {
	// 按照 id 查询
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("CloseEvent param: %+v\n", err),
		})
		return
	}

	err = models.CloseEvent(id, true)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("CloseEvent: %+v\n", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rtnCode": ok,
	})
}

// FindEventByParam 根据参数查询事件
func FindEventByParam(c *gin.Context) {
	// 解析参数
	type paramT struct {
		Start  string `json:"start_dt"`
		End    string `json:"end_dt"`
		Status string `json:"status"`
		Cond   string `json:"cond"`
	}
	param := paramT{}

	// post 过来的 json，绑定到 对象上；
	err := c.BindJSON(&param)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("FindEventByParam BindJson error: %+v\n", err),
		})
		return
	}

	// 根据参数查询
	items, err := models.FindEventByParam(param.Start, param.End, param.Status, param.Cond)

	if err != nil {
		c.JSON(200, gin.H{
			"rtnCode": err2,
			"message": fmt.Sprintf("FindEventByParam error: %+v\n", err),
		})
		return
	}

	c.JSON(200, gin.H{
		"rtnCode": ok,
		"items":   items,
	})
}
