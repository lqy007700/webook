.PHONY: docker
docker:
	@rm webook || true
	@GOOS=linux GOARCH=arm go build -tags=k8s -o webook .
	@docker rmi -f lqy007700/webook:v0.0.1
	@docker build -t lqy007700/webook:v0.0.1 .