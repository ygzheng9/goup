package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"pickup/dajun/todoList/models"
)

// MatWLGetAll 取得监控的磁盘空间列表
func MatWLGetAll(c *gin.Context) {
	items, err := models.MatWLFindAll()

	if err != nil {
		c.JSON(200, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("FindAll error: %+v\n", err),
		})
		return
	}

	boms, err := models.MatWLLoadBOM(items)
	if err != nil {
		c.JSON(200, gin.H{
			"rtnCode": err2,
			"message": fmt.Sprintf("MatWLLoadBOM error: %+v\n", err),
		})
		return
	}

	c.JSON(200, gin.H{
		"rtnCode": ok,
		"items":   items,
		"boms":    boms,
	})
}

// MatWLFindByID 根据 ID 查找
func MatWLFindByID(c *gin.Context) {
	// 从 url path 中提取 id
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("path param error: %+v\n", err),
		})
		return
	}

	item, err := models.MatWLFindByID(id)
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

// MatWLInsert 当前对象插入到数据库
func MatWLInsert(c *gin.Context) {
	// 创建一个对象，用以绑定 post 过来的 json
	param := models.MaterialWatchItem{}

	// post 过来的 json，绑定到 对象上；
	err := c.BindJSON(&param)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("Update BindJson error: %+v\n", err),
		})
		return
	}

	item, err := models.MatWLInsert(param)
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

// MatWLBatchInsert 批量插入到数据库
func MatWLBatchInsert(c *gin.Context) {
	param := struct {
		User   string   `json:"user"`
		Remark string   `json:"remark"`
		Items  []string `json:"items"`
	}{}

	// post 过来的 json，绑定到 对象上；
	err := c.BindJSON(&param)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("MatWLBatchInsert BindJson error: %+v\n", err),
		})
		return
	}

	failed := 0
	errMsg := ""

	for _, v := range param.Items {
		mat := models.MaterialWatchItem{}
		mat.CreateUser = param.User
		mat.Remark = param.Remark
		mat.InvCode = v

		_, err := models.MatWLInsert(mat)
		if err != nil {
			failed++

			errMsg = errMsg + fmt.Sprintf("%s: %+v\n", v, err)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"rtnCode": ok,
		"total":   len(param.Items),
		"failed":  failed,
		"msg":     errMsg,
	})
}

// MatWLDelete 从数据库删除当前对象
func MatWLDelete(c *gin.Context) {
	// 按照 id 查询
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("path param error: %+v\n", err),
		})
		return
	}

	item := models.MaterialWatchItem{}
	item.ID = id
	err = models.MatWLDelete(item)
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
