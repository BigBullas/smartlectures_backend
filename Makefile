.PHONY: build down prod recognizer storage


build:
	@docker-compose up --build -d
	
down:
	@docker-compose down

prod:
	@mv ./docker-compose.yml ./docker-compose-dev.yml 
	@cp ./prod/docker-compose.yml ./docker-compose.yml
	@docker compose up --build -d
	@rm ./docker-compose.yml
	@mv ./docker-compose-dev.yml ./docker-compose.yml  

build-prod-ml-search:
	@mv ./docker-compose.yml ./docker-compose-dev.yml 
	@cp ./prod/docker-compose.ml.search.yml ./docker-compose.yml
	@docker compose up --build -d
	@rm ./docker-compose.yml
	@mv ./docker-compose-dev.yml ./docker-compose.yml  
	
prod-ml-search:
	@mv ./docker-compose.yml ./docker-compose-dev.yml 
	@cp ./prod/docker-compose.ml.search.yml ./docker-compose.yml
	@docker compose up -d
	@rm ./docker-compose.yml
	@mv ./docker-compose-dev.yml ./docker-compose.yml  



swagger:
	@docker-compose up swagger -d
	@docker-compose up nginx -d

lint:
	@golangci-lint run ./...


storage:
	@protoc --go_out=./internal/ --go_opt=paths=source_relative \
    --go-grpc_out=./internal/ --go-grpc_opt=paths=source_relative \
    -Iprotos protos/storage/storage.proto
	
auth:
	@protoc --go_out=./internal/ --go_opt=paths=source_relative \
    --go-grpc_out=./internal/ --go-grpc_opt=paths=source_relative \
    -Iprotos protos/auth/auth.proto
	
recognizer:
	@cp ./configs/recognizer/Dockerfile ./recognizer/Dockerfile
	
recognizer-server:
	@python3 -m grpc_tools.protoc -Irecognizer=protos \
  	--python_out=. --pyi_out=. --grpc_python_out=. \
  	./protos/recognizer.proto
	
recognizer-client:
	@protoc --go_out=./internal/recognizer/service --go_opt=paths=source_relative \
    --go-grpc_out=./internal/recognizer/service --go-grpc_opt=paths=source_relative \
    -Iprotos protos/recognizer.proto

