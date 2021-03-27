package funcFP

import (
	"fmt"
	"log"
)

type Handler interface {
	Do(k string, v int) (string, int)
}

// 为这种类型的函数，定义一个类型，方便使用
type HandlerFuncGeneral func(k string, v int) (string, int)

// 由于需要为该类型的函数增加一些特性，所以采用一个 adapter
//type HandlerFunc func(k string, v int) (string, int)
type HandlerFunc HandlerFuncGeneral

func (f HandlerFunc) Do(k string, v int) (string, int) {
	return f(k, v)
}

func Each(m map[string]int, h Handler) {
	if m != nil && len(m) > 0 {
		for k, v := range m {
			s, i := h.Do(k, v)
			log.Printf("(%q, %d) ==> (%q, %d)", k, v, s, i)
		}
	}
}

// HandlerFunc(f) 为类型转换，结果就是 为 f 增加了 Do 的方法，从而变成了一个 Handler
//func EachFunc(m map[string]int, f func(k string, v int) (string, int)) {
//	Each(m, HandlerFunc(f))
//}
func EachFunc(m map[string]int, f HandlerFuncGeneral) {
	Each(m, HandlerFunc(f))
}

func Reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

// 满足通用类型  HandlerFuncGeneral
func SelfInfo(k string, v int) (result string, n int) {
	fmt.Printf("大家好，我是 %s, 今年 %d 岁\n", k, v)

	result = Reverse(k)
	n = len(k)
	return
}
