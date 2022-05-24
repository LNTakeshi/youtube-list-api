SHELL=bash
.PHONY: install
install:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/google/ko@latest
	go install github.com/makiuchi-d/arelo@latest
	go install github.com/google/wire/cmd/wire@latest

.PHONY: gen
gen:
	wire ./...

.PHONY: build
build:
	cd react && npm run build

.PHONY: deploy-%
deploy-%:
	$(eval IMAGE := $(shell Set KO_DOCKER_REPO=asia.gcr.io/$(*)/api&& ko publish registry/api/run/main.go))
	gcloud run deploy api --image=$(IMAGE) --region asia-northeast1 --allow-unauthenticated --cpu=1

.PHONY: deploy-linux-%
deploy-linux-%:
	KO_DOCKER_REPO=gcr.io/$(*)/api gcloud run deploy api --image=`ko publish registry/api/run/main.go` --region asia-northeast1 --allow-unauthenticated --cpu=1

.PHONY: start
start:
	arelo go run ./registry/api/local/main.go

.PHONY: npm-update
npm-update:
	cd react && npx npm-check-updates -u

.PHONY: npm-install
npm-install:
	cd react && npm install