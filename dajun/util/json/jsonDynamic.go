package json

import (
	"encoding/json"
	"fmt"
	"log"
)

// 动态数据，在 msg 一个节点下
const input = `
{
	"type": "sound",
	"msg": {
		"description": "dynamite",
		"authority": "the Bruce Dickinson"
	}
}
`

// Envelope 最外层
type Envelope struct {
	Type string      `json:"type"`
	Msg  interface{} `json:"msg"`
}

// Sound 消息1
type Sound struct {
	Description string `json:"description"`
	Authority   string `json:"authority"`
}

// Cowbell 消息2
type Cowbell struct {
	More bool `json:"more"`
}

// Simple 基本用法
func Simple() {
	s := Envelope{
		Type: "sound",
		Msg: Sound{
			Description: "dynamite",
			Authority:   "the Bruce Dickinson",
		},
	}
	buf, err := json.Marshal(s)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", buf)

	c := Envelope{
		Type: "cowbell",
		Msg: Cowbell{
			More: true,
		},
	}
	buf, err = json.Marshal(c)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", buf)
}

// DynamicParse 根据类型解析
func DynamicParse() {
	// Evvelope 中有 Type 和 Msg，这里对 Msg 进行赋值，使其等于 RawMessage，目的是为了进行第二次解析
	var msg json.RawMessage
	env := Envelope{
		Msg: &msg,
	}

	// 第一次解析，目的是获取 Envelope.Type
	if err := json.Unmarshal([]byte(input), &env); err != nil {
		log.Fatal(err)
	}

	switch env.Type {
	case "sound":
		// 已经知道了类型，再对 Msg 进行第二次解析
		var s Sound
		if err := json.Unmarshal(msg, &s); err != nil {
			log.Fatal(err)
		}
		desc := s.Description
		fmt.Println(desc)
	default:
		log.Fatalf("unknown message type: %q", env.Type)
	}
}

// 动态部分是和 type 平行，而且是多个值
const input2 = `
{
	"type": "sound",
	"description": "dynamite",
	"authority": "the Bruce Dickinson"
}
`

// Envelope2 仅包含消息类型
type Envelope2 struct {
	Type string `json:"type"`
}

// DynamicSameLevel 同一级元素的动态处理
func DynamicSameLevel() {
	// 即使是同一个值，也可以被处理多次
	buf := []byte(input2)

	// 第一次解析，仅解析类型
	var env Envelope2
	if err := json.Unmarshal(buf, &env); err != nil {
		log.Fatal(err)
	}

	switch env.Type {
	case "sound":
		// 组合类型，完整的数据，而不是仅动态部分（上一个例子中，仅是动态部分）
		var s struct {
			Envelope
			Sound
		}
		if err := json.Unmarshal(buf, &s); err != nil {
			log.Fatal(err)
		}
		desc := s.Description
		fmt.Println(desc)
	default:
		log.Fatalf("unknown message type: %q", env.Type)
	}
}
