.PHONY: build test run

build:
	go build -o server .

test:
	go test ./...

# ローカル動作確認
run:
	docker build --platform linux/amd64 -t g2 .
	docker run --rm -p 8080:8080 g2
