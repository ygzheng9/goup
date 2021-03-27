package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"pickup/dajun/todoList/models"
)

// DailyTaskGetAll 取得监控的磁盘空间列表
func DailyTaskGetAll(c *gin.Context) {
	items, err := models.DailyTaskFindAll()

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

// DailyTaskSearchBy 取得监控的磁盘空间列表
func DailyTaskSearchBy(c *gin.Context) {
	param := struct {
		Cond  string `json:"cond"`
		Start string `json:"start_dt"`
		End   string `json:"end_dt"`
	}{}

	// post 过来的 json，绑定到 对象上；
	err := c.BindJSON(&param)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("DailyTaskSearchBy BindJson error: %+v\n", err),
		})
		return
	}

	// 拼接查询条件
	cond := ` where (BIZ_DTE >='` + param.Start + `' and BIZ_DTE <= '` + param.End + `')`
	if len(param.Cond) > 0 {
		cond = cond + `and (USR_NME like '%` + param.Cond + `%' or WORK_RMK like '%` + param.Cond + `%')`
	}
	cond = cond + ` order by BIZ_DTE, USR_NME`

	items, err := models.DailyTaskFindBy(cond)
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

// DailyTaskFindByID 根据 ID 查找
func DailyTaskFindByID(c *gin.Context) {
	// 从 url path 中提取 id
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("path param error: %+v\n", err),
		})
		return
	}

	item, err := models.DailyTaskFindByID(id)
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

// DailyTaskInsert 当前对象插入到数据库
func DailyTaskInsert(c *gin.Context) {
	// 创建一个对象，用以绑定 post 过来的 json
	param := models.DailyTask{}

	// post 过来的 json，绑定到 对象上；
	err := c.BindJSON(&param)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("Update BindJson error: %+v\n", err),
		})
		return
	}

	item, err := models.DailyTaskInsert(param)
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

// DailyTaskUpdate 更新当前对象
func DailyTaskUpdate(c *gin.Context) {
	// 创建一个对象，用以绑定 post 过来的 json
	param := models.DailyTask{}

	// post 过来的 json，绑定到 对象上；
	err := c.BindJSON(&param)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("Update BindJson error: %+v\n", err),
		})
		return
	}

	item, err := models.DailyTaskUpdate(param)
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
