

gen: clear
	protoc --proto_path=./proto proto/*.proto --go_out=. --go-grpc_out=. && mv ./grpc/pb Messages
	protoc --proto_path=./proto proto/*.proto --go_out=. --go-grpc_out=. && mv ./grpc/pb Users
	rm -rf grpc

clear: 
	rm -rf Messages/pb Users/pb
