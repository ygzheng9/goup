package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"pickup/dajun/todoList/models"
)

// LaborClaimFindByWk 查找一周的记录
func LaborClaimFindByWk(c *gin.Context) {
	type paramT struct {
		UserName string `json:"userName"`
		Monday   string `json:"monday"`
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

	monday, err := time.Parse("2006-01-02", param.Monday)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err2,
			"message": fmt.Sprintf("time.Parse: %+v\n", err),
		})
		return
	}

	// 加载一周的数据
	sunday := monday.AddDate(0, 0, 6).Format("2006-01-02")
	items, err := models.LaborClaimFindByWk(param.UserName, param.Monday, sunday)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err3,
			"message": fmt.Sprintf("LaborClaimFindByWk: %+v\n", err),
		})
		return
	}

	// 上个星期填过的项目
	from := monday.AddDate(0, 0, -6).Format("2006-01-02")
	projList, err := models.LaborClaimFindProj(param.UserName, from, sunday)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err4,
			"message": fmt.Sprintf("LaborClaimFindProj %+v\n", err),
		})
		return
	}

	c.JSON(200, gin.H{
		"rtnCode":  ok,
		"items":    items,
		"projList": projList,
	})
}

// LaborClaimBatchSave 批量保存
func LaborClaimBatchSave(c *gin.Context) {
	type itemT struct {
		ProjCode  string `json:"projCode"`
		BizDate   string `json:"bizDate"`
		HourCount int    `json:"hourCount"`
		Remark    string `json:"remark"`
	}

	type paramT struct {
		UserName string  `json:"userName"`
		Monday   string  `json:"monday"`
		Items    []itemT `json:"items"`
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

	monday, err := time.Parse("2006-01-02", param.Monday)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err2,
			"message": fmt.Sprintf("time.Parse: %+v\n", err),
		})
		return
	}

	// 删除一周的数据
	sunday := monday.AddDate(0, 0, 6).Format("2006-01-02")
	err = models.LaborClaimDeleteWk(param.UserName, param.Monday, sunday)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err3,
			"message": fmt.Sprintf("LaborClaimDeleteWk: %+v\n", err),
		})
		return
	}

	// 循环插入
	for _, v := range param.Items {
		labor := models.LaborClaim{}
		labor.BizDate = v.BizDate
		labor.ProjCode = v.ProjCode
		labor.HourCount = v.HourCount
		labor.Remark = v.Remark
		labor.UserName = param.UserName
		labor.UpdateUser = param.UserName

		err = models.LaborClaimInsert(labor)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"rtnCode": err4,
				"message": fmt.Sprintf("LaborClaimInsert: %+v\n", err),
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"rtnCode": ok,
	})
}

// LaborClaimSearchByPeriod 期间内搜索
func LaborClaimSearchByPeriod(c *gin.Context) {
	param := struct {
		Start    string `json:"start"`
		End      string `json:"end"`
		UserName string `json:"user"`
	}{}

	err := c.BindJSON(&param)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("BindJSON: %+v\n", err),
		})
		return
	}

	items, err := models.LaborClaimSearchByPeriod(param.UserName, param.Start, param.End)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err2,
			"message": fmt.Sprintf("LaborClaimSearchByPeriod: %+v\n", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rtnCode": ok,
		"items":   items,
	})
}
