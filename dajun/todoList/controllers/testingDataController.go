package controllers

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"pickup/dajun/todoList/models"
	"strconv"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"

	"github.com/gin-gonic/gin"
)

// TestingDataUpload 上载测试数据
func TestingDataUpload(c *gin.Context) {
	// 获取上载的文件
	file, _, err := c.Request.FormFile("file")
	defer file.Close()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("FormFile error: %+v\n", err),
		})
		return
	}

	// 把 file 转换成成 []byte
	fileBuf := bytes.NewBuffer(nil)
	if _, err = io.Copy(fileBuf, file); err != nil {
		fmt.Println("can not read file")

		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err2,
			"message": fmt.Sprintf("buffer copy: %+v\n", err),
		})

		return
	}

	// utf8Buf, err := GbkToUtf8(fileBuf.Bytes())
	// if err != nil {
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"rtnCode": err3,
	// 		"message": fmt.Sprintf("utf8 error: %+v\n", err),
	// 	})

	// 	return
	// }

	_, err = sendTestingData(fileBuf.Bytes())
	// err = models.TestingDataParse(utf8Buf)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err3,
			"message": fmt.Sprintf("TestingDataParse: %+v\n", err),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rtnCode": ok,
		// "message": msg,
	})
}

func sendTestingData(input []byte) (string, error) {
	// connect to this socket, 本地
	addr := "127.0.0.1:8090"
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	// 发送信息给 server，一次性发送
	_, err = conn.Write(input)
	if err != nil {
		return "", err
	}

	// 从 tcp server 取得返回状态
	// 获取反馈 server，一次性读取
	resBuf := make([]byte, 1024)
	reqLen, err := conn.Read(resBuf)
	if err != nil {
		return "", err
	}

	msg := string(resBuf[:reqLen])
	return msg, nil
}

// GbkToUtf8 GBK 到 UTF-8 的转换
func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

// Utf8ToGbk UTF-8 到 GBK 的转换
func Utf8ToGbk(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

// TestingDataSave 把测试数据保存到数据库
func TestingDataSave(c *gin.Context) {
	param := struct {
		Data string `json:"data"`
	}{}

	// post 过来的 json，绑定到 对象上；
	err := c.BindJSON(&param)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("Update BindJson error: %+v\n", err),
		})
		return
	}

	err = models.TestingDataParse([]byte(param.Data))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err2,
			"message": fmt.Sprintf("TestingDataParse: %+v\n", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rtnCode": ok,
	})
}

// TestingDataSearch 根据条件查找头信息
func TestingDataSearch(c *gin.Context) {
	param := struct {
		Product   string `json:"product"`
		ProductSN string `json:"productSN"`
		Start     string `json:"start_dt"`
		End       string `json:"end_dt"`
	}{}

	// post 过来的 json，绑定到 对象上；
	err := c.BindJSON(&param)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("BindJson error: %+v\n", err),
		})
		return
	}

	// 拼接查询条件
	var cond bytes.Buffer
	cond.WriteString(" where 1 = 1 ")

	if param.Product != "" {
		cond.WriteString(" and PRODCT like '%" + param.Product + "%' ")
	}
	if param.ProductSN != "" {
		cond.WriteString(" and PRODCT_SN like '%" + param.ProductSN + "%' ")
	}
	if param.Start != "" {
		cond.WriteString(" and (TESTING_DTE >= '" + param.Start + "' and TESTING_DTE <= '" + param.End + "') ")
	}

	items, err := models.TestingDataHeadFind(cond.String())
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err2,
			"message": fmt.Sprintf("TestingDataHeadFind error: %+v\n", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rtnCode": ok,
		"items":   items,
	})
}

// TestingDataDelete 删除测试数据
func TestingDataDelete(c *gin.Context) {
	// 从 url path 中提取 id
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("path param error: %+v\n", err),
		})
		return
	}

	err = models.TestingDataHeadDelete(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err2,
			"message": fmt.Sprintf("TestingDataHeadDelete error: %+v\n", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rtnCode": ok,
	})
}

// TestingDataLoadDeatail 查询明细测试结果
func TestingDataLoadDeatail(c *gin.Context) {
	// 从 url path 中提取 id
	headID, err := strconv.Atoi(c.Param("headid"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("path param error: %+v\n", err),
		})
		return
	}

	items, err := models.TestingDataItemFind(headID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err2,
			"message": fmt.Sprintf("TestingDataItemFind error: %+v\n", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rtnCode": ok,
		"items":   items,
	})

}
