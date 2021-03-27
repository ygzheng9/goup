package models

import (
	"fmt"
	"testing"
)

func Test_reviewInsert(t *testing.T) {
	r := CreateReview()
	r.BizType = "PBC"
	r.TxName = "2018年度目标"
	r.FormContent = "akdsfkakkskdf"
	r.FormUser = "郑永刚"
	r.FormDate = "2017-11-20"

	_, err := ReviewInsert(r)
	if err != nil {
		t.Errorf("ReviewInsert error: %+v\n", err)
	}

	r.FormContent = "bababa"
	r.FormUser = "zhengyg"
	r.FormDate = "2017-11-21"
	_, err = ReviewInsert(r)
	if err != nil {
		t.Errorf("ReviewInsert error: %+v\n", err)
	}

	fmt.Printf("boo")
}

func Test_reviewFind(t *testing.T) {
	items, err := ReviewFindAll()
	if err != nil {
		t.Errorf("ReviewFindAll error: %+v\n", err)
	}

	for _, i := range items {
		fmt.Printf("%+v\n", i)
	}

	item, err := ReviewFindByID(2)
	if err != nil {
		t.Errorf("ReviewFindByID error: %+v\n", err)
	}

	item.ReviewUser = "郑永刚"
	item.ReviewContent = "继续努力"
	_, err = ReviewUpdate(item)
	if err != nil {
		t.Errorf("ReviewUpdate error: %+v\n", err)
	}

	err = ReviewDelete(item)
	if err != nil {
		t.Errorf("ReviewDelete error: %+v\n", err)
	}

	fmt.Printf("item: %+v\n", item)
}

func Test_reviewLogInsert(t *testing.T) {
	var r ReviewLog
	r.RefID = 1
	r.FromStatus = "草稿"
	r.ToStatus = "待审批"
	r.OpUser = "郑永刚"
	r.OpDate = "2017-11-21"

	err := ReviewLogInsert(r)
	if err != nil {
		t.Errorf("ReviewLogInsert error: %+v\n", err)
	}

	r.RefID = 1
	r.FromStatus = "待审批"
	r.ToStatus = "审批通过"
	r.OpUser = "郑永刚"
	r.OpDate = "2017-11-22"

	err = ReviewLogInsert(r)
	if err != nil {
		t.Errorf("ReviewLogInsert error: %+v\n", err)
	}
}

func Test_reviewLogFind(t *testing.T) {
	items, err := ReviewLogFindByRefID(1)
	if err != nil {
		t.Errorf("ReviewLogFindByRefID error: %+v\n", err)
	}

	for _, i := range items {
		fmt.Printf("%+v\n", i)
	}
}
