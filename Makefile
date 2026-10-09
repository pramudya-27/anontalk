SHELL := /bin/bash

COMPOSE := docker compose -f backend/docker-compose.yml
DB_URL := root:root@tcp(127.0.0.1:3306)/anontalk

.PHONY: help setup services dev frontend backend build lint clean

help: ## Tampilkan perintah yang tersedia
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-12s\033[0m %s\n", $$1, $$2}'

setup: ## Siapkan dependency dan konfigurasi frontend lokal
	npm ci
	cd backend && go mod download
	@if [ ! -f .env.local ]; then printf 'NEXT_PUBLIC_API_URL=http://localhost:8080\n' > .env.local; fi
	@echo "Setup selesai. Jalankan: make dev"

services: ## Jalankan MySQL dan Redis menggunakan Docker
	$(COMPOSE) up -d mysql redis

dev: services ## Jalankan backend Go dan frontend Next.js
	@$(MAKE) -j2 frontend backend

frontend: ## Jalankan frontend pada http://localhost:3000
	@if [ ! -x node_modules/.bin/next ]; then npm ci; fi
	npm run dev

backend: ## Jalankan backend pada http://localhost:8080
	cd backend && DATABASE_URL='$(DB_URL)' JWT_SECRET='local-only-change-me' REDIS_ADDR='127.0.0.1:6380' SERVER_PORT=8080 go run main.go

build: ## Build frontend untuk produksi
	npm run build

lint: ## Jalankan ESLint
	npm run lint

clean: ## Hentikan MySQL dan Redis lokal
	$(COMPOSE) down
