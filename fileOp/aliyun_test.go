package main

import (
	"fmt"
	"testing"
)

func Test_connectAliB(t *testing.T) {
	connectAli()
}

func Test_setCors(t *testing.T) {
	setCors()
}

// "300Wx300H/YVA03-01-0004"

func Test_getFileList(t *testing.T) {
	files, err := getFileList("YVA03-01-0004")

	if err != nil {
		t.Errorf("err: %s", err)
	}

	for _, v := range files {
		fmt.Println(v)
	}
}
