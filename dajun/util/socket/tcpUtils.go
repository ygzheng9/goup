package main

import (
	"bufio"
	"encoding/gob"
	"io"
	"log"
	"net"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

type complexData struct {
	N int
	S string
	M map[string]int
	P []byte
	C *complexData
}

const (
	// Port 服务器端口
	Port = ":61000"
)

// Open 连接到服务器端口
func Open(addr string) (*bufio.ReadWriter, error) {
	log.Println("Dial " + addr)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, errors.Wrap(err, "Dialing "+addr+" failed")
	}
	return bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn)), nil
}

// HandleFunc 服务器端处理函数
type HandleFunc func(*bufio.ReadWriter)

// Endpoint 服务器的抽象：根据不同类型，对应不同的处理函数
type Endpoint struct {
	listener net.Listener
	handler  map[string]HandleFunc
	m        sync.RWMutex
}

// NewEndpoint 新建一个服务器
func NewEndpoint() *Endpoint {
	return &Endpoint{
		handler: map[string]HandleFunc{},
	}
}

// AddHandleFunc 为服务器绑定一个处理程序
func (e *Endpoint) AddHandleFunc(name string, f HandleFunc) {
	e.m.Lock()
	e.handler[name] = f
	e.m.Unlock()
}

// Listen 服务器端的监听
func (e *Endpoint) Listen() error {
	var err error
	e.listener, err = net.Listen("tcp", Port)
	if err != nil {
		return errors.Wrap(err, "Unable to listen on "+e.listener.Addr().String()+"\n")
	}
	log.Println("Listen on", e.listener.Addr().String())

	for {
		log.Println("Accept a connection request.")
		conn, err := e.listener.Accept()
		if err != nil {
			log.Println("Failed accepting a connection request:", err)
			continue
		}
		log.Println("Handle incoming messages.")
		go e.handleMessages(conn)
	}
}

// handleMessages 先读取第一行的 command，然后根据 command 找到对应的处理函数，并执行
func (e *Endpoint) handleMessages(conn net.Conn) {
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	defer conn.Close()

	for {
		log.Print("Receive command '")
		cmd, err := rw.ReadString('\n')
		switch {
		case err == io.EOF:
			log.Println("Reached EOF - close this connection.\n   ---")
			return
		case err != nil:
			log.Println("\nError reading command. Got: '"+cmd+"'\n", err)
			return
		}

		// 取出客户端传来的命令
		cmd = strings.Trim(cmd, "\n")
		log.Println(cmd + "'")

		// 根据命令找到对应的处理程序
		e.m.RLock()
		handleCommand, ok := e.handler[cmd]
		e.m.RUnlock()
		if !ok {
			log.Println("Command '" + cmd + "' is not registered.")
			return
		}

		// 执行命令对应的程序函数
		handleCommand(rw)
	}
}

// handleStrings 处理 STRING：读入一行，写入一行
func handleStrings(rw *bufio.ReadWriter) {
	log.Print("Receive STRING message:")

	// 读取一行：客户端发来的信息只有一行
	s, err := rw.ReadString('\n')
	if err != nil {
		log.Println("Cannot read from connection.\n", err)
	}
	s = strings.Trim(s, "\n")
	log.Println(s)

	// 写入一行
	_, err = rw.WriteString("Thank you.\n")
	if err != nil {
		log.Println("Cannot write to connection.\n", err)
	}
	err = rw.Flush()
	if err != nil {
		log.Println("Flush failed.", err)
	}
}

// handleGob 处理 GOB
func handleGob(rw *bufio.ReadWriter) {
	log.Print("Receive GOB data:")
	var data complexData

	// 从 Reader 创建一个 Decorder，再 Decode 到 struct
	dec := gob.NewDecoder(rw)
	err := dec.Decode(&data)
	if err != nil {
		log.Println("Error decoding GOB data:", err)
		return
	}

	log.Printf("Outer complexData struct: \n%#v\n", data)
	log.Printf("Inner complexData struct: \n%#v\n", data.C)
}
