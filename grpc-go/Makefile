gen:
	protoc --go_out=. --go_opt=paths=source_relative \
		   --go-grpc_out=. --go-grpc_opt=paths=source_relative */*.proto

s:
	@go run server/main.go $(ARGS)

c:
	@go run client/main.go $(ARGS)

nginx-up:
	docker run -it -d --rm \
		-v ${PWD}/nginx/nginx.conf:/etc/nginx/nginx.conf \
		-p 80:80 \
		--name nginx nginx

nginx-rm:
	docker stop nginx
	docker rm nginx --force

nginx-reload:
	docker exec nginx nginx -s reload

