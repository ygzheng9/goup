package funcFP

import "testing"

func TestABA(t *testing.T) {
	cat := VoicerFunc(CatMiao)
	cat.Bark("通过Bark调用：成功变形")
	cat("直接调用：成功了")

	ShowUp(VoicerFunc(CatMiao), "balabala")

	ShowUp(VoicerFunc(DogWong), "wowowowo")

	a := VoicerNum(71)
	a.Bark("hahaha")

	ShowUp(a, "hahaha")

	// 类型不匹配，编译时会检查出来
	// b := VoicerNum("12")
	//log.Println(b)
	//ShowUp(b, "hahaha")
}
