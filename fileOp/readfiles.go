package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// func getAllFileNames(baseDir: string) {
// 	dir_list, e := ioutil.ReadDir(baseDir)
// 	if e != nil {
//         fmt.Println("read dir error")
//         return []string
//     }

// 		for i, v := range dir_list {
//         fmt.Println(i, "=", v.Name())
// 		}

// }

func main2() {
	fileName := "E:/99.localDev/qtProj/DJPlan5.6/availabilitycheck.cpp"
	// fileName := "E:/99.localDev/qtProj/DJPlan5.6/mainwindow.cpp"

	f, err := os.Open(fileName)

	if err != nil {
		panic(err)
	}
	defer f.Close()

	rd := bufio.NewReader(f)
	for {
		line, err := rd.ReadString('\n') //以'\n'为结束符读入一行

		if err != nil || io.EOF == err {
			break
		}

		fmt.Println(line)
	}
}
