default: get binary 

fmt:
	gofmt -s -w .

get:
	go get ./...

test:
	go test -v ./...

fmt-test:
	gofmt -l . | wc -c | grep -E ^0$

binary:
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' .

clean:
	rm -f kaiif 
