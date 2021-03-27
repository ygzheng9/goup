package controllers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"

	"pickup/dajun/todoList/models"
)

// UploadController Feedback 相关的操作
type UploadController struct{}

const (
	// 修改成 web 目录下，而不是绝对路径'
	uploadDir = "./attachments"
)

// CreateUploadController 返回一个空对象
func CreateUploadController() UploadController {
	return UploadController{}
}

// FindAll 取得事项清单
func (d UploadController) FindAll(c *gin.Context) {
	item := models.CreateUploadItem()
	items, err := item.FindAll()

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

// FindByID 根据 ID 查找
func (d UploadController) FindByID(c *gin.Context) {
	// 从 url path 中提取 id
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("path param error: %+v\n", err),
		})
		return
	}

	item := models.CreateUploadItem()
	item, err = item.FindByID(id)
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

// FindByRef 根据 ref 查找列表
func (d UploadController) FindByRef(c *gin.Context) {
	// 取得查询参数
	item := models.CreateUploadItem()
	err := c.BindJSON(&item)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("FindByRef BindJson error: %+v\n", err),
		})
		return
	}

	items, err := item.FindByRef()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err2,
			"message": fmt.Sprintf("FindByRef error: %+v\n", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rtnCode": ok,
		"items":   items,
	})
}

// Insert 当前对象插入到数据库
func (d UploadController) Insert(c *gin.Context) {
	// 创建一个对象，接收参数
	param := models.CreateUploadItem()
	// 这里是上载文件，是 form 的 post，不是 json；
	err := c.Bind(&param)
	if err != nil {
		fmt.Printf("Insert Bind: %#v", err)
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("Insert Bind: %#v", err),
		})
		return
	}

	// 获取上载的文件
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err2,
			"message": fmt.Sprintf("FormFile: %+v", err),
		})
		return
	}

	// 在数据库中做记录，取得插入记录的 id
	param.FileName = header.Filename
	id, err := param.Insert()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err3,
			"message": fmt.Sprintf("Insert: %+v\n", err),
		})
		return
	}
	// 服务器上文件的全路径
	fileLoc := fmt.Sprintf("%s/at_%d", uploadDir, id)
	log.Printf("fullpath: %s\n", fileLoc)

	// 把文件上载的服务器
	out, err := os.Create(fileLoc)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err4,
			"message": fmt.Sprintf("Create: %+v", err),
		})
		return
	}
	defer out.Close()

	// 把上传文件复制到指定位置
	_, err = io.Copy(out, file)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err5,
			"message": fmt.Sprintf("Copy: %+v", err),
		})
		return
	}

	// 把完整路径，更新到数据库
	i := int(id)
	param, err = param.FindByID(i)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err6,
			"message": fmt.Sprintf("FindByID: %+v", err),
		})
		return
	}

	// 设置完整路径，并保存
	param.FileLocation = fileLoc
	err = param.Update()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err7,
			"message": fmt.Sprintf("Update: %+v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rtnCode": ok,
	})
}

// Update 更新当前对象
func (d UploadController) Update(c *gin.Context) {
	// 创建一个对象，用以绑定 post 过来的 json
	param := models.CreateUploadItem()

	// post 过来的 json，绑定到 对象上；
	err := c.BindJSON(&param)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("Update BindJson error: %+v\n", err),
		})
		return
	}

	err = param.Update()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err2,
			"message": fmt.Sprintf("Update error: %+v\n", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rtnCode": ok,
	})
}

// DownloadFile 下载选中的附件，一个一个下载
func (d UploadController) DownloadFile(c *gin.Context) {
	item := models.CreateUploadItem()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("%#v", err),
		})
		return
	}

	item, err = item.FindByID(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err2,
			"message": fmt.Sprintf("FindByID: %#v", err),
		})
		return
	}

	fullpath := item.FileLocation
	f, err := os.Open(fullpath)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err3,
			"message": fmt.Sprintf("open: %#v", err),
		})
		return
	}
	defer f.Close()

	basename := item.FileName
	c.Writer.Header().Set("Content-Disposition", "attachment; filename="+basename)
	c.Writer.Header().Set("Content-Type", c.Request.Header.Get("Content-Type"))

	io.Copy(c.Writer, f)
}

// Delete 从数据库删除当前对象
func (d UploadController) Delete(c *gin.Context) {
	// 按照 id 查询
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("path param error: %+v\n", err),
		})
		return
	}

	item := models.CreateUploadItem()
	item.ID = id
	err = item.Delete()
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
