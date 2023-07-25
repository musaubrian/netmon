mkdir logs
go mod tidy
go build  . -o bin\netmon.exe
cls
bin\netmon.exe
