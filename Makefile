.PHONY: build test run deploy

build:
	go build -o bootstrap .

test:
	go test ./...

# Lambda Runtime Interface Emulator (RIE) でローカル動作確認
run:
	docker build --platform linux/amd64 -t g2-lambda .
	docker run --rm -p 9000:8080 g2-lambda

# 別ターミナルで実行: make invoke
invoke:
	curl -s -XPOST "http://localhost:9000/2015-03-31/functions/function/invocations" \
	  -d '{"name": "Lambda"}' | jq .

deploy:
	chmod +x deploy.sh
	./deploy.sh
