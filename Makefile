
.PHONY: protos
protos:
	protoc --include_imports -Ischema --descriptor_set_out=build/root.pb schema/root.proto schema/imports/sub_message.proto --go_out=publisher


