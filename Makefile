REMOTE_USER ?= root
REMOTE_HOST ?= 165.232.166.176
SSH_KEY ?= id_ed25519
REMOTE_ROOT ?= /var/www/ryangel

.PHONY: all build build-frontend build-backend deploy deploy-frontend deploy-backend

all: build deploy

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

deploy-backend:
	@echo "Deploying backend to $(REMOTE_HOST)..."
	scp -i $(SSH_KEY) backend/bin/server $(REMOTE_USER)@$(REMOTE_HOST):$(REMOTE_ROOT)/backend/
	@echo "Backend uploaded. You may need to restart the service with: ssh -i $(SSH_KEY) $(REMOTE_USER)@$(REMOTE_HOST) 'systemctl restart ryangel-backend'"

stop-prod-backend:
	@echo "stopping backend in $(REMOTE_HOST)..."
	ssh -i $(SSH_KEY) $(REMOTE_USER)@$(REMOTE_HOST) 'systemctl stop ryangel-backend'

start-prod-backend:
	@echo "starting backend in $(REMOTE_HOST)..."
	ssh -i $(SSH_KEY) $(REMOTE_USER)@$(REMOTE_HOST) 'systemctl start ryangel-backend'
