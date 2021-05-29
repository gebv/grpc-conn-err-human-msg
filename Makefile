
export CAPATH=./ca

init:
	mkdir -p ./ssl
	mkdir -p ./ca

gen-ssl: init
	CANAME=ok ./ssl-self-signed.sh c

	CANAME=ok OUTPATH=./ssl ./ssl-self-signed.sh d localhost-ok localhost
	CANAME=ok OUTPATH=./ssl ./ssl-self-signed.sh d localhost-expired localhost 0
	CANAME=ok OUTPATH=./ssl ./ssl-self-signed.sh s localhost-simple localhost

proto-gen-install:
	GO111MODULE=on go get google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
	GO111MODULE=on go get google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1

gen-proto: proto-gen-install
	protoc \
		--go_out=. \
    	--go-grpc_out=. \
		simple.proto

run:
	docker-compose up --build -d app
	docker-compose up -d sslok sslexpired sslsimple

test:
	go test -v ./...
