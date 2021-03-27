package models

import (
	"fmt"
	"testing"
)

func Test_Fedback_FindByRef(t *testing.T) {
	item := Feedback{}

	item.RefID = 5
	item.RefType = "SoftIns"

	items, err := item.FindByRef()
	if err != nil {
		fmt.Printf("FindByRef error: %+v\n", err)
		return
	}

	for _, v := range items {
		fmt.Printf("%+v\n", v)
	}
}
