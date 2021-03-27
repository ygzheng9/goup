package models

import (
	"fmt"
	"testing"
)

func Test_FindByParam(t *testing.T) {
	item := CreateEventTodo()

	items, err := item.FindByParam("郑")
	if err != nil {
		t.Errorf("FindByParam: %+v\n", err)
		return
	}

	for _, v := range items {
		fmt.Printf("%#v\n", v)
	}
}

func Test_SetMatRule(t *testing.T) {
	err := SetMatRule(7, "0", "999-002000011-000", ">= 500")
	if err != nil {
		t.Errorf("Test_SetMatRule: %+v\n", err)
		return
	}
}

func Test_ClseEvent(t *testing.T) {
	err := CloseEvent(19, false)
	if err != nil {
		t.Errorf("Test_ClseEvent: %+v\n", err)
		return
	}
}
