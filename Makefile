up:
	docker-compose up -d
start:
	docker-compose start
stop:
	docker-compose stop
state:
	docker-compose ps
rebuild:
	docker-compose stop
	docker-compose pull
	docker-compose rm --force app
	docker-compose build --no-cache --pull
	docker-compose up -d --force-recreate
shell:
	docker-compose exec app /bin/bash

# Server is the one who defines protocols, so we put generated code in server pkg's
protobuf:
	cd proto && protoc --go_out=../pubsub-server/pkg/proto --go_opt=paths=source_relative \
    --go-grpc_out=../pubsub-server/pkg/proto --go-grpc_opt=paths=source_relative \
    --proto_path=. \
    pubsub.proto

test:
	go test -v -race -cover -count 10 ./pubsub-server/internal