package main

import (
	"bufio"
	"fmt"
	"os"

	"pickup/dajun/todoList/u8"
)

func main() {
	reply, err := u8.ReadOutboundFile("./Book1.xlsx")
	if err != nil {
		fmt.Printf("err: %+v\n", err)
		return
	}

	fmt.Printf("reply: %+v\n", reply)

	fmt.Print("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	return
}
