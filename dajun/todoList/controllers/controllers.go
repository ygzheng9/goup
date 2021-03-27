package controllers

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	ok   = 0
	err1 = 1
	err2 = 2
	err3 = 3
	err4 = 4
	err5 = 5
	err6 = 6
	err7 = 7
	err8 = 8
	err9 = 9
)

// GeneratePassword 生成密码的密文
func GeneratePassword(c *gin.Context) {
	type paramT struct {
		Plain string `json:"plain" form:"plain"`
	}

	var param paramT
	err := c.BindJSON(&param)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("GeneratePassword error: %+v\n", err),
		})
		return
	}

	// 生成密文, base64编码
	encodeString := base64.StdEncoding.EncodeToString([]byte(param.Plain))
	c.JSON(http.StatusOK, gin.H{
		"rtnCode": ok,
		"value":   encodeString,
	})
}

// TestReload 测试用
func TestReload(c *gin.Context) {
	t := time.Now()
	msg := "test...." + t.Format("2006-01-02 15:04:05")
	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": msg})

	// can not do the second json
	// c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "second json"})

}

// TestPost 测试用
func TestPost(c *gin.Context) {
	// userName := c.PostForm("userName")
	// comment := c.PostForm("comment")

	type param struct {
		UserName string `form:"userName" json:"userName"`
		Comment  string `form:"comment" json:"comment"`
	}

	p := param{}

	err := c.BindJSON(&p)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "bad", "message": err})
		return
	}

	t := time.Now()
	msg := p.UserName + " " + p.Comment + " " + t.Format("2006-01-02 15:04:05")
	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": msg, "userName": p.UserName, "comment": p.Comment})
}
