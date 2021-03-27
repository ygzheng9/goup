package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"pickup/dajun/todoList/models"
)

// BpmProcessNodeGetAll 取得监控的磁盘空间列表
func BpmProcessNodeGetAll(c *gin.Context) {
	items, err := models.BpmProcessNodeFindAll()

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

// BpmProcessNodeFindByID 根据 ID 查找
func BpmProcessNodeFindByID(c *gin.Context) {
	// 从 url path 中提取 id
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("path param error: %+v\n", err),
		})
		return
	}

	item, err := models.BpmProcessNodeFindByID(id)
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

// BpmProcessNodeCheckRule 检查节点规则
func BpmProcessNodeCheckRule(c *gin.Context) {
	// 创建一个对象，用以绑定 post 过来的 json
	param := models.BpmProcessNode{}

	// post 过来的 json，绑定到 对象上；
	err := c.BindJSON(&param)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("BindJson error: %+v", err),
		})
		return
	}

	isValid := models.BpmProcessNodeCheckRule(param.BizRule)
	if !isValid {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err2,
			"message": fmt.Sprintf("invalid rule: %s", param.BizRule),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rtnCode": ok,
	})
}

// BpmProcessNodeInsert 当前对象插入到数据库
func BpmProcessNodeInsert(c *gin.Context) {
	// 创建一个对象，用以绑定 post 过来的 json
	param := models.BpmProcessNode{}

	// post 过来的 json，绑定到 对象上；
	err := c.BindJSON(&param)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf(" BpmProcessNodeInsert BindJson error: %+v\n", err),
		})
		return
	}

	item, err := models.BpmProcessNodeInsert(param)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err2,
			"message": fmt.Sprintf("BpmProcessNodeInsert error: %+v\n", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rtnCode": ok,
		"item":    item,
	})
}

// BpmProcessNodeUpdate 更新当前对象
func BpmProcessNodeUpdate(c *gin.Context) {
	// 创建一个对象，用以绑定 post 过来的 json
	param := models.BpmProcessNode{}

	// post 过来的 json，绑定到 对象上；
	err := c.BindJSON(&param)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("Update BindJson error: %+v\n", err),
		})
		return
	}

	item, err := models.BpmProcessNodeUpdate(param)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err2,
			"message": fmt.Sprintf("Update error: %+v\n", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rtnCode": ok,
		"item":    item,
	})
}

// BpmProcessNodeDeleteByID 根据 ID 删除
func BpmProcessNodeDeleteByID(c *gin.Context) {
	// 从 url path 中提取 id
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("path param error: %+v\n", err),
		})
		return
	}

	err = models.BpmProcessNodeDelete(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err2,
			"message": fmt.Sprintf("FindByID error: %+v\n", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rtnCode": ok,
	})
}
