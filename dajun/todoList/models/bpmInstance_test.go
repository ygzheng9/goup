package models

import "testing"
import "fmt"

func Test_BpmInstanceStart(t *testing.T) {
	// 业务类型，业务ID
	soft, err := SoftInstFindByID(9)
	if err != nil {
		t.Errorf("Err:SoftInstFindByID %+v\n", err)
		return
	}
	r := soft.StartBPM()
	fmt.Printf("result: %t\n", r)
}

func Test_BpmInstanceNodeComplete(t *testing.T) {
	BpmInstanceNodeComplete(35, BpmApproved, "已处理完毕")
}

func Test_Printf(t *testing.T) {
	str := "asdfads"
	fmt.Printf("%q -- %s", str, str)
}
