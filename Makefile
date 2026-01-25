REMOTE_USER ?= root
REMOTE_HOST ?= 165.232.166.176
SSH_KEY ?= id_ed25519
REMOTE_ROOT ?= /var/www/ryangel
BACKEND_DEV_TMUX_SESSION ?= backend-dev

.PHONY: all build build-frontend build-backend deploy deploy-frontend deploy-backend dev-frontend stop-prod-backend start-prod-backend

all: build deploy

dev: dev-frontend dev-backend

dev-frontend:
	@echo "Starting frontend development server..."
	cd frontend && npm install && npm run dev

dev-backend:
	@echo "Starting backend development server in tmux..."
	@if tmux has-session -t "$(BACKEND_DEV_TMUX_SESSION)" 2>/dev/null; then \
		echo "Session '$(BACKEND_DEV_TMUX_SESSION)' exists."; \
		tmux send-keys -t "$(BACKEND_DEV_TMUX_SESSION)" C-c ' make run' Enter; \
	else \
		echo "Session '$(BACKEND_DEV_TMUX_SESSION)' does not exist. Creating a new one."; \
		tmux new-session -d -s "$(BACKEND_DEV_TMUX_SESSION)"; \
		tmux send-keys -t "$(BACKEND_DEV_TMUX_SESSION)" C-c ' make run' Enter; \
	fi

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
