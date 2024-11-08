build-ui:
	cd ui && yarn build

build-go:
	go build -o bin/goggle .

build: build-ui build-go
