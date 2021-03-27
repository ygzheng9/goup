package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"pickup/dajun/todoList/models"
)

// MilestoneGetAll 取得列表
func MilestoneGetAll(c *gin.Context) {
	items, err := models.MilestoneFindAll()

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

// MilestoneFindByContract 合同下的付款计划
func MilestoneFindByContract(c *gin.Context) {
	// 从 url path 中提取 id
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("path param error: %+v\n", err),
		})
		return
	}

	items, err := models.MilestoneFindByContract(id)

	if err != nil {
		c.JSON(200, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("MilestoneFindByContract error: %+v\n", err),
		})
		return
	}

	c.JSON(200, gin.H{
		"rtnCode": ok,
		"items":   items,
	})
}

// MilestoneFindByID 根据 ID 查找
func MilestoneFindByID(c *gin.Context) {
	// 从 url path 中提取 id
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("path param error: %+v\n", err),
		})
		return
	}

	item, err := models.MilestoneFindByID(id)
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

// MilestoneInsert 当前对象插入到数据库
func MilestoneInsert(c *gin.Context) {
	// 创建一个对象，用以绑定 post 过来的 json
	param := models.Milestone{}

	// post 过来的 json，绑定到 对象上；
	err := c.BindJSON(&param)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("Update BindJson error: %+v\n", err),
		})
		return
	}

	item, err := models.MilestoneInsert(param)
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

// MilestoneUpdate 更新当前对象
func MilestoneUpdate(c *gin.Context) {
	// 创建一个对象，用以绑定 post 过来的 json
	param := models.Milestone{}

	// post 过来的 json，绑定到 对象上；
	err := c.BindJSON(&param)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("Update BindJson error: %+v\n", err),
		})
		return
	}

	item, err := models.MilestoneUpdate(param)
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

// MilestoneDeleteByID 根据 ID 查找
func MilestoneDeleteByID(c *gin.Context) {
	// 从 url path 中提取 id
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("path param error: %+v\n", err),
		})
		return
	}

	i := models.Milestone{}
	i.ID = id
	err = models.MilestoneDelete(i)
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
