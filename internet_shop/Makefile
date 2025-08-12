.PHONY: run
run:
	go run app/cmd/app/main.go

.PHONY: lint
lint:
	golangci-lint run

swagger:
	swag init -g ./app/cmd/app/main.go -o app/docs

.PHONY: clean
clean:
	rm -rf ./app/build || true

.PHONY: up-local-env
up-local-env: down-local-env
	@docker-compose -f docker-compose.local.yml up -d

.PHONY: down-local-env
down-local-env:
	@docker-compose -f docker-compose.local.yml stop