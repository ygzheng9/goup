package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"pickup/dajun/todoList/models"
)

// ProjectFindAll 取得项目清单
func ProjectFindAll(c *gin.Context) {
	items, err := models.ProjectFindBy(" order by PROJ_GRP, PROJ_CDE ")

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("ProjectFindBy: %+v\n", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rtnCode": ok,
		"items":   items,
	})
}

// ProjectInsert 当前对象插入到数据库
func ProjectInsert(c *gin.Context) {
	// 创建一个对象，用以绑定 post 过来的 json
	param := models.Project{}

	// post 过来的 json，绑定到 对象上；
	err := c.BindJSON(&param)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("Insert error: %+v\n", err),
		})
		return
	}

	err = models.ProjectInsert(param)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err2,
			"message": fmt.Sprintf("Insert error: %+v\n", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rtnCode": ok,
	})
}

// ProjectUpdate 更新当前对象
func ProjectUpdate(c *gin.Context) {
	// 创建一个对象，用以绑定 post 过来的 json
	param := models.Project{}

	// post 过来的 json，绑定到 对象上；
	err := c.BindJSON(&param)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("Update error: %+v\n", err),
		})
		return
	}

	updOne, err := models.ProjectUpdate(param)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err2,
			"message": fmt.Sprintf("Update error: %+v\n", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rtnCode": ok,
		"item":    updOne,
	})
}

// ProjectDelete 从数据库删除当前对象
func ProjectDelete(c *gin.Context) {
	// 按照 id 查询
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("path error: %+v\n", err),
		})
		return
	}

	err = models.ProjectDelete(id)
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
