# Docker Compose Setup
# ====================
# Start, stop and watch docker processes for local and integration testing.
.PHONY: test cover up down status check wait run

# providers: wm, bn (Java), rpk (Go)
KAFKA_PROVIDER  = rpk.v23
KAFKA_HOST      = localhost
KAFKA_PORT      = 9092
KAFKA_BROKER    = $(KAFKA_HOST):$(KAFKA_PORT)
COMPOSE_FILE    = compose-$(KAFKA_PROVIDER).yml
WAIT_TIME       = 15
HEATH_CHECK     = nc -z $(KAFKA_HOST) $(KAFKA_PORT)
NUM_EVENTS      = 10
TIMEOUT         = 10s
ARGS            = -demo

test:
	go test -race ./...

cover:
	../../scripts/cover.sh ./...

up:     ; docker-compose -f $(COMPOSE_FILE) up --no-recreate -d
down:   ; docker-compose -f $(COMPOSE_FILE) down --remove-orphans
status: ; docker-compose -f $(COMPOSE_FILE) ps
check:  ; $(HEATH_CHECK) -G 1
wait:
	# waiting for brokers...
	for i in `seq $(WAIT_TIME)`; do printf "."; $(HEATH_CHECK) && break || false; done

run:
	go run -race cmd/buffertest/main.go -c $(NUM_EVENTS) -b $(KAFKA_BROKER) -d $(TIMEOUT)

.PHONY: run-tables
run-tables: bin/tables
	go build -o $^ cmd/tables/main.go && $^ $(ARGS)
