build:
	docker build --no-cache -t example-webhook .
deploy-only:
	kubectl config use-context local
	kubectl delete -f deployment/service.yaml || true
	kubectl delete -f deployment/pod.yaml || true
	kubectl delete externaladmissionhookconfiguration example-config || true
	kubectl create -f deployment/service.yaml
	# It's necessary because the webhook needs to access the extension-apiserver-authentication configmap
	kubectl create clusterrolebinding admission-webhook --clusterrole=cluster-admin --serviceaccount=default:default 
	sleep 5
	kubectl create -f deployment/pod.yaml	
deploy: build deploy-only

