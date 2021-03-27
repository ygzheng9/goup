package funcFP

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSum(t *testing.T) {
	// ex1 是一个函数
	ex1 := Sum(Const(1), Const(2), Const(3))
	t.Logf("result: %+v", ex1)

	// result 是一个 float64
	result := ex1()
	assert.Equal(t, result, 6.0, "they should be equal")

}

func TestCollect(t *testing.T) {
	fmt.Println(Collect(Mk(1, 2, 3))) // prints [1 2 3]

	expects := []float64{1, 2, 3}

	assert.Equal(t, expects, Collect(Mk(1, 2, 3)), "collect failed")
}

func TestTake(t *testing.T) {
	fmt.Println(Collect(Take(10, Naturals())))

	a := Mk(1, 2, 3)
	b := Take(2, a)
	c := Collect(b)
	fmt.Println(c)

	d, ok := b.Next()
	if !ok {
		fmt.Printf("no more")
	} else {
		fmt.Printf("next is %f\n", d)
	}

}

func TestTake2(t *testing.T) {
	e := Mk(4, 5, 6, 7)

	// 这里不是一次性取出前两个，而是返回了一个 Seq，
	// 说白了，就是根本没有做任何事，一直等到 Next() 调用才取数
	f := Take(2, e)

	// 执行一次 Next()，就取出一个数；然后再执行一次 Next().再取一个数
	idx := 0
	for v, ok := f.Next(); ok; {
		fmt.Printf("idx: %d ; val: %f\n", idx, v)
		idx++
		v, ok = f.Next()

		if idx > 100 {
			break
		}
	}
}

func TestSplit(t *testing.T) {
	s := Take(7, Naturals())
	l, r := Split(4, s)
	fmt.Println(Collect(l), Collect(r))

}
