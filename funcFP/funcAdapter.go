package funcFP

import "log"

// CatMiao 具体的函数
func CatMiao(s string) {
	log.Printf("Cat miao miao: %q", s)
}

// DogWong 具体的函数
func DogWong(s string) {
	log.Printf("Dog wong wong: %q", s)
}

// Voicer 定义一个通用接口，只有一个方法 Bark
type Voicer interface {
	Bark(s string)
}

// VoicerFunc 是一个 function adapter: 只要参数、返回值相同的函数，都可以转换为同一个函数
// VoicerFunc 是一个 type，所以可以用作类型转换使用:
// a := VoicerFunc(afunc)
// 转换后，可以执行 a.Bark("msg")
type VoicerFunc func(s string)

// Bark 使得 VoicerFunc 满足 Voicer interface 的要求
// go 的惯用手法，为一个 type 附加方法，是原本的 type 满足某个 interface
func (f VoicerFunc) Bark(s string) {
	// f 实质上是一个函数，所以可以直接调用
	f(s)
}

// VoicerNum go 中的惯用手法，目的是增强现有类型的表现
type VoicerNum int

// Bark 相当于为 int 增加了 Bark 方法，使得 int 可以作为 Voicer 来使用
func (i VoicerNum) Bark(s string) {
	// 因为 VoicerNum 实质是 int，不是 func，所以不能像函数一样使用 i(s)
	log.Printf("Bark %q for %d times. \n", s, i)
}

// ShowUp 参数是 interface，而不是 struct
func ShowUp(v Voicer, s string) {
	v.Bark(s)
}
