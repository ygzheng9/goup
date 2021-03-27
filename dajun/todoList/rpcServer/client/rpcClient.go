package main

import (
	"fmt"
	"log"
	"net/rpc"

	"pickup/dajun/todoList/server"
)

func main() {
	// serverAddress := "localhost"
	serverAddress := "10.10.21.50"
	// serverAddress := "10.10.10.222"

	client, err := rpc.DialHTTP("tcp", serverAddress+":8999")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	fmt.Println("Call Real.GetDiskSpace")

	// Synchronous call
	space := &server.DiskSpace{}
	err = client.Call("Real.GetDiskSpace", "C", &space)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	fmt.Printf("GetDiskSpace: %+v\n", space)

}
