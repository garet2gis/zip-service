.PHONY: build
build:
	docker build -t zip .

.PHONY: run
run:
	docker run -dp 8080:8080 \
	-v ${PWD}/root:/root \
  	--rm \
  	--memory=512m \
   	--name ziptest \
  	zip