all: windows-64 linux-64
windows-64:
	GOOS=windows go build -o bin/blinkserver.exe main.go
linux-64:
	GOOS=linux go build -o bin/blinkserver main.go
clean:
	rm blinkserver