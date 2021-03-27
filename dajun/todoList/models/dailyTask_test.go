package models

import "testing"
import "fmt"

func Test_DailyTaskInsert(t *testing.T) {
	r := DailyTask{}

	r.UserName = "郑永刚"
	r.BizDate = "2017-11-22"
	r.WorkRemark = "bababal"

	_, err := DailyTaskInsert(r)
	if err != nil {
		t.Errorf("DailyTaskInsert error: %+v\n", err)
	}

	r.UserName = "郑永刚"
	r.BizDate = "2017-11-21"
	r.WorkRemark = "gogogooo"

	item, err := DailyTaskInsert(r)
	if err != nil {
		t.Errorf("DailyTaskInsert error: %+v\n", err)
	}
	fmt.Printf("item: %+v\n", item)

	items, err := DailyTaskFindAll()
	if err != nil {
		t.Errorf("DailyTaskFindAll error: %+v\n", err)
	}

	for _, v := range items {
		fmt.Printf("%+v\n", v)
	}
}
