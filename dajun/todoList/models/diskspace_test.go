package models

import (
	"fmt"
	"testing"
	"time"
)

func Test_diskspaceInsert(t *testing.T) {
	item := CreateDiskSpace()

	item.DiskName = "localhost"
	item.ServerIP = "10.10.10.18"
	item.DiskName = "C"
	item.WarningPoint = 75
	item.NoticeUser = "babab"
	item.Status = "ACTIVE"

	err := item.Insert()
	if err != nil {
		t.Errorf("insert failed: %+v\n", err)
		return
	}
}

func Test_diskspaceUpdate(t *testing.T) {
	item, err := DiskSpaceFindByID(5)
	if err != nil {
		t.Errorf("DiskSpaceFindByID failed: %+v\n", err)
		return
	}

	item.LastFreeAmt = 7.234234
	item.LastTotalAmt = 60.000
	item.LastTick = time.Now().Format("2006-01-02 15:04:05")
	err = item.UpdateSpace()
	if err != nil {
		t.Errorf("UpdateSpace failed: %+v\n", err)
		return
	}

	item.Status = "DEACTIVE"
	err = item.UpdateStatus()
	if err != nil {
		t.Errorf("UpdateStatus failed: %+v\n", err)
		return
	}
}

func Test_diskspaceSelect(t *testing.T) {
	items, err := DiskSpaceFindAll()
	if err != nil {
		t.Errorf("DiskSpaceFindAll failed: %+v\n", err)
		return
	}
	for _, v := range items {
		fmt.Printf("%+v\n", v)
	}
}

func Test_diskspaceMonitor(t *testing.T) {
	err := DiskSpaceMonitor()
	if err != nil {
		t.Errorf("DiskSpaceMonitor failed: %+v\n", err)
		return
	}
}
