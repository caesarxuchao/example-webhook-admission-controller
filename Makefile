all:
	docker build --no-cache -t example-webhook .
deploy-only:
	kubectl config use-context local
	kubectl delete -f deployment/service.yaml || true
	kubectl delete -f deployment/pod.yaml || true
	kubectl create -f deployment/service.yaml
	sleep 5
	kubectl create -f deployment/pod.yaml	
deploy: all deploy-only

