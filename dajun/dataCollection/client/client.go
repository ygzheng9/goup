package main

import (
	"fmt"
	"net"
)

func main() {
	// connect to this socket
	conn, err := net.Dial("tcp", "localhost:8090")
	if err != nil {
		fmt.Printf("can not connect to server. %+v\n", err)
		return
	}
	fmt.Println("connected to server.")
	defer conn.Close()

	// 发送信息给 server，一次性发送
	n, err := conn.Write([]byte("hello, I'm from client, can you see me?"))
	if err != nil {
		fmt.Printf("err: %+v\n", err)
	}
	fmt.Printf("send %d bytes.\n", n)

	// 获取反馈 server，一次性读取
	buf := make([]byte, 1024)
	reqLen, err := conn.Read(buf)
	if err != nil {
		fmt.Printf("Read error: %+v\n ", err)
		return
	}
	fmt.Printf("from server: %s\n", string(buf[:reqLen]))

}

func handleWrite(conn net.Conn, done chan string) {
	n, err := conn.Write([]byte("hello, I'm from client, can you see me?"))
	if err != nil {
		fmt.Printf("err: %+v\n", err)
	}
	fmt.Printf("send %d bytes.  Wait for replay.\n", n)

	done <- "Sent"
}

func handleRead(conn net.Conn, done chan string) {
	buf := make([]byte, 1024)
	reqLen, err := conn.Read(buf)
	if err != nil {
		fmt.Printf("Read error: %+v\n ", err)
		return
	}
	fmt.Printf("from server: %s\n", string(buf[:reqLen-1]))
	done <- "Read"
}
