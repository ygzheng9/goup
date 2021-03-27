package bookingCancel

import (
	"log"
	"testing"
	"time"
)

func TestLoadData(t *testing.T) {

	s := time.Now()

	//ProcessData()

	BatchCheck()

	log.Printf("finished: %s", time.Since(s))

	//ScanNumber()
}

func TestLoadDataParallel(t *testing.T) {
	s := time.Now()

	BatchCheckParallel()

	log.Printf("finished: %s", time.Since(s))
}
