package matShortage

import (
	"fmt"
	"strings"
	"testing"
)

func Test_calcDiff(t *testing.T) {
	s1 := "ED25-D1-02电机总成（已包装）"

	p := strings.Index(s1, "-")
	fmt.Printf("Index: %d\n ", p)
	fmt.Printf("%s\n", string([]rune(s1)[:p+3]))

	left := s1[:p+3]

	fmt.Printf("%s\n", left)

}

func Test_loadForecast(t *testing.T) {
	from := "2018-03-08"
	to := "2018-03-15"

	items1, err := loadForecast(from)
	if err != nil {
		fmt.Printf("err: %+v\n", err)
		return
	}
	fmt.Printf("items1: %d\n", len(items1))

	items2, err := loadForecast(to)
	if err != nil {
		fmt.Printf("err: %+v\n", err)
		return
	}
	fmt.Printf("items2: %d\n", len(items2))

	items3, err := loadStockIn(from, to)
	if err != nil {
		fmt.Printf("err: %+v\n", err)
		return
	}
	fmt.Printf("items3: %d\n", len(items3))

}

func Test_calcForecastDiff(t *testing.T) {
	outFile := "E:/99.localDev/easypy/u8/forecast_out.xlsx"
	from := "2018-03-08"
	to := "2018-03-15"

	CalcForecastDiff(from, to, outFile)
}

func Test_GetDuration(t *testing.T) {
	inFile := "E:/99.localDev/easypy/u8/forecast.xlsx"
	from, to, err := GetDuration(inFile)
	if err != nil {
		fmt.Printf("err: %+v\n", err)
		return
	}

	fmt.Printf("from: %s, to: %s\n", from, to)

	outFile := "E:/99.localDev/easypy/u8/forecast_out.xlsx"

	CalcForecastDiff(from, to, outFile)
}
