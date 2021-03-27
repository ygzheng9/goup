package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/gommon/log"
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
	"time"
)

// 全局连接
var db *sqlx.DB

func setupDB() {
	// 连接数据库
	var err error
	db, err = sqlx.Open("mysql", "djuser:P@ss1234@tcp(10.10.10.222:3306)/world?parseTime=true")
	if err != nil {
		panic(err)
	}

	// Test the connection to the database
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}
}

type Tick struct {
	ID          int    `db:"ID" json:"id"`
	TickDt      string `db:"TICK_DT" json:"tick_dt"`
	MAC         string `db:"MAC" json:"mac"`
	CPU         string `db:"CPU_LOAD" json:"cpu_load"`
	CPUNumber   string `db:"CPU_NUMBER" json:"cpu_number"`
	CPUType     string `db:"CPU_TYPE" json:"cpu_type"`
	CPULevel    string `db:"CPU_LEVEL" json:"cpu_level"`
	CPURevision string `db:"CPU_REVISION" json:"cpu_revision"`
	MemoryTotal string `db:"MEMORY_TOTAL" json:"memory_total"`
	MemoryUsed  string `db:"MEMORY_USED" json:"memory_used"`
	AppTitle    string `db:"APP_TITLE" json:"app_title"`
}

func saveTick(tick Tick) (err error) {
	cmd := `insert T_TICK (TICK_DT, MAC, CPU_LOAD, CPU_NUMBER, CPU_TYPE, CPU_LEVEL, CPU_REVISION, MEMORY_TOTAL, MEMORY_USED, APP_TITLE)
			values (:TICK_DT, :MAC, :CPU_LOAD, :CPU_NUMBER, :CPU_TYPE, :CPU_LEVEL, :CPU_REVISION, :MEMORY_TOTAL, :MEMORY_USED, :APP_TITLE)`

	tick.TickDt = time.Now().Format("2006-01-02 15:04:05")
	_, err = db.NamedExec(cmd, &tick)

	return err
}

func processTick(c *gin.Context) {
	tick := Tick{}
	if c.BindJSON(&tick) == nil {
		//log.Printf("%+v", tick)

		err := saveTick(tick)

		if err != nil {
			log.Print(err)
			c.JSON(http.StatusOK, gin.H{
				"status":  "error",
				"message": "can not save",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		}

	} else {
		log.Print("error in processTick")
		c.JSON(http.StatusOK, gin.H{
			"status":  "error",
			"message": "can not parse tick.",
		})
	}
}

func dummy(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
				"status":  "ok",
				"message": "dummy.",
			})
}

func main() {
	setupDB()

	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	router.Static("/", "./")

	router.POST("/api/tick", processTick)
	// router.POST("/api/tick", dummy)
	

	router.Run(":8088")
}
