REMOTE_USER ?= root
REMOTE_HOST ?= 165.232.166.176
SSH_KEY ?= id_ed25519
REMOTE_ROOT ?= /var/www/ryangel

.PHONY: all build build-frontend build-backend deploy deploy-frontend deploy-backend dev-frontend stop-prod-backend start-prod-backend

all: build deploy

dev: dev-frontend dev-backend

dev-frontend:
	@echo "Starting frontend development server..."
	cd frontend && npm install && npm run dev

dev-backend:
	@echo "Starting backend development server..."
	cd $(CURDIR)/backend && go run ./cmd/server

build: build-frontend build-backend

build-frontend:
	@echo "Building frontend..."
	cd frontend && npm install && npm run build

build-backend:
	@echo "Building backend..."
	cd backend && GOOS=linux GOARCH=amd64 go build -o bin/server ./cmd/server

deploy: deploy-frontend deploy-backend

deploy-frontend:
	@echo "Deploying frontend to $(REMOTE_HOST)..."
	scp -r -i $(SSH_KEY) frontend/dist/* $(REMOTE_USER)@$(REMOTE_HOST):$(REMOTE_ROOT)/frontend/

deploy-backend: stop-prod-backend upload-prod-backend start-prod-backend

upload-prod-backend:
	@echo "Deploying backend to $(REMOTE_HOST)..."
	scp -i $(SSH_KEY) backend/bin/server $(REMOTE_USER)@$(REMOTE_HOST):$(REMOTE_ROOT)/backend/
	scp -i $(SSH_KEY) backend/.env.prod $(REMOTE_USER)@$(REMOTE_HOST):$(REMOTE_ROOT)/backend/.env.prod
	scp -i $(SSH_KEY) backend/ryangel-backend.service $(REMOTE_USER)@$(REMOTE_HOST):/etc/systemd/system/ryangel-backend.service
	@echo "Backend uploaded. You may need to restart the service with: ssh -i $(SSH_KEY) $(REMOTE_USER)@$(REMOTE_HOST) 'systemctl daemon-reload && systemctl restart ryangel-backend'"

restart-prod-backend: stop-prod-backend start-prod-backend

stop-prod-backend:
	@echo "stopping backend in $(REMOTE_HOST)..."
	ssh -i $(SSH_KEY) $(REMOTE_USER)@$(REMOTE_HOST) 'systemctl stop ryangel-backend'

start-prod-backend:
	@echo "starting backend in $(REMOTE_HOST)..."
	ssh -i $(SSH_KEY) $(REMOTE_USER)@$(REMOTE_HOST) 'systemctl start ryangel-backend'
