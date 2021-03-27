package models

import (
	"fmt"
	"testing"
	"time"
)

func Test_ActivityLogInsert(t *testing.T) {
	r := ActivityLog{}
	r.Msg1 = "wowowowo"
	r.OpUser = "郑永刚"
	r.OpDate = time.Now().Format("2006-01-02 15:04:05")

	err := ActivityLogInsert(r)
	if err != nil {
		t.Errorf("ActivityInsert error: %+v\n", err)
	}

	fmt.Printf("OK. ")
}

func Test_ActivityLogFindByRefID(t *testing.T) {
	items, err := ActivityLogFindByRefID(0, "")
	if err != nil {
		t.Errorf("ActivityInsert error: %+v\n", err)
	}

	for _, v := range items {
		fmt.Printf("item: %+v\n", v)
	}
}
