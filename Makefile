.PHONY: build run test lint docker-build docker-run docker-stop clean

BINARY_NAME=mail-service
DOCKER_COMPOSE=docker-compose

build:
	go build -o $(BINARY_NAME) ./cmd/app

run:
	go run ./cmd/app

test:
	go test -v ./...

lint:
	golangci-lint run

docker-build:
	$(DOCKER_COMPOSE) build

docker-run:
	$(DOCKER_COMPOSE) up -d

docker-stop:
	$(DOCKER_COMPOSE) down

clean:
	go clean
	rm -f $(BINARY_NAME)

create-topic:
	docker exec kafka kafka-topics --create --topic mail.send-email --bootstrap-server kafka:29092 --partitions 1 --replication-factor 1

send-test-message:
	echo '{"user_id": "test123", "email_type": "registration"}' | docker exec -i kafka kafka-console-producer --broker-list kafka:29092 --topic mail.send-email

watch-messages:
	docker exec kafka kafka-console-consumer --bootstrap-server kafka:29092 --topic mail.send-email --from-beginning 