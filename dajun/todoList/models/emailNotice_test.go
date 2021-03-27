package models

import (
	"fmt"
	"testing"
)

func Test_email_Insert(t *testing.T) {
	item := CreateEmailNotice()

	item.UserName = "郑永刚"
	item.SendTo = "zhengyg@dajuntech.com"
	item.Subject = "流程通知"
	item.Content = "<b>第二封邮件</b>"

	err := item.Insert()
	if err != nil {
		fmt.Printf("Insert: %+v\n", err)
		return
	}
}

func Test_email_FindByID(t *testing.T) {
	item := CreateEmailNotice()
	item, err := item.FindByID(1)
	if err != nil {
		fmt.Printf("FindByID: %+v\n", err)
		return
	}

	fmt.Printf("%#v\n", item)
}

func Test_FindAll(t *testing.T) {
	item := CreateEmailNotice()
	items, err := item.FindAll()
	if err != nil {
		fmt.Printf("FindAll: %+v\n", err)
		return
	}

	for _, i := range items {
		fmt.Printf("%+v\n", i)
	}
}

func Test_email_Send(t *testing.T) {
	item := CreateEmailNotice()
	item, err := item.FindByID(107)
	if err != nil {
		fmt.Printf("FindByID: %+v\n", err)
		return
	}

	fmt.Printf("%#v\n", item)

	err = item.Send()
	if err != nil {
		t.Errorf("SentOut: %#v\n", err)
		return
	}
}

func Test_email_SendAll(t *testing.T) {
	err := SendAllEmails()
	if err != nil {
		t.Errorf("SentOut: %#v\n", err)
		return
	}
}
