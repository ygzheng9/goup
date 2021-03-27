package main

import (
	"fmt"
	"testing"
)

func Test_copyOutboundFile(t *testing.T) {
	err := copyOutboundFile()
	if err != nil {
		fmt.Printf("err: %+v\n", err)
	}
}
