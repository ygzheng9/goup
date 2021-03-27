package util

import (
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

func Test_GenProcutReturnSQL(t *testing.T) {
	GenProcutReturnSQL("productReturn.tmpl", "productReturn.csv")
}

func Test_padLeft(t *testing.T) {
	fmt.Printf("%s\n", padLeft("2366", "0000000000"))
	fmt.Printf("%s\n", padLeft("10896", "0000000000"))

	str := "hello"

	convStr := func(input string) string {
		return fmt.Sprintf("'%s'", input)
	}
	out := convStr(str)
	fmt.Println(out)
}

func Test_foo(T *testing.T) {
	reader := strings.NewReader("Clear is better than clever")
	p := make([]byte, 4)
	for {
		n, err := reader.Read(p)
		// if err == io.EOF {
		// 	break
		// }

		if err != nil {
			if err == io.EOF {
				fmt.Println(string(p[:n])) //should handle any remainding bytes.
				break
			}
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println(string(p[:n]))
	}
}
