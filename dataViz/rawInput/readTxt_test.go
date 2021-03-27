package rawInput

import (
	"fmt"
	"log"
	"testing"
)

func Test_ReadPOHead(t *testing.T) {
	items, err := ReadPOHead("E:/99.localDev/tmp/1.txt")

	if err != nil {
		t.Fatal(err)
	}

	t.Logf("items: %d\n", len(items))

	for _, v := range items {
		log.Println(v)
	}
}

func Test_ReadPOItem(t *testing.T) {
	fmt.Println("start...")

	items, err := ReadPOItem("E:/99.localDev/tmp/2.txt")

	if err != nil {
		t.Fatal(err)
	}

	log.Printf("items: %d\n", len(items))

	for _, v := range items {
		t.Log(v)
		// fmt.Println(v)
	}
}
func Test_LoadPOItems(t *testing.T) {
	items, _ := LoadPOItems("E:/99.localDev/tmp/3.txt")
	fmt.Printf("items: %d \n", len(items))
}

func Test_POItemsByDate(t *testing.T) {
	items, _ := LoadPOItems("E:/99.localDev/tmp/3.txt")
	fmt.Printf("items: %d \n", len(items))

	merged := POItemsByDate(items, "2017-11-01", "2017-11-30")
	fmt.Printf("items: %d \n", len(merged))

}

func Test_loadMatByMonth(t *testing.T) {
	items, _ := loadMatByMonth("E:/99.localDev/tmp/4.txt")
	fmt.Printf("items: %d \n", len(items))
}
