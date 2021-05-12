all: win linux
win:
	GOOS=windows GOARCH=amd64 go build -o bin/blinkcryptor.exe main.go 
linux:
	GOOS=linux GOARCH=amd64 go build -o bin/blinkcryptor main.go 