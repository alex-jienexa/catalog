.PHONY: run build test clean migrate

run:
	go run cmd/main.go

build:
	go build -o bin/catalog-backend.exe cmd/main.go

test:
	go test ./... -v

clean:
	rm -rf bin/
	rm -f catalog.db

migrate:
	go run cmd/migrate/main.go

docker-build:
	docker build -t catalog-backend .

docker-run:
	docker run -p 8080:8080 catalog-backend