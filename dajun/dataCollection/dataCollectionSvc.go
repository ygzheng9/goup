package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

const (
	connPort = ":8090"
	connType = "tcp"
	saveURL  = "http://localhost:8099/api/testingData_save"
)

func main() {
	// Listen for incoming connections.
	ln, err := net.Listen(connType, connPort)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer ln.Close()
	fmt.Println("Listening on " + connPort)

	for {
		// Listen for an incoming connection.
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		fmt.Printf("Received from: %s\n", conn.RemoteAddr())

		// Handle connections in a new goroutine.
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	// 读取 client 发来的信息，一次性读取
	result, err := readFully(conn)
	if err != nil {
		fmt.Printf("readFully error: %+v", err)
		return
	}

	resultUTF8, err := GbkToUtf8(result)

	// 解析数据
	fmt.Printf("will parse: %s\n", resultUTF8)
	err = postTestingData(resultUTF8)
	if err != nil {
		fmt.Printf("postTestingData error: %+v", err)
		conn.Write([]byte("post error"))
		return
	}

	// 返回给 client 的信息，一次性发送
	var msg bytes.Buffer
	msg.WriteString("ok")
	// msg.Write(result)

	conn.Write(msg.Bytes())
}

func readFully(conn net.Conn) ([]byte, error) {
	// 每次一读取的长度
	const size = 1024
	buf := make([]byte, size)

	// 循环读取的次数
	total := 0
	// 最终结果
	var result bytes.Buffer
	for {
		n, err := conn.Read(buf)
		total++
		result.Write(buf[:n])

		if err != nil {
			if err != io.EOF {
				log.Printf("Read error: %s", err)
				return nil, err
			}
			break
		}

		// 读取的长度，小于 buf 的长度，表示已读取完毕
		if n < size {
			// fmt.Printf("total: %d, %s\n", total, result.String())
			break
		}
	}

	return result.Bytes(), nil
}

// postTestingData 把数据 post 给 api
func postTestingData(b []byte) error {
	values := map[string]string{"data": string(b)}

	jsonValue, _ := json.Marshal(values)
	_, err := http.Post(saveURL, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Printf("post error: %+v\n", err)
		return err
	}

	return nil
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
