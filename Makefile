MAKEFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
SCHEMA_URI := "file:///$(dir $(MAKEFILE_PATH))build/root.pb"
MESSAGE_TYPE ?= SimpleMessage


.PHONY: protos
protos:	build/root.pb

build/root.pb:
	mkdir -p build
	protoc --include_imports -Ischema --descriptor_set_out=build/root.pb schema/root.proto schema/imports/sub_message.proto --go_out=publisher

.PHONY: publisher
publisher: build/publisher

build/publisher: build/root.pb
	cd publisher && go build && mv publisher ../build/


.PHONY: publish-event
publish-event: build/publisher
	mkdir -p tmp
	build/publisher emit --schema_uri $(SCHEMA_URI) --path tmp/cloud-event.json --type $(MESSAGE_TYPE)

.PHONY: consumer
consumer: build/consumer.jar

build/consumer.jar: build/root.pb
	cd consumer && mvn package  && mv dynamic-serde-jar-with-dependencies.jar ../build/consumer.jar
	
.PHONY: consume-event
consume-event: build/consumer.jar
	java -jar build/consumer.jar --source_path tmp/cloud-event.json --target_path tmp/cloud-event-deserialised.json

.PHONY: clean
clean:
	-rm -rf build
	-rm -rf tmp
	-rm -rf publisher/pkg/events
	cd consumer && mvn clean	
