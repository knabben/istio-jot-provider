push:
	docker build --tag knabben/proxy:latest . 
	docker tag knabben/proxy:latest knabben/proxy:latest
	docker push knabben/proxy:latest 

install:
	kubectl apply -f services/auth-proxy.yaml