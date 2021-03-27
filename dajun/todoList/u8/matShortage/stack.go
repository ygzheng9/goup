package matShortage

import (
	"fmt"
	"sync"
)

// OneLevel BOM 的一个节点, stack 中的元素
type OneLevel struct {
	// 料号
	InvCode string
	// 相对于根节点的耗量
	BaseQty float64
}

// Stack 自己实现一个
type Stack struct {
	lock  sync.Mutex // you don't have to do this if you don't want thread safety
	items []OneLevel
}

// NewStack 创建对象
func NewStack() *Stack {
	return &Stack{sync.Mutex{}, make([]OneLevel, 0)}
}

// 从 stack 中弹出一个
func (s *Stack) pop() OneLevel {
	s.lock.Lock()
	defer s.lock.Unlock()

	l := len(s.items)
	if l == 0 {
		return OneLevel{}
	}

	res := s.items[l-1]
	s.items = s.items[:l-1]
	return res
}

// 压入一个到 stack
func (s *Stack) push(v OneLevel) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.items = append(s.items, v)
}

// 判断 stack 是否为 空
func (s *Stack) isEmpty() bool {
	return len(s.items) == 0
}

// 打印
func (s *Stack) dump() {
	for _, v := range s.items {
		fmt.Printf("invCode: %s, baseQty: %0.2f\n", v.InvCode, v.BaseQty)
	}
}
