build-ui:
	cd ui && yarn build

build-go:
	go build -o bin/goggle .

test-go:
	@./scripts/test-go . --strict

test-go-cover:
	@./scripts/test-go .
	@go tool cover -html=coverage.out

generate-mock:
	@./scripts/genmock .

build: build-ui build-go

watch:
	@gow run ./cmd/goggle -debug | ./scripts/logform -p

