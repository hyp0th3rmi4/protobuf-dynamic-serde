# Dynamic Protobuf Deserialisation Example

This repository contains a simple producer consumer application to implement the concept of dynamic serialisation of Protobuf binaries into a corresponding JSON structure with information driven by a FileDescriptor instance for the Protobuf.

The application is composed by:
- a collection of `.proto` files containing the schema definition of the binary protobuf under test
- a producer application written in go that:
  - has a static linking to the protobuf message types and can generate both CloudEvent instances and raw protobuf messages
  - has the ability to dynamically parse the generated output without using the statically linked type to convert protobuf into JSON
- a consumer application written in Java that relies upon the file descriptor to convert the binary protobuf into a corresponding JSON document

The application uses an envelope for the protobuf binary that carries information about where to fetch the FileDescriptor for the content to deserialise along with an indication of the encoded root message. We use the CloudEvent specification to implement the envelop as it natively supports encoding in base64 of binary data and means to transport schema information, but the structure of the envelope does not need to necessarily be a CloudEvent.


## Run the Example


To try out the example (once it will be completed) do the following:

- install Golang with modules support
- install protobuf bindings for go (i.e. `go install google.golang.org/protobuf/cmd/protoc-gen-go@latest`)
- install Java 1.8
- install Maven (see [here](https://maven.apache.org/install.html))


Once your environment is fully set up do the following:

```bash
    
   git clone git@github.com/hyp0th3rmi4/protobuf-dynamic-serde.git 
   cd protobuf-dynamic-serde

   # builds the protobuf bindings and binaries
   make protos
   make publisher
   make consumer

   # creates a cloud event and saves it 
   # to: tmp/cloud-event.json
   #
   make publish-event
   
   # creates a raw protobuf binary file and
   # saves it: tmp/message.bin
   #
   make publish-raw
   
   # parses the cloud event previously created
   # and unpack the base64 binary protobuf into
   # JSON structure and saves the entire event
   # to: tmp/cloud-event-deserialised-go.json
   #
   make parse-event
   
   # parses the raw protobuf binaty file and 
   # converts it into a JSON structure that is
   # then saved into the file:
   # tmp/cloud-event-deserialised-go.json
   #
   make parse-raw

   # reads the previous cloud event and
   # deserialises its payload a corresponding
   # JSON structure and then saves it into
   # the file: tmp/cloud-event-deserialised.json
   #
   make consume-event
   
   # reads the previous protobuf binary and
   # deserialises it into the file:
   # tmp/cloud-event-deserialised.json
   #
   make consume-raw


  # clean all builds artefacts
  # and temporary files
  make clean

```

## Operating Modes

The sample allows you to perform different tasks in terms of serialisation/deserialisation:

- 💥 `make protos`: generating Go bindings for protobuf definitions of a set of test messages, with an associated file descriptor set 
- 💥 `make publish-event`: publishing a cloud event in Go, with a protobuf binary encoded as base64 string in the data attribute
- 💥 `make publish-raw`: publishing a byte array representing a serialised protobuf message in Go (i.e. `make publish-raw`)
- 💥 `make parse-event`: parsing a cloud event in Go, and converting the data payload from a base64 protobuf binary into a corresponding JSON representation by using a file descriptor for the type resolution 
- 💥 `make parse-raw`: parsing a byte array representing a serialised protobuf message in Go and converting it into a corresponding JSON representation by using a file descriptor for the type resolution
- 🚧 `make consume-event`: consuming a cloud event in Java, with a protobuf binary encoded as based64 string in the data attribute and converting it into the corresponding JSON structure by using a file descriptor for the type resolution
- 🚧 `make consume-raw`: consuming a byte array representing a serialised protobuf message in Java and converting it into the corresponding JSON structure by using a file descriptor for the type resolution

The parsing behaviour produces a JSON document where the following types are rendered as strings:

- `bytes`
- `int64` 
- `sint64`
- `fixed64`
- `sfixed64`

This seems to be a bug in the `protojson` library. An alternative would be creating a tree walker to produce the JSON by hand, and address these issues.

In addition, it is also possible to run the parsing behaviour by resolving the type descriptor from a static type representing the type serialised. This is accomplished by setting `--dynamic=false` and what this does is ignoring the file descriptor set, and resolving the type descriptor by mapping the type name encoded in the URL fragment of the schema to the corresponding statically linked type of the messages used for the purpose of testing. This is rather uninteresting, but primarily used for the purpose of testing during development.

## Notes

- This is a __work in progress__ and not production code. 
- The Java consumer uses the same philosophy implemented for parsing in Go, but uses a custom JSON converter that has incremental support:
  - Can parse message with simple types (primitive protobuf types)
  - Can parse repeated and map fields
  - Can parse composed messages (messages with non primitive types, but other defined messages)


