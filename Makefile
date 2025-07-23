build:
	docker build --platform=linux/amd64 -t oliverkra/http-hello:latest .

release: build
	docker push --platform=linux/amd64 oliverkra/http-hello:latest