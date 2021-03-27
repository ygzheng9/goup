package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"pickup/dajun/todoList/models"
)

// 操作日志

// ActivityLogFindByRefID 根据 ref 取得操作日志
func ActivityLogFindByRefID(c *gin.Context) {
	param := struct {
		RefID   int    `json:"ref_id"`
		RefType string `json:"ref_type"`
	}{}

	// post 过来的 json，绑定到 对象上；
	err := c.BindJSON(&param)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("ActivityLogFindByRefID BindJson error: %+v\n", err),
		})
		return
	}

	items, err := models.ActivityLogFindByRefID(param.RefID, param.RefType)

	if err != nil {
		c.JSON(200, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("ActivityLogFindByRefID error: %+v\n", err),
		})
		return
	}

	c.JSON(200, gin.H{
		"rtnCode": ok,
		"items":   items,
	})
}
