package main

import (
	"log"
	"time"

	"pickup/bookingCancel"
)

func main() {
	s := time.Now()

	bookingCancel.BatchCheck()

	log.Printf("finished: %s", time.Since(s))
}
