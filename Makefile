.PHONY: build test run help

help: ## ヘルプを表示
	@grep -E '^[a-zA-Z_-]+:.*##' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*##"}; {printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2}'

build: ## バイナリをビルド (./server)
	go build -o server .

test: ## テストを実行
	go test ./...

run: ## Docker でローカル起動 (:8080)
	docker build --platform linux/amd64 -t g2 .
	docker run --rm -p 8080:8080 g2
