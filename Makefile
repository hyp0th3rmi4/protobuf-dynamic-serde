
.PHONY: protos
protos:	
	mkdir -p build
	protoc --include_imports -Ischema --descriptor_set_out=build/root.pb schema/root.proto schema/imports/sub_message.proto --go_out=publisher

publisher:
	cd publisher && go build && mv publisher ../build/

consumer:
	cd consumer && mvn package  && mv consumer ../build/
	
