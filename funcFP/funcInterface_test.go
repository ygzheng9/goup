package funcFP

import "testing"

func TestSelfInfo(t *testing.T) {
	persons := make(map[string]int)
	persons["张三"] = 20
	persons["李四"] = 23
	persons["王五"] = 26

	EachFunc(persons, SelfInfo)
}
