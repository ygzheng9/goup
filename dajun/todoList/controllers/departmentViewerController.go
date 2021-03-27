package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"pickup/dajun/todoList/models"
)

// DeptViewerFindAll 加载全部关系
func DeptViewerFindAll(c *gin.Context) {
	items, err := models.DeptViewerFindAll()

	if err != nil {
		c.JSON(200, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("DeptViewerFindAll: %+v\n", err),
		})
		return
	}

	c.JSON(200, gin.H{
		"rtnCode": ok,
		"items":   items,
	})
}

// DeptViewerBatchAdd 批量建立关系
func DeptViewerBatchAdd(c *gin.Context) {
	param := struct {
		DeptName  string   `json:"deptName"`
		UserNames []string `json:"userNames"`
	}{}

	err := c.BindJSON(&param)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("BindJson error: %+v\n", err),
		})
		return
	}

	err = models.DeptViewerBatchAdd(param.DeptName, param.UserNames)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err2,
			"message": fmt.Sprintf("DeptViewerBatchAdd error: %+v\n", err),
		})
		return
	}

	items, err := models.DeptViewerFindAll()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err3,
			"message": fmt.Sprintf("DeptViewerFindAll error: %+v\n", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rtnCode": ok,
		"items":   items,
	})
}

// DeptViewerRemove 移除关联
func DeptViewerRemove(c *gin.Context) {
	param := struct {
		DeptName string `json:"deptName"`
		UserName string `json:"userName"`
	}{}

	err := c.BindJSON(&param)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("BindJson error: %+v\n", err),
		})
		return
	}

	err = models.DeptViewerRemove(param.DeptName, param.UserName)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err2,
			"message": fmt.Sprintf("DeptViewerRemove error: %+v\n", err),
		})
		return
	}

	items, err := models.DeptViewerFindAll()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err3,
			"message": fmt.Sprintf("DeptViewerFindAll error: %+v\n", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rtnCode": ok,
		"items":   items,
	})
}
