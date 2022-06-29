package parser

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"
)

// FullNameFormat enables the generation of the fullly qualified
// name of the protobuf message given its simple name.
const FullNameFormat = "hyp0th3rmi4.protobuf.sample.%s"

// ParseRaw reads the content of the file specified by `sourcePath`
// and interprets it as protobuf binary containing an instance of
// the message whose schema is defined in the location pointed by
// `schemaUri` (the fragment includes the type of the message). It
// then renders a map containing the values and attributes for this
// message.
func ParseRaw(sourcePath string, schemaUri string) (map[string]interface{}, error) {

	data, err := os.ReadFile(sourcePath)
	if err != nil {
		return nil, err
	}

	return deserialize(data, schemaUri)
}

// ParseCloudEvent reads the content of the file specified by `sourcePath`
// and interprets it as a JSON document containing the definition of a
// CloudEvent, whose payload is a base64 binary of a protobuf. It then
// deserialises the payload based on the type information contained in
// the event and the supplied `schemaUri` and converts it into a map,
// which is used to replace the original payload. The method returns
// the map representation of the entire CloudEvent, whose payload has
// been exploded into JSON.
func ParseCloudEvent(sourcePath string, schemaUri string) (map[string]interface{}, error) {

	ce := cloudevents.Event{}
	data, err := os.ReadFile(sourcePath)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &ce)
	if err != nil {
		return nil, err
	}

	structure, err := deserialize(ce.Data(), ce.DataSchema())
	if err != nil {
		return nil, err
	}
	container := map[string]interface{}{}
	json.Unmarshal(data, &container)
	container["datacontenttype"] = "application/json"
	container["data"] = structure

	return container, nil

}

// deserilize interprets the content of the given protobuf binary array according
// to the given message type specified by `schemaUri` (the fragment is the message
// type). The implementation of the method first constructs a file descriptor set
// from the given schema and uses it to setup a protobuf registry used to lookup
// the message descriptor mapped by the given type. It then constructs a dynamic
// message with the given `protobuf` array and the resolved descriptor to convert
// it then into a JSON document, returned as a map.
func deserialize(protobuf []byte, schemaUri string) (map[string]interface{}, error) {

	typeSchemaUrl, err := url.Parse(schemaUri)
	if err != nil {
		return nil, err
	}

	buffer, err := os.ReadFile(typeSchemaUrl.Path)
	if err != nil {
		return nil, err
	}
	fds := descriptorpb.FileDescriptorSet{}
	err = proto.Unmarshal(buffer, &fds)
	if err != nil {
		return nil, err
	}

	// we need to resolve the specific type thart
	// is serialised within the protobuf binary.
	files, err := protodesc.NewFiles(&fds)
	if err != nil {
		return nil, err
	}

	messageType := typeSchemaUrl.Fragment
	messageTypeFullName := protoreflect.FullName(fmt.Sprintf(FullNameFormat, messageType))

	descriptor, err := files.FindDescriptorByName(messageTypeFullName)
	if err != nil {
		return nil, err
	}

	msg := dynamicpb.NewMessage(descriptor.(protoreflect.MessageDescriptor))
	err = proto.Unmarshal(buffer, msg)
	if err != nil {
		return nil, err
	}
	data, err := protojson.Marshal(msg)
	if err != nil {
		return nil, err
	}
	structure := map[string]interface{}{}
	err = json.Unmarshal(data, &structure)
	if err != nil {
		return nil, err
	}
	return structure, nil

}
