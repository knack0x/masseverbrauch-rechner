.PHONY: help build build-api build-web up down test clean logs deploy dev-api dev-web

# Project paths
PROJECT_ROOT := $(shell pwd)
API_DIR := $(PROJECT_ROOT)/api
WEB_DIR := $(PROJECT_ROOT)/web

help:
	@echo "Masseverbrauch-Rechner - Available commands:"
	@echo ""
	@echo "  build       - Build all Docker images"
	@echo "  build-api   - Build Go API binary"
	@echo "  build-web   - Build PHP Docker image"
	@echo "  up          - Start all services with docker-compose"
	@echo "  down        - Stop all services"
	@echo "  test        - Test the API endpoint"
	@echo "  test-web    - Test the full stack"
	@echo "  logs        - View logs from all services"
	@echo "  clean       - Remove containers and images"
	@echo "  deploy      - Deploy to Mac Studio"
	@echo "  dev-api     - Run Go API locally (without Docker)"
	@echo "  dev-web     - Run PHP web locally (without Docker)"
	@echo ""

build: build-api build-web

build-api:
	@echo "Building Go API..."
	cd $(API_DIR) && docker build -t masseverbrauch-api .

build-web:
	@echo "Building PHP Web..."
	cd $(WEB_DIR) && docker build --no-cache -t masseverbrauch-web .

up:
	@echo "Starting services..."
	docker-compose up -d --build
	@echo "API running at http://localhost:50570"
	@echo "Web running at http://localhost:50571"

down:
	@echo "Stopping services..."
	docker-compose down

test:
	@echo "Testing API..."
	curl -X POST http://localhost:50570/api/calculate \
		-H "Content-Type: application/json" \
		-d '{"flow":1.5,"runtime_minutes":10,"slots":[{"before":100,"after":85},{"before":50,"after":42}]}' \
		| python3 -m json.tool || cat

test-web: up
	@echo "Waiting for services to start..."
	@sleep 3
	@echo "Testing full stack..."
	@curl -s http://localhost:50571 | head -20

logs:
	docker-compose logs -f

clean: down
	@echo "Cleaning up..."
	docker rmi masseverbrauch-api masseverbrauch-web 2>/dev/null || true
	docker system prune -f

deploy:
	@echo "Deploying to Mac Studio..."
	@echo "NOTE: Update DEPLOY_HOST with your Mac Studio IP/hostname"
	@echo "Run: make deploy DEPLOY_HOST=your-mac-studio.local"
	if [ -z "$(DEPLOY_HOST)" ]; then \
		echo "Error: DEPLOY_HOST not set. Usage: make deploy DEPLOY_HOST=hostname"; \
		exit 1; \
	fi
	ssh $(DEPLOY_HOST) "mkdir -p ~/masseverbrauch-rechner"
	scp -r $(PROJECT_ROOT)/docker-compose.yml $(DEPLOY_HOST):~/masseverbrauch-rechner/
	scp -r $(API_DIR) $(DEPLOY_HOST):~/masseverbrauch-rechner/
	scp -r $(WEB_DIR) $(DEPLOY_HOST):~/masseverbrauch-rechner/
	ssh $(DEPLOY_HOST) "cd ~/masseverbrauch-rechner && docker-compose down && docker-compose up -d"
	@echo "Deployment complete!"

# For local Go development (without Docker)
dev-api:
	@echo "Running API locally..."
	cd $(API_DIR) && go run main.go

# For local PHP development (without Docker)
dev-web:
	@echo "Running PHP web locally..."
	cd $(WEB_DIR) && API_URL=http://localhost:50570/api/calculate php -S localhost:50571

# Quick restart
restart: down up
