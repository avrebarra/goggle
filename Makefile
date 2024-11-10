build-ui:
	cd ui && yarn build

build-go:
	go build -o bin/goggle .

generate_mock:
	./scripts/genmock .

build: build-ui build-go

