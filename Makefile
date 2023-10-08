SHELL=bash
PROJECT_ID=$(shell gcloud config get core/project)

.PHONY: install
install:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/google/ko@latest
	go install github.com/makiuchi-d/arelo@latest
	go install github.com/google/wire/cmd/wire@latest

.PHONY: gen
gen:
	wire ./...

.PHONY: docker-build
docker-build:
	cd infra/youtube-dl/ && docker-compose build

.PHONY: docker-up
docker-up:
	cd infra/youtube-dl/ && docker-compose up -d

.PHONY: build
build:
	cd react && npm run build

.PHONY: docker-deploy
docker-deploy:
	gcloud auth configure-docker asia-northeast1-docker.pkg.dev
	docker build -f infra/youtube-dl/Dockerfile -t asia-northeast1-docker.pkg.dev/$(PROJECT_ID)/api/latest .
	docker push asia-northeast1-docker.pkg.dev/$(PROJECT_ID)/api/latest
	@gcloud run deploy api --image=asia-northeast1-docker.pkg.dev/$(PROJECT_ID)/api/latest:latest --region asia-northeast1 --allow-unauthenticated --cpu=1 \
	--set-env-vars=YOUTUBE_API_KEY=$(shell gcloud secrets versions access latest --secret=youtube-api-key) \
	--set-env-vars=REDIS_ADDR=$(shell gcloud secrets versions access latest --secret=redis-addr) \
	--set-env-vars=REDIS_PASSWORD=$(shell gcloud secrets versions access latest --secret=redis-password) \
	--set-env-vars=SPOTIFY_CLIENT_ID=$(shell gcloud secrets versions access latest --secret=spotify-client-id) \
	--set-env-vars=SPOTIFY_CLIENT_SECRET=$(shell gcloud secrets versions access latest --secret=spotify-client-secret)

.PHONY: start
start:
	arelo go run ./registry/api/local/main.go

.PHONY: npm-update
npm-update:
	cd react && npx npm-check-updates -u

.PHONY: npm-install
npm-install:
	cd react && npm install