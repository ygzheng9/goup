package u8

import (
	"fmt"
	"testing"
)

func Test_findSNInfo(t *testing.T) {
	param := TraceSNParam{}
	param.StartDate = "2018-04-01"
	param.EndDate = "2018-03-20"
	param.CustCode = ""
	param.SoCode = ""
	param.SN = "KLBEVMC260 1142J30057"

	items, err := findSNInfo(param)
	if err != nil {
		t.Errorf("\n%#v\n", err)
	}

	fmt.Printf("count: %d\n", len(items))
	for _, v := range items {
		fmt.Printf("%+v", v)
	}
}
