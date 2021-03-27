package bookingCancel

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type Stat struct {
	Positive int
	Negative int
}

func (s Stat) String() string {
	return fmt.Sprintf("%5.2d / %5.2d", s.Positive, s.Negative)
}

func ProcessData() {
	infile, err := os.Open(inputFile)
	if err != nil {
		panic(err)
	}
	defer infile.Close()

	dumpOutput(infile)
}

// 通过循环，检查多列的值
func BatchCheck() {
	cols := []Index{CGO_TRADE_LANE_CDE, ITS_TRADE_LANE_CDE, CGO_TRADE_CDE, BKG_OFCE_CDE, ITS_TRADE_CDE, CNTR_TYPE, TEU}

	//cols := []Index{CONTROL_UUID, SHIPPER_UUID, CONSIGN_UUID, NOTIFY_UUID, FORWARD_UUID, LAST_APPROVER}

	//cols := []Index{SHIPPER_UUID}
	for _, col := range cols {
		checkData(col)
	}
}

func BatchCheckParallel() {
	//cols := []Index{CGO_TRADE_LANE_CDE, ITS_TRADE_LANE_CDE, CGO_TRADE_CDE, BKG_OFCE_CDE, ITS_TRADE_CDE, CNTR_TYPE, TEU}

	cols := []Index{CONTROL_UUID, SHIPPER_UUID, CONSIGN_UUID, NOTIFY_UUID, FORWARD_UUID, LAST_APPROVER}

	//cols := []Index{SHIPPER_UUID}

	c := make(chan map[string]Stat, len(cols))

	for _, col := range cols {
		go checkDataParallel(col, c)
	}

	var out bytes.Buffer
	for i := 0; i < len(cols); i++ {
		values := <-c

		for k, v := range values {
			fmt.Fprintf(&out, "%q = %s\n", strings.TrimSpace(k), v)
		}

		fmt.Fprintln(&out, "--------------------------------------")
	}

	fmt.Println(out.String())
}

// 检查某一列的值，分析某一列的值，与目标列的对应关系
func checkData(i Index) {
	infile, err := os.Open(inputFile)
	if err != nil {
		panic(err)
	}
	defer infile.Close()

	values := make(map[string]Stat)

	scanner := bufio.NewScanner(infile)
	for scanner.Scan() {
		cols := strings.Split(scanner.Text(), sep)

		k := cols[i]
		a := values[k]

		if _, ok := TargetMap[cols[COMMIT_STATUS]]; ok {
			a.Positive = a.Positive + 1
		} else {
			a.Negative = a.Negative + 1
		}

		values[k] = a
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	var o bytes.Buffer
	for k, v := range values {
		fmt.Fprintf(&o, "%q = %s\n", strings.TrimSpace(k), v)
	}

	fmt.Println("--------------------------------------")
	fmt.Println(o.String())
}

func checkDataParallel(i Index, c chan map[string]Stat) {
	infile, err := os.Open(inputFile)
	if err != nil {
		panic(err)
	}
	defer infile.Close()

	values := make(map[string]Stat)

	scanner := bufio.NewScanner(infile)
	for scanner.Scan() {
		cols := strings.Split(scanner.Text(), sep)

		k := cols[i]
		a := values[k]

		if _, ok := TargetMap[cols[COMMIT_STATUS]]; ok {
			a.Positive = a.Positive + 1
		} else {
			a.Negative = a.Negative + 1
		}

		values[k] = a
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	c <- values
}

// 读入文件，对每一行进行处理，然后再写入新文件
func dumpOutput(r io.Reader) {
	outfile, err := os.Create(outputFile)
	if err != nil {
		panic(err)
	}
	defer outfile.Close()

	b := bufio.NewWriter(outfile)

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()

		_, err = b.ReadFrom(strings.NewReader(parseRecord(line)))
		if err != nil {
			panic(err)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	if err = b.Flush(); err != nil {
		panic(err)
	}
}

// 对每一行，先进行拆分，再对某些列做逻辑处理，处理的结果作为新列，加在最后
// 对每一行，在最后增加一个 换行符
func parseRecord(s string) string {
	fields := strings.Split(s, sep)

	fields = append(fields, calcStatus(fields[COMMIT_STATUS]))

	fields = append(fields, "\n")

	return strings.Join(fields, sep)
}

// 没有 set ，但是通过 map 可以模拟出相似的效果
func calcStatus(s string) string {
	if _, ok := TargetMap[s]; ok {
		return "1"
	} else {
		return "0"
	}
}

func ScanNumber() {
	// An artificial input source.
	const input = "1234  1234567901234567890 5678"
	scanner := bufio.NewScanner(strings.NewReader(input))

	// Create a custom split function by wrapping the existing ScanNumber function.
	split := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		advance, token, err = bufio.ScanWords(data, atEOF)
		if err == nil && token != nil {
			_, err = strconv.ParseInt(string(token), 10, 32)
		}
		return
	}

	// Set the split function for the scanning operation.
	scanner.Split(split)

	// Validate the input
	for scanner.Scan() {
		fmt.Printf("%s\n", scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Invalid input: %s", err)
	}
}
