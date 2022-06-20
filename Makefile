
.PHONY: protos
protos:	
	mkdir -p build
	protoc --include_imports -Ischema --descriptor_set_out=build/root.pb schema/root.proto schema/imports/sub_message.proto --go_out=publisher

publisher:
	go build publisher/*.go
