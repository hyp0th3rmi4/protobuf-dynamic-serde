MAKEFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
SCHEMA_URI := "file://$(dir $(MAKEFILE_PATH))build/root.pb"
MESSAGE_TYPE ?= SimpleMessage
MESSAGE_SCHEMA_URI := "file://$(dir $(MAKEFILE_PATH))build/root.pb\#$(MESSAGE_TYPE)"


.PHONY: protos
protos:	build/root.pb

build/root.pb:
	mkdir -p build
	protoc --include_imports -Ischema --descriptor_set_out=build/root.pb schema/root.proto schema/imports/sub_message.proto --go_out=publisher

.PHONY: publisher
publisher: build/publisher

build/publisher: build/root.pb
	cd publisher && go build && mv publisher ../build/

.PHONY: publish-raw
publish-raw: build/publisher
	mkdir -p tmp
	build/publisher emit --raw --schema_uri $(SCHEMA_URI) --target_path tmp/message.bin --type $(MESSAGE_TYPE)

.PHONY: publish-event
publish-event: build/publisher
	mkdir -p tmp
	build/publisher emit --schema_uri $(SCHEMA_URI) --target_path tmp/cloud-event.json --type $(MESSAGE_TYPE)

.PHONY: parse-event
parse-event: build/publisher
	build/publisher parse --schema_uri $(SCHEMA_URI) --source_path tmp/cloud-event.json --target_path tmp/cloud-event-deserialised-go.json

.PHONY: parse-raw
parse-raw: build/publisher
	build/publisher parse --raw --dynamic=true --schema_uri $(MESSAGE_SCHEMA_URI) --source_path tmp/message.bin --target_path tmp/cloud-event-deserialised-go.json

.PHONY: consumer
consumer: build/consumer.jar

build/consumer.jar: build/root.pb
	cd consumer && mvn package  && mv target/dynamic-serde-jar-with-dependencies.jar ../build/consumer.jar

.PHONY: consume-bin
consume-raw: build/consumer.jar
	java -jar build/consumer.jar --raw --schema_uri $(MESSAGE_SCHEMA_URI) --source_path tmp/message.bin --target_path tmp/message-deserialised.json

	
.PHONY: consume-event
consume-event: build/consumer.jar
	java -jar build/consumer.jar --source_path tmp/cloud-event.json --target_path tmp/cloud-event-deserialised.json

.PHONY: clean
clean:
	-rm -rf build
	-rm -rf tmp
	-rm -rf publisher/pkg/events
	cd consumer && mvn clean	
