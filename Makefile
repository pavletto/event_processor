all: build

build:
	go build -o ./dist/app ./cmd/main.go

run:
	go run ./cmd/main.go

docker-build:
	docker-compose build

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

build-debug:
	go build -gcflags "all=-N -l" -o ./dist/app ./cmd/main.go

debug:
	dlv --listen=:2345 --headless=true --api-version=2 --continue --accept-multiclient --check-go-version=false --log exec ./dist/app

air:
	air

.PHONY: build run docker-build docker-up docker-down debug air
