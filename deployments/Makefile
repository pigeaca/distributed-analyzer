IMAGE_NAME=distributed-analyzer/api-gateway
TAG=latest

build:
	docker build -t $(IMAGE_NAME):$(TAG) -f ../services/api-gateway/Dockerfile ../

load:
	minikube image load $(IMAGE_NAME):$(TAG)

deploy:
	kubectl apply -f k8s/services/api-gateway.yaml

dev: build load deploy