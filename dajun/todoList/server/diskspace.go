package server

// DiskSpace 磁盘空间容量
type DiskSpace struct {
	Disk     string
	TotalAmt float64
	FreeAMt  float64
}

// Real 注册的服务
type Real int

// GetDiskSpace windows 下取得磁盘的容量信息
func (t *Real) GetDiskSpace(disk string, space *DiskSpace) error {
	if len(disk) == 0 {
		return nil
	}

	d := disk
	if len(disk) == 1 {
		d = d + ":"
	}

	// h := syscall.MustLoadDLL("kernel32.dll")
	// c := h.MustFindProc("GetDiskFreeSpaceExW")
	// lpFreeBytesAvailable := int64(0)
	// lpTotalNumberOfBytes := int64(0)
	// lpTotalNumberOfFreeBytes := int64(0)
	// _, _, _ = c.Call(uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(d))),
	// 	uintptr(unsafe.Pointer(&lpFreeBytesAvailable)),
	// 	uintptr(unsafe.Pointer(&lpTotalNumberOfBytes)),
	// 	uintptr(unsafe.Pointer(&lpTotalNumberOfFreeBytes)))

	// G := 1024 * 1024 * 1024 * 1.0

	// space.Disk = disk
	// space.FreeAMt = float64(lpFreeBytesAvailable) / G
	// space.TotalAmt = float64(lpTotalNumberOfBytes) / G

	// now := time.Now().Format("2006-01-02 15:04:05")
	// fmt.Printf("GetDiskSpace %s: %+v\n", now, space)

	return nil
}
