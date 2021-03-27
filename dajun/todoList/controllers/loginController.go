package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"pickup/dajun/todoList/models"
)

// ValidateLogin 验证用户
func ValidateLogin(c *gin.Context) {
	type paramT struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var param paramT

	// 按照 id 查询用户
	err := c.BindJSON(&param)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("BindJSON error: %+v\n", err),
		})
		return
	}

	user, err := models.ValidateLogin(param.Email, param.Password)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err2,
			"message": fmt.Sprintf("ValidateLogin error: %+v\n", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rtnCode": ok,
		"item":    user,
	})
}
