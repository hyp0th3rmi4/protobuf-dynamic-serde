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


## Notes

- This is a __work in progress__ and not production code. 
- At present time, any attempt to read the protobuf binary, besides the use of protoc --decode-raw has been unsuccessful. :(
