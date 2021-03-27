// 编译成 win 32 位
set CGO_ENABLED=0
set GOOS=windows
set GOARCH=386
go build cangku.go