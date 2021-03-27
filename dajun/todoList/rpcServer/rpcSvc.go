package main

import (
	"fmt"
	"net/http"
	"net/rpc"
	"time"

	"pickup/dajun/todoList/server"
)

func main() {
	// 注册 rpc 对象
	real := new(server.Real)
	rpc.Register(real)

	// rpc over HTTP
	rpc.HandleHTTP()

	now := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("service running from %s\n", now)

	// 监听端口
	port := ":8999"
	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}
