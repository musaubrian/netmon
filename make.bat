mkdir -p logs
go mod tidy
go build  .
cls
netmon.exe
