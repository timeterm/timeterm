.PHONY: generate
generate:
	rm -rf go/**/*.pb.go 
	protoc --proto_path=src/timeterm_proto --go_out=go --go_opt=paths=source_relative src/**/**/*.proto
