# ---------------------------------------
# Makefile for minimal_sns_app project
# ---------------------------------------

# Docker Compose base commands
COMPOSE=docker compose -f docker-compose.yml
TEST_COMPOSE=$(COMPOSE) -f docker-compose.test.yml

up:
	$(COMPOSE) up -d

down:
	$(COMPOSE) down

restart:
	$(COMPOSE) down && $(COMPOSE) up -d

test-up:
	$(TEST_COMPOSE) up -d

test-down:
	$(TEST_COMPOSE) down -v

test-restart:
	$(TEST_COMPOSE) down -v && $(TEST_COMPOSE) up -d

test:
	$(COMPOSE) exec app go test ./test -v


log:
	$(COMPOSE) logs app

.PHONY: up down restart test-up test-down test-restart test mysql log build
