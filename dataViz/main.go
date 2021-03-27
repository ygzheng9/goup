package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"pickup/dataViz/rawInput"
)

// DATAFOLDER 数据文件存放的路径
// const DATAFOLDER = "C:/localWork/dataViz/data/"
// const DATAFOLDER = "E:/99.localDev/tmp/"
// const DATAFOLDER = "/Users/ygzheng/Documents/playground/rawData/"
const DATAFOLDER = "./data/"

// getPOHead 读取 PO 头信息
func getPOHead(c *gin.Context) {
	items, err := rawInput.ReadPOHead(DATAFOLDER + "1.dt")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": -1,
			"message": fmt.Sprintf("ReadPOHead error. %s", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rtnCode": 0,
		"items":   items,
	})
}

// loadPOItems 加载采购行项目
func loadPOItems(c *gin.Context) {
	// 取得传入参数
	param := struct {
		Start string `json:"start"`
		End   string `json:"end"`
	}{}

	err := c.BindJSON(&param)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": -1,
			"message": fmt.Sprintf("BindJSON: %+v\n", err),
		})
		return
	}

	// 行项目
	items, err := rawInput.LoadPOItems(DATAFOLDER + "3.dt")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": -2,
			"message": fmt.Sprintf("LoadPOItems: %s", err),
		})

		return
	}

	results := rawInput.POItemsByDate(items, param.Start, param.End)
	c.JSON(http.StatusOK, gin.H{
		"rtnCode": 0,
		"items":   results,
	})
}

// loadMatByMonth 读取按月汇总的物料信息
func loadMatByMonth(c *gin.Context) {
	items, err := rawInput.LoadMatByMonth(DATAFOLDER + "4.dt")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": -1,
			"message": fmt.Sprintf("LoadMatByMonth: %s", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rtnCode": 0,
		"items":   items,
	})
}

// loadBomComponent 读取单层 BOM
func loadBomComponent(c *gin.Context) {
	items, err := rawInput.LoadBOMComponent(DATAFOLDER + "5.dt")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": -1,
			"message": fmt.Sprintf("LoadBOMComponent: %s", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rtnCode": 0,
		"items":   items,
	})
}

// loadMatInfo 读取物料基本信息
func loadMatInfo(c *gin.Context) {
	items, err := rawInput.LoadMatInfo(DATAFOLDER + "6.dt")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, gin.H{
			"rtnCode": -1,
			"message": fmt.Sprintf("LoadMatInfo: %s", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rtnCode": 0,
		"items":   items,
	})
}

func main() {
	router := gin.Default()

	// 允许跨域访问
	router.Use(cors.Default())

	// router.StaticFS("/", http.Dir("./build"))

	// router.GET("/api/hello", func(c *gin.Context) {
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"rtnCode": 0,
	// 		"message": "server is running....",
	// 	})
	// })

	router.GET("/api/dataviz/poHeads", getPOHead)
	router.POST("/api/dataviz/poItems", loadPOItems)
	router.GET("/api/dataviz/matByMonth", loadMatByMonth)
	router.GET("/api/dataviz/bomComponent", loadBomComponent)
	router.GET("/api/dataviz/loadMatInfo", loadMatInfo)

	// router.StaticFS("/", http.Dir("./build"))
	// router.Static("/", "./build")

	router.Run(":8077")
}
