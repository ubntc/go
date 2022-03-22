# Docker Compose Setup
# ====================
# Start, stop and watch docker processes for local and integration testing.
.PHONY: test cover up down status check wait run

BN_FILE = docker-compose-bitnami.yml
WM_FILE = docker-compose-wurstmeister.yml
COMPOSE_FILE = $(BN_FILE)
WAIT_TIME = 15
HEATH_CHECK = nc -z localhost 9092

test:
	go test -race ./...

cover:
	../scripts/cover.sh ./...

up:     ; docker-compose -f $(COMPOSE_FILE) up --no-recreate -d
down:   ; docker-compose -f $(COMPOSE_FILE) down --remove-orphans
status: ; docker-compose -f $(COMPOSE_FILE) ps
check:  ; $(HEATH_CHECK)
wait:
	# waiting for brokers...
	for i in `seq $(WAIT_TIME)`; do \
		printf "."; $(HEATH_CHECK) && break || sleep 1; false; done

run:
	go run -race cmd/buffertest/main.go -c 10
