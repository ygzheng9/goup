package controllers

import (
	"fmt"
	"net/http"
	"pickup/dajun/todoList/models"
	"time"

	"github.com/gin-gonic/gin"
)

// ListEmails 显示邮件列表
func ListEmails(c *gin.Context) {
	// 参数暂时没用
	type paramT struct {
		Cond string `json:"cond"`
	}
	var param paramT

	err := c.BindJSON(&param)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("BindJSON error: %+v\n", err),
		})
		return
	}

	// 过滤留到客户端做，这里默认返回 90 天之内的邮件
	start := time.Now().AddDate(0, 0, -90).Format("2006-01-02")
	where := " where CRE_DTE >= '" + start + "' order by CRE_DTE desc"

	item := models.CreateEmailNotice()
	items, err := item.FindBy(where)
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

// ListMatNotice 显示物料提醒
func ListMatNotice(c *gin.Context) {
	items, err := models.ListMatNotice()
	if err != nil {
		c.JSON(200, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("ListMatNotice error: %+v\n", err),
		})
		return
	}

	c.JSON(200, gin.H{
		"rtnCode": ok,
		"items":   items,
	})
}
