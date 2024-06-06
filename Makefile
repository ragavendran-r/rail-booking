# Makefile


protoc:
	cd proto && protoc --go_out=../protogen --go_opt=paths=source_relative \
	--go-grpc_out=../protogen --go-grpc_opt=paths=source_relative user/*.proto booking/*.proto seats/seats.proto