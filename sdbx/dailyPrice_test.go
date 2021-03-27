package sdbx

import (
	"bytes"
	"fmt"
	"log"
	"testing"
	"time"
)

func Test_getDailyPrice(t *testing.T) {
	s := time.Now()

	items := getDailyPrice("SH600491")
	for _, v := range items {
		fmt.Println(v)
	}

	log.Printf("finished: %s", time.Since(s))
}

func Test_getPriceTable(t *testing.T) {
	// items := getPriceTable([]string{"SH600491", "SH601398", "SZ000985"})
	// items := getPriceTable([]string{"SH603288", "SZ000651", "SH600535"})
	items := getPriceTable([]string{"SH603288", "SZ002230", "SH600535"})

	for _, v := range items {
		var buf bytes.Buffer
		for _, t := range v {
			buf.WriteString(t + " ")
		}

		fmt.Println(buf.String())
	}
}
