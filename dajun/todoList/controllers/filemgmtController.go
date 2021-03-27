package controllers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"pickup/dajun/todoList/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// FileMgmtController 文件管理的操作逻辑
type FileMgmtController struct{}

// CreateFileMgmtController 返回一个空对象
func CreateFileMgmtController() FileMgmtController {
	return FileMgmtController{}
}

// ListFiles 先从数据库中，根据 parent_id 检索，如果没有记录，再到fs中查找
func (t FileMgmtController) ListFiles(c *gin.Context) {
	item := models.CreateFileMgmt()
	err := c.BindJSON(&item)

	// 先从数据库中找子元素
	items, err := item.FindSubs()

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("%#v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rtnCode": ok,
		"files":   items,
	})
}

// DeleteFile 删除当前选中的文件
func (t FileMgmtController) DeleteFile(c *gin.Context) {
	item := models.CreateFileMgmt()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("Atoi error: %#v", err),
		})
		return
	}

	item.LoadByID(id)

	err = item.RemoveFile()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err2,
			"message": fmt.Sprintf("RemoveFile error: %#v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rtnCode": ok,
	})
}

// FindHistory 查找全部修改历史
func (t FileMgmtController) FindHistory(c *gin.Context) {
	item := models.CreateFileMgmt()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("Atoi error: %#v", err),
		})
		return
	}

	item.LoadByID(id)

	items, err := item.FindHistory()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err2,
			"message": fmt.Sprintf("FindHistory error: %#v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rtnCode": ok,
		"files":   items,
	})
}

// ChangeName 重命名当前选中的文件，新名字是传入的 file_name；
func (t FileMgmtController) ChangeName(c *gin.Context) {
	item := models.CreateFileMgmt()
	err := c.BindJSON(&item)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("BindJson: %#v", err),
		})
		return
	}

	err = item.RenameFile(item.FileName)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err2,
			"message": fmt.Sprintf("RenameFile: %#v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rtnCode": ok,
	})
}

// DownloadFile 下载文件
func (t FileMgmtController) DownloadFile(c *gin.Context) {
	item := models.CreateFileMgmt()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("%#v", err),
		})
		return
	}

	item.LoadByID(id)

	fullpath, err := item.GetFullName()
	basename := path.Base(fullpath)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("GetFullName: %#v", err),
		})
		return
	}

	f, err := os.Open(fullpath)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err2,
			"message": fmt.Sprintf("open: %#v", err),
		})
		return
	}
	defer f.Close()

	c.Writer.Header().Set("Content-Disposition", "attachment; filename="+basename)
	c.Writer.Header().Set("Content-Type", c.Request.Header.Get("Content-Type"))

	io.Copy(c.Writer, f)
}

// UploadNewFile 上载一个新文件，名字和当前目录下文件都不相同
func (t FileMgmtController) UploadNewFile(c *gin.Context) {
	// 这里是上载文件，是 form 的 post，不是 json；
	// post 过来的是 当前目录
	currPath := models.CreateFileMgmt()
	err := c.Bind(&currPath)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("bind: %+v", err),
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

	// 把上载的文件，保存到当前目录中
	// 创建一个新文件
	filename := header.Filename
	fmt.Println(header.Filename)
	fullName, err := currPath.GetFullName()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("GetFullName: %+v", err),
		})
		return
	}
	fullName = fullName + "/" + filename
	out, err := os.Create(fullName)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("Create: %+v", err),
		})
		return
	}
	defer out.Close()

	// 把上传文件复制到指定位置
	_, err = io.Copy(out, file)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("Copy: %+v", err),
		})
		return
	}

	// 在数据库中，插入一条新纪录
	// 参数是当前目录，上载的文件放到当前目录下
	item := models.CreateFileMgmt()
	item.ParentID = currPath.ID
	item.FileName = header.Filename
	item.IsDir = "false"
	item.Depth = currPath.Depth + 1
	item.Status = "A"
	item.UpdateDate = time.Now().Format("2006-01-02 15:04:05")

	err = item.Insert()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("Insert %+v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "上传成功",
	})
}

// UploadModFile 对已有文件，上载一个新版本
func (t FileMgmtController) UploadModFile(c *gin.Context) {
	// 这里是上载文件，是 form 的 post，不是 json；
	// post 过来的是 已有的文件
	oldItem := models.CreateFileMgmt()
	err := c.Bind(&oldItem)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("Bind %+v", err),
		})
		return
	}
	log.Printf("ref: %+v\n", oldItem)

	// 旧文件备份
	err = oldItem.BackupFile()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("BackupFile： %+v", err),
		})
		return
	}

	// 获取上载的文件
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("FormFile: %+v", err),
		})
		return
	}

	// oldItem 是已经存在的同名文件
	fullName, err := oldItem.GetFullName()

	// oldItem 已经做过 rename，所以再创建一个新文件
	out, err := os.Create(fullName)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("Create %+v", err),
		})
		return
	}
	defer out.Close()

	// 把上传文件复制到指定位置
	_, err = io.Copy(out, file)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("Copy %+v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rtnCode": ok,
	})
}

// CreateFolder 创建一个新目录
func (t FileMgmtController) CreateFolder(c *gin.Context) {
	// post 过来的是 id 是当前目录，fileName 是新folder的名字
	currPath := models.CreateFileMgmt()
	err := c.BindJSON(&currPath)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("BindJSON: %+v", err),
		})
		return
	}
	log.Printf("ref: %+v\n", currPath)

	// 新目录的完整路径
	newFolderName := currPath.FileName
	currPath.LoadByID(currPath.ID)
	fullName, err := currPath.GetFullName()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("GetFullName: %+v", err),
		})
		return
	}
	newFullName := fullName + "/" + newFolderName

	err = os.MkdirAll(newFullName, 0777)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("MkdirAll: %+v", err),
		})
		return
	}

	// 参数是当前目录，上载的文件放到当前目录下
	item := models.CreateFileMgmt()
	item.ParentID = currPath.ID
	item.FileName = newFolderName
	item.IsDir = "true"
	item.Depth = currPath.Depth + 1
	item.Status = "A"
	item.UpdateDate = time.Now().Format("2006-01-02 15:04:05")

	err = item.Insert()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("Insert: %+v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rtnCode": ok,
		"message": "目录创建成功",
	})
}

//GetAllSubs 从 db 中，获取当前目录下所有文件，包括打删除标记的文件
func (t FileMgmtController) GetAllSubs(c *gin.Context) {
	// id 是当前目录
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("GetAllSubs: %#v", err),
		})
		return
	}

	currPath := models.CreateFileMgmt()
	currPath.LoadByID(id)

	items, err := currPath.NextLevelRaw()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf(" NextLevelRaw: %#v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rtnCode": ok,
		"files":   items,
	})
}
