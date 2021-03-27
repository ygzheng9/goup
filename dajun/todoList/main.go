package main

import (
	"fmt"
	"os"
	"path/filepath"

	// _ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	// cors "gopkg.in/gin-contrib/cors.v1"
	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"

	"pickup/dajun/todoList/config"
	"pickup/dajun/todoList/controllers"
	"pickup/dajun/todoList/models"
	"pickup/dajun/todoList/u8"
)

// 显示程序当前目录
func printPath() {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	fmt.Printf("Path: %s\n", exPath)
}

// 主程序
func main() {
	// 从文件加载配置信息，exe 的相同目录下；
	appConfig, err := config.LoadConfiguration("./todoList.json")
	if err != nil {
		fmt.Printf("LoadConfiguration failed. %+v\n", err)
		return
	}

	// 初始化数据库连接，邮件模板路径
	models.Setup(config.ConnectDB(), appConfig.TemplateDir)

	// U8 可用时，连接 U8
	if appConfig.U8.IsAvl {
		// 初始化 U8 连接
		u8.SetupDB(config.ConnectU8())
		fmt.Println("connect to  u8.")
	}

	if appConfig.Schedule {
		// 开启定时任务
		models.StartNoticeSvc()
		fmt.Println("start notice service.")
	}

	if appConfig.Release {
		// 发布模式
		gin.SetMode(gin.ReleaseMode)
	}

	// 初始化路由
	router := gin.Default()

	if appConfig.Release {
		// 只允许白名单
		config := cors.DefaultConfig()
		config.AllowOrigins = []string{"http://127.0.0.1:8000", "http://localhost:8000", "http://10.10.21.43:8000",
			"http://server22:8000", "http://10.10.10.222:8000"}
		router.Use(cors.New(config))

		fmt.Println("set to release.")
	} else {
		// 开发模式，允许所有
		router.Use(cors.Default())
		fmt.Printf("config: %+v\n", appConfig)
	}

	// 全部使用 api 访问，所以不需要使用 server-side template
	// router.Static("/assets", "./assets")
	// router.StaticFile("/favicon.ico", "./assets/favicon.ico")

	// router.LoadHTMLGlob("./templates/**/*")
	// router.GET("/ecn", func(c *gin.Context) {
	// 	c.HTML(http.StatusOK, "todolist/index.html", gin.H{
	// 		"title": "Posts",
	// 	})
	// })

	// 正式逻辑
	router.POST("/api/validateLogin", controllers.ValidateLogin)
	router.POST("/api/changePassword", controllers.ChangePassword)
	router.POST("/api/genPassword", controllers.GeneratePassword)

	userCtrl := controllers.CreateUserController()
	user := router.Group("/api/users")
	{
		user.GET("/", userCtrl.FindAll)
		user.GET("/:id", userCtrl.FindByID)
		user.DELETE("/:id", userCtrl.Delete)
		user.POST("/", userCtrl.Insert)
		user.PUT("/", userCtrl.Update)
	}
	router.POST("/api/users_search", userCtrl.Search)

	deptCtrl := controllers.CreateDepartmentController()
	department := router.Group("/api/departments")
	{
		department.GET("/", deptCtrl.FindAll)
		department.GET("/:id", deptCtrl.FindByID)
		department.DELETE("/:id", deptCtrl.Delete)
		department.POST("/", deptCtrl.Insert)
		department.PUT("/", deptCtrl.Update)
	}
	// 为部门批量增加用户
	router.POST("/api/departments_batchAddSubs", deptCtrl.BatchAddSubs)
	router.POST("/api/departments_removeSub", deptCtrl.RemoveSub)

	// 部门与用户关联
	router.GET("/api/departments_viewers", controllers.DeptViewerFindAll)
	router.POST("/api/departments_batchAddViewers", controllers.DeptViewerBatchAdd)
	router.POST("/api/departments_removeViewer", controllers.DeptViewerRemove)

	eventCtrl := controllers.CreateEventController()
	event := router.Group("/api/events")
	{
		event.GET("/", eventCtrl.FindAll)
		event.GET("/:id", eventCtrl.FindByID)
		event.DELETE("/:id", eventCtrl.Delete)
		event.POST("/", eventCtrl.Insert)
		event.PUT("/", eventCtrl.Update)
	}
	router.GET("/api/events_notify/:id", eventCtrl.NotifyOpen)
	router.GET("/api/events_close/:id", controllers.CloseEvent)
	router.POST("/api/events_search", controllers.FindEventByParam)

	todoCtrl := controllers.CreateEventTodoController()
	todo := router.Group("/api/todos")
	{
		todo.GET("/", todoCtrl.FindAll)
		todo.GET("/:id", todoCtrl.FindByID)
		todo.DELETE("/:id", todoCtrl.Delete)
		todo.POST("/", todoCtrl.Insert)
		todo.PUT("/", todoCtrl.Update)
	}
	router.POST("/api/todos_byref", todoCtrl.FindByRef)
	router.POST("/api/todos_mark", todoCtrl.MarkTodo)
	router.POST("/api/todos_upload", todoCtrl.UploadTodos)
	router.POST("/api/todos_search", todoCtrl.FindByParam)
	router.POST("/api/todos_matrule", controllers.SetMatRule)

	fdbkCtrl := controllers.CreateFeedbackController()
	feedback := router.Group("/api/feedback")
	{
		feedback.GET("/", fdbkCtrl.FindAll)
		feedback.GET("/:id", fdbkCtrl.FindByID)
		feedback.DELETE("/:id", fdbkCtrl.Delete)
		feedback.POST("/", fdbkCtrl.Insert)
		feedback.PUT("/", fdbkCtrl.Update)
	}
	router.POST("/api/feedback_ref", fdbkCtrl.FindByRef)

	uploadCtrl := controllers.CreateUploadController()
	upload := router.Group("/api/uploads")
	{
		upload.GET("/", uploadCtrl.FindAll)
		upload.GET("/:id", uploadCtrl.FindByID)
		upload.DELETE("/:id", uploadCtrl.Delete)
		upload.POST("/", uploadCtrl.Insert)
		//upload.POST("/", uploadHandler)
		upload.PUT("/", uploadCtrl.Update)
	}
	router.POST("/api/uploads_byref", uploadCtrl.FindByRef)
	router.GET("/api/uploads_download/:id", uploadCtrl.DownloadFile)

	filemgmtCtrl := controllers.CreateFileMgmtController()
	filemgmt := router.Group("/api/filemgmt")
	{
		filemgmt.DELETE("/:id", filemgmtCtrl.DeleteFile)
	}
	router.POST("/api/filemgmt_list", filemgmtCtrl.ListFiles)
	router.GET("/api/filemgmt_history/:id", filemgmtCtrl.FindHistory)
	router.POST("/api/filemgmt_rename", filemgmtCtrl.ChangeName)
	router.GET("/api/filemgmt_download/:id", filemgmtCtrl.DownloadFile)
	router.POST("/api/filemgmt_uploadNew", filemgmtCtrl.UploadNewFile)
	router.POST("/api/filemgmt_uploadMod", filemgmtCtrl.UploadModFile)
	router.POST("/api/filemgmt_createFolder", filemgmtCtrl.CreateFolder)
	router.GET("/api/filemgmt_allsubs/:id", filemgmtCtrl.GetAllSubs)

	assetCtrl := controllers.CreateAssetController()
	asset := router.Group("/api/assets")
	{
		asset.GET("/", assetCtrl.FindAll)
		asset.GET("/:id", assetCtrl.FindByID)
		asset.DELETE("/:id", assetCtrl.Delete)
		asset.POST("/", assetCtrl.Insert)
		asset.PUT("/", assetCtrl.Update)
	}

	softCtrl := controllers.CreateSoftInstController()
	soft := router.Group("/api/soft")
	{
		soft.GET("/", softCtrl.FindAll)
		soft.GET("/:id", softCtrl.FindByID)
		soft.DELETE("/:id", softCtrl.Delete)
		soft.POST("/", softCtrl.Insert)
		soft.PUT("/", softCtrl.Update)
	}
	router.POST("/api/soft_search", softCtrl.FindBy)

	// U8 检查库存数量
	router.POST("/api/u8/checkStock", u8.CheckStock)
	// U8 销售发货单，批改 批次号
	router.POST("/api/u8/outboundUpload", u8.OutboundUpload)
	router.POST("/api/u8/outboundUploadSvc", u8.OutboundUploadSvc)

	// U8 批量修改保质期
	router.POST("/api/u8/massDateUpload", u8.MassDateUpload)

	// U8 生产订单子件需求 - 生产备料仓现有库存
	router.GET("/api/u8/poinvdiff", u8.FetchPoInvDiff)

	// U8 钢印号追溯
	router.POST("/api/u8/traceSN", u8.TraceSN)

	// 邮件通知
	router.POST("/api/listMails", controllers.ListEmails)
	// 物料提醒
	router.GET("/api/listMatNotice", controllers.ListMatNotice)

	// 显示磁盘空间监控的列表
	diskspace := router.Group("/api/diskspace")
	{
		diskspace.GET("/", controllers.DiskSpaceGetAll)
		diskspace.GET("/:id", controllers.DiskSpaceFindByID)
		diskspace.DELETE("/:id", controllers.DiskSpaceDelete)
		diskspace.POST("/", controllers.DiskSpaceInsert)
		diskspace.PUT("/", controllers.DiskSpaceUpdate)
	}
	// 修改邮件通知激活状态
	router.POST("/api/diskspace_status", controllers.DiskSpaceToggleStatus)
	router.GET("/api/diskspace_scan/:id", controllers.DiskSpaceScan)

	// Review 单个节点的审批
	review := router.Group("/api/review")
	{
		review.GET("/", controllers.ReviewGetAll)
		review.GET("/:id", controllers.ReviewFindByID)
		review.DELETE("/:id", controllers.ReviewDelete)
		review.POST("/", controllers.ReviewInsert)
		review.PUT("/", controllers.ReviewUpdate)
	}
	router.POST("/api/review_doAction", controllers.ReviewDoAction)
	router.GET("/api/review_loadlogs/:id", controllers.ReviewLoadLogs)

	// 每天工作记录
	dailyTask := router.Group("/api/dailytask")
	{
		dailyTask.GET("/", controllers.DailyTaskGetAll)
		dailyTask.GET("/:id", controllers.DailyTaskFindByID)
		dailyTask.POST("/", controllers.DailyTaskInsert)
		dailyTask.PUT("/", controllers.DailyTaskUpdate)
	}
	router.POST("/api/dailytask_searchby", controllers.DailyTaskSearchBy)

	// 操作日志
	router.POST("/api/activitylog_searchby", controllers.ActivityLogFindByRefID)

	// 监控的 BOM 清单
	matWatch := router.Group("/api/matWatch")
	{
		matWatch.GET("/", controllers.MatWLGetAll)
		matWatch.GET("/:id", controllers.MatWLFindByID)
		matWatch.POST("/", controllers.MatWLInsert)
		matWatch.DELETE("/:id", controllers.MatWLDelete)
	}
	router.POST("/api/matWatch_batch", controllers.MatWLBatchInsert)

	// 公告信息
	article := router.Group("/api/article")
	{
		article.GET("/", controllers.ArticleGetAll)
		article.GET("/:id", controllers.ArticleFindByID)
		article.POST("/", controllers.ArticleInsert)
		article.PUT("/", controllers.ArticleUpdate)
		article.DELETE("/:id", controllers.ArticleDeleteByID)
	}

	// 合同管理
	contract := router.Group("/api/contract")
	{
		contract.GET("/", controllers.ContractGetAll)
		contract.GET("/:id", controllers.ContractFindByID)
		contract.POST("/", controllers.ContractInsert)
		contract.PUT("/", controllers.ContractUpdate)
		contract.DELETE("/:id", controllers.ContractDeleteByID)
	}

	// 付款计划
	milestone := router.Group("/api/milestone")
	{
		milestone.GET("/", controllers.MilestoneGetAll)
		milestone.GET("/:id", controllers.MilestoneFindByID)
		milestone.POST("/", controllers.MilestoneInsert)
		milestone.PUT("/", controllers.MilestoneUpdate)
		milestone.DELETE("/:id", controllers.MilestoneDeleteByID)
	}
	router.GET("/api/milestone_byContract/:id", controllers.MilestoneFindByContract)

	// 发票
	invoice := router.Group("/api/invoice")
	{
		invoice.GET("/", controllers.InvoiceGetAll)
		invoice.GET("/:id", controllers.InvoiceFindByID)
		invoice.POST("/", controllers.InvoiceInsert)
		invoice.PUT("/", controllers.InvoiceUpdate)
		invoice.DELETE("/:id", controllers.InvoiceDeleteByID)
	}
	router.GET("/api/invoice_byContract/:id", controllers.InvoiceFindByContract)
	router.POST("/api/invoice_handover", controllers.InvoiceHandOver)
	router.POST("/api/invoice_paymentrequest", controllers.InvoicePaymentRequest)

	// 上载测试数据
	router.POST("/api/testingData_upload", controllers.TestingDataUpload)
	// 网络服务，保存测试数据
	router.POST("/api/testingData_save", controllers.TestingDataSave)
	// 查询测试记录头信息
	router.POST("/api/testingData_search", controllers.TestingDataSearch)
	// 删除测试记录，头信息+行信息
	router.DELETE("/api/testingData/:id", controllers.TestingDataDelete)

	// 根据测试头，取得测试项明细
	router.GET("/api/testingData/:headid/details", controllers.TestingDataLoadDeatail)

	// 审批流程配置
	bpmProcess := router.Group("/api/bpmProcess")
	{
		bpmProcess.GET("/", controllers.BpmProcessGetAll)
		bpmProcess.GET("/:id", controllers.BpmProcessFindByID)
		bpmProcess.POST("/", controllers.BpmProcessInsert)
		bpmProcess.PUT("/", controllers.BpmProcessUpdate)
		bpmProcess.DELETE("/:id", controllers.BpmProcessDeleteByID)
	}
	router.POST("/api/bpmProcess_checkrule", controllers.BpmProcessCheckRule)

	bpmProcessNode := router.Group("/api/bpmNode")
	{
		bpmProcessNode.GET("/", controllers.BpmProcessNodeGetAll)
		bpmProcessNode.GET("/:id", controllers.BpmProcessNodeFindByID)
		bpmProcessNode.POST("/", controllers.BpmProcessNodeInsert)
		bpmProcessNode.PUT("/", controllers.BpmProcessNodeUpdate)
		bpmProcessNode.DELETE("/:id", controllers.BpmProcessNodeDeleteByID)
	}
	router.POST("/api/bpmNode_checkrule", controllers.BpmProcessNodeCheckRule)

	router.POST("/api/laborclaim_search", controllers.LaborClaimFindByWk)
	router.POST("/api/laborclaim_batchSave", controllers.LaborClaimBatchSave)
	router.POST("/api/laborclaim_searchByPeriod", controllers.LaborClaimSearchByPeriod)

	projects := router.Group("/api/projects")
	{
		projects.GET("/", controllers.ProjectFindAll)
		projects.POST("/", controllers.ProjectInsert)
		projects.PUT("/", controllers.ProjectUpdate)
		projects.DELETE("/:id", controllers.ProjectDelete)
	}

	router.Run(appConfig.Port)
}
