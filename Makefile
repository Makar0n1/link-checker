.PHONY: help infra-up infra-down logs clean

# Default target
help:
	@echo "Link Tracker - Available commands:"
	@echo ""
	@echo "  Infrastructure:"
	@echo "    make infra-up      - Start PostgreSQL and Redis"
	@echo "    make infra-down    - Stop all containers"
	@echo "    make infra-restart - Restart infrastructure"
	@echo "    make logs          - Show container logs"
	@echo "    make logs-f        - Follow container logs"
	@echo ""
	@echo "  Cleanup:"
	@echo "    make clean         - Remove containers and volumes"
	@echo ""

# =============================================================================
# Infrastructure
# =============================================================================

infra-up:
	docker-compose up -d postgres redis
	@echo "Infrastructure started. PostgreSQL: 5432, Redis: 6379"

infra-down:
	docker-compose down

infra-restart: infra-down infra-up

logs:
	docker-compose logs

logs-f:
	docker-compose logs -f

# =============================================================================
# Cleanup
# =============================================================================

clean:
	docker-compose down -v --remove-orphans
	@echo "Containers and volumes removed"
