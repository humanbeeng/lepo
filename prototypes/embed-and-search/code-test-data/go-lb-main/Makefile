build:
	go build -o ./bin/server

server: build
	./bin/lb

run: build server s1 s2

s1: 
	cd /home/personal/vscode-workspace/go/dummy-server && go run . -id=1 -addr=localhost:5001

s2:
	cd /home/personal/vscode-workspace/go/dummy-server && go run . -id=2 -addr=localhost:5002
.PHONY: s1 s2 run build