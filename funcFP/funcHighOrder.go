package funcFP

// 第一个例子

// Expr 是一个类型，代表一个函数，该函数没有参数，返回一个 float64
type Expr func() float64

// Const 是一个 high order function，也即：返回值是一个函数，类型是 Expr
// Const 自身是一个函数，参数是 float64；
func Const(x float64) Expr {
	// 定义一个函数，没有参数，返回值是 float64
	a := func() float64 {
		return x
	}

	// 返回这个函数
	return a

	// 上述逻辑可以简化成如下形式，两者表达的意思是相同的
	// return func() float64 {
	// 	return x
	// }
}

// Sum 参数是一个不定长的数组，每个元素都是 Expr，
// 返回值，是一个函数，类型也是 Expr, 也即：
// a := Sum(..)  // a 其实是一个 Expr，不是一个数，
// b := a() // a 执行后，b 才是一个数字
func Sum(exs ...Expr) Expr {
	// 返回值类型是 Expr，所以这里 return 的就是这样的类型
	return func() float64 {
		sum := 0.0

		// 参数不定长度，使用 range 遍历
		for _, ex := range exs {
			// 每一个元素都是一个 Expr，本质是一个函数，所以可以执行调用
			sum += ex()
		}

		// 这里 sum 是一个 float64，和函数返回值类型相同
		return sum
	}
}

// 第二个例子

// Seq 是一个协议，只要提供 Next() 就可以当做 Seq 来使用
type Seq interface {
	Next() (x float64, ok bool)
}

// SeqFunc 是一个新的 type，可以直接作为 类型转换 的运算；
// 其本质是一个 函数，该函数没有参数，返回两个值
type SeqFunc func() (float64, bool)

// Next 使得 SeqFunc 符合 interface Seq 的要求，也即：任何一个 SeqFunc 都可以作为 Seq 使用
func (f SeqFunc) Next() (float64, bool) {
	return f()
}

// Mk 返回值是一个 interface，具有 Next() 方法
func Mk(vals ...float64) Seq {
	i := 0

	// 这里是函数定义，该函数没有执行；
	// 在调用返回结果的 Next() 时，才会执行；执行时，每次读取的 vals，i 都是同一个对象
	f1 := func() (result float64, ok bool) {
		if i >= len(vals) {
			return 0, false
		}

		result = vals[i]
		i++
		return result, true
	}

	// SeqFunc 相当于一个类型转换的操作符，把 f1 变成了 SeqFunc 类型；
	// SeqFunc 定义了 Next()，这样 SeqFunc 就满足了 interface Seq;
	f2 := SeqFunc(f1)

	return f2
}

// Collect 把一个 Seq 变成 slice
//
func Collect(s Seq) []float64 {
	var vals []float64
	for {
		// 调用了 Next()，而这个 Next() 执行 SeqFunc(f) 中的 f；而且只执行一次
		x, ok := s.Next()
		if !ok {
			break
		}
		vals = append(vals, x)
	}

	return vals
}

// Naturals 无穷数列
func Naturals() Seq {
	n := 0.0
	return SeqFunc(func() (x float64, ok bool) {
		n++
		return n, true
	})

}

// Take 取前 n 个，这里返回的是 Seq，而不是取出来的数；要得到具体的数，还需要调用 Next()
func Take(n int, s Seq) Seq {
	f1 := func() (float64, bool) {
		if n <= 0 {
			return 0, false
		}
		n--

		// lazy evaluation: 一直到返回的 Seq 的 Next() 被调用时，才会执行
		return s.Next()
	}

	f2 := SeqFunc(f1)

	return f2
}

// Split 一个拆两个
func Split(n int, s Seq) (l, r Seq) {
	return Mk(Collect(Take(n, s))...), s
}
