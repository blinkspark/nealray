all: win linux
win:
	GOOS=windows GOARCH=amd64 go build -o bin/blinkcrypto.exe main.go 
linux:
	GOOS=linux GOARCH=amd64 go build -o bin/blinkcrypto main.go 