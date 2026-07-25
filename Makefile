.PHONY: proto
proto:
	protoc --go_out=. --go_opt=module=pulse \
	       --go-grpc_out=. --go-grpc_opt=module=pulse \
	       docs/proto/v1/*.proto