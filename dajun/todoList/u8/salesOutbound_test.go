package u8

import (
	"fmt"
	"os"
	"testing"
)

func Test_fetchOutboundItems(t *testing.T) {
	items, err := fetchOutboundItems("7158")
	if err != nil {
		t.Errorf("\n%#v\n", err)
	}

	fmt.Printf("count: %d\n", len(items))
	for _, v := range items {
		fmt.Printf("%+v", v)
	}
}

func Test_updateOutboundSeq(t *testing.T) {
	autoID := 1001899546

	err := updateOutboundSeq(autoID, "hahahah")
	if err != nil {
		t.Errorf("\n updateOutboundSeq: %#v\n", err)
	}
}

func Test_matchOutboundItems(t *testing.T) {
	items, err := fetchOutboundItems("7158")
	if err != nil {
		t.Errorf("\n%#v\n", err)
	}

	newItems := []OutboundT{
		{InvCode: "D15-C10700015-000", SeqNo: "aaaaaa", Matched: 0},
		{InvCode: "D15-C10700015-000", SeqNo: "bbbbbb", Matched: 0},
		{InvCode: "D15-C10700015-000", SeqNo: "cccccc", Matched: 0},
		{InvCode: "D15-C10700014-000", SeqNo: "dddddd", Matched: 0},
		{InvCode: "D15-C10700014-000", SeqNo: "eeeeee", Matched: 0},
	}

	cnt, err := matchOutboundItems(items, newItems)
	if err != nil {
		t.Errorf("\n updateOutboundSeq: %#v\n", err)
	}

	fmt.Printf("matched: %d", cnt)
}

func Test_parseUploadOutbound(t *testing.T) {
	file, err := os.Open("E:/temp/Book2.xlsx")
	defer file.Close()
	if err != nil {
		t.Errorf("\n open file: %#v\n", err)
	}

	orderNum, items, err := ParseUploadOutbound(file)
	if err != nil {
		t.Errorf("\n parseUploadOutbound: %#v\n", err)
	}

	fmt.Printf("orderNum: %s\n", orderNum)
	for _, v := range items {
		fmt.Printf("%+v", v)
	}
}
