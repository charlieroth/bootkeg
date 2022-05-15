all:
	go build -o bk main.go && mv bk bin

clean:
	rm bin/bk
