proto:
	@ docker run --rm -v ${PWD}:${PWD} -w ${PWD} rvolosatovs/protoc \
		--go_out=. \
		--proto_path=. \
		--go-grpc_out=. \
		--grpc-gateway_out=. \
		account.proto \
		events.proto

	@ cp -r ./pb/ ../account/ && cp -r ./pb/ ../logger/