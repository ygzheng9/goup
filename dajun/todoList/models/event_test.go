package models

import (
	"fmt"
	"testing"
)

func Test_updateOpenCnt(t *testing.T) {
	a := CreateEvent()
	items, err := a.FindAll()
	if err != nil {
		if err != nil {
			t.Errorf("FindAll: %+v\n", err)
			return
		}
	}

	// items, err = updateOpenCnt(items)
	// if err != nil {
	// 	if err != nil {
	// 		t.Errorf("updateOpenCnt: %+v\n", err)
	// 		return
	// 	}
	// }

	for _, i := range items {
		fmt.Printf("%+v\n", i)
	}
}

func Test_FindEventByParam(t *testing.T) {
	items, err := FindEventByParam("2017-10-10", "2017-10-30", "ALL", "王仁军")
	if err != nil {
		if err != nil {
			t.Errorf("FindAll: %+v\n", err)
			return
		}
	}

	for _, i := range items {
		fmt.Printf("%+v\n", i)
	}
}
