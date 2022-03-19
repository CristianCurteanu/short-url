
run-grpc:
	docker-compose build grpc
	docker-compose up grpc cache mongodb

kill-grpc:
	docker-compose down grpc

run-http:
	docker-compose build server
	docker-compose up server cache mongodb

run-all:
	docker-compose build && docker-compose up

run-grpc-standalone:
	go run cmd/grpc/main.go

run-http-standalone:
	go run cmd/res/main.go
