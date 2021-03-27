package u8

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// OutboundReply api 的返回结果
type OutboundReply struct {
	RtnCode    int    `json:"rtnCode"`
	OrderNum   string `json:"orderNum"`
	ItemCnt    int    `json:"itemCnt"`
	UploadCnt  int    `json:"uploadCnt"`
	MatchedCnt int    `json:"matchedCnt"`
}

// ReadOutboundFile 从文件中读取数据, 然后 post 给 api
func ReadOutboundFile(infile string) (OutboundReply, error) {
	// 解析 post 的返回结果
	reply := OutboundReply{}

	file, err := os.Open(infile)
	defer file.Close()
	if err != nil {
		fmt.Printf("\n open file: %#v\n", err)
		return reply, err
	}

	orderNum, items, err := ParseUploadOutbound(file)
	if err != nil {
		fmt.Printf("\n parseUploadOutbound: %#v\n", err)
		return reply, err
	}

	// 准备 post 参数
	param := OutboundUploadSvcParam{
		OrderNum: orderNum,
		Items:    items,
	}
	data, _ := json.Marshal(param)

	// post 请求
	resp, err := http.Post("http://server22:8000/api/u8/outboundUploadSvc", "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Printf("post: %+v\n", err)
		return reply, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	// fmt.Printf("body: %s\n", body)

	err = json.Unmarshal([]byte(body), &reply)
	if err != nil {
		fmt.Printf("reply err: %+v\n", err)
		return reply, err
	}
	// fmt.Printf("ok: %+v\n", reply)

	return reply, nil
}
