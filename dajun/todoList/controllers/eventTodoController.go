package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"pickup/dajun/todoList/models"
)

// EventTodoController todo 相关的操作
type EventTodoController struct{}

// CreateEventTodoController 返回一个空对象
func CreateEventTodoController() EventTodoController {
	return EventTodoController{}
}

// FindAll 取得待办清单
func (d EventTodoController) FindAll(c *gin.Context) {
	item := models.CreateEventTodo()
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
func (d EventTodoController) FindByID(c *gin.Context) {
	// 从 url path 中提取 id
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("path param error: %+v\n", err),
		})
		return
	}

	item := models.CreateEventTodo()
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
func (d EventTodoController) FindByRef(c *gin.Context) {
	// 取得查询参数
	item := models.CreateEventTodo()
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

// FindByParam 根据 责任人 待办内容 查找列表
func (d EventTodoController) FindByParam(c *gin.Context) {

	type paramT struct {
		Cond string `json:"param" form:"param"`
	}

	param := paramT{}

	err := c.BindJSON(&param)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("FindByParam BindJson error: %+v\n", err),
		})
		return
	}

	item := models.CreateEventTodo()
	items, err := item.FindByParam(param.Cond)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err2,
			"message": fmt.Sprintf("FindByParam error: %+v\n", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rtnCode": ok,
		"items":   items,
	})
}

// Insert 当前对象插入到数据库
func (d EventTodoController) Insert(c *gin.Context) {
	// 创建一个对象，用以绑定 post 过来的 json
	param := models.CreateEventTodo()

	// post 过来的 json，绑定到 对象上；
	err := c.BindJSON(&param)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("Insert BindJson error: %+v\n", err),
		})
		return
	}

	err = param.Insert()
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

// Update 更新当前对象
func (d EventTodoController) Update(c *gin.Context) {
	// 创建一个对象，用以绑定 post 过来的 json
	param := models.CreateEventTodo()

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

// Delete 从数据库删除当前对象
func (d EventTodoController) Delete(c *gin.Context) {
	// 按照 id 查询
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("Delete param error: %+v\n", err),
		})
		return
	}

	item := models.CreateEventTodo()
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

// MarkTodo 更新 todo 状态
func (d EventTodoController) MarkTodo(c *gin.Context) {
	// post 过来的 json，绑定到 对象上；
	item := models.CreateEventTodo()
	err := c.BindJSON(&item)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("MarkTodo BindJson error: %+v\n", err),
		})
		return
	}

	err = item.UpdateStatus()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err2,
			"message": fmt.Sprintf("MarkTodo error: %+v\n", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rtnCode": ok,
	})
}

// UploadTodos 在一个 event 下，批量上载 todolist，
func (d EventTodoController) UploadTodos(c *gin.Context) {
	// 解析 request 中的 postData，里面记录了父节点
	type paramType struct {
		RefID    int    `json:"ref_id" form:"ref_id"`
		RefType  string `json:"ref_type" form:"ref_type"`
		RefTitle string `json:"ref_title" form:"ref_title"`
		UserName string `json:"user_name" form:"user_name"`
	}

	param := paramType{}
	err := c.Bind(&param)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("UploadTodos Bind: %#v", err),
		})
		return
	}

	// 获取上载的文件
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("FormFile error: %+v\n", err),
		})
		return
	}

	// 从上载的文件中，解析多条 todo，并保存到数据库中
	item := models.CreateEventTodo()
	items, err := item.LoadTodos(file)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err2,
			"message": fmt.Sprintf("LoadTodos error: %+v\n", err),
		})
		return
	}

	// 上载的 todos 保存到数据库
	for _, v := range items {
		v.RefID = param.RefID
		v.RefType = param.RefType
		v.RefTitle = param.RefTitle
		v.CreateUser = param.UserName

		err = v.Insert()
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"rtnCode": err3,
				"message": fmt.Sprintf("Insert error: %+v\n", err),
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"rtnCode": ok,
		"loadCnt": len(items),
	})
}

// SetMatRule 设置 物料数量规则
func SetMatRule(c *gin.Context) {
	// post 过来的参数
	type paramT struct {
		ID      int    `json:"id" form:"id"`
		MatRule string `json:"mat_rule" form:"mat_rule"`
		MatCode string `json:"mat_code" form:"mat_code"`
		MatCond string `json:"mat_cond" form:"mat_cond"`
	}

	p := paramT{}
	err := c.BindJSON(&p)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err1,
			"message": fmt.Sprintf("SetMatRule Bind: %#v", err),
		})
		return
	}
	fmt.Printf("\n%#v\n", p)

	err = models.SetMatRule(p.ID, p.MatRule, p.MatCode, p.MatCond)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err2,
			"message": fmt.Sprintf("SetMatRule: %#v", err),
		})
		return
	}

	// 规则设置好后，对本规则，自动扫描一次
	err = models.ScanMatRule(p.ID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": err3,
			"message": fmt.Sprintf("ScanMatRule: %#v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rtnCode": ok,
	})
}
