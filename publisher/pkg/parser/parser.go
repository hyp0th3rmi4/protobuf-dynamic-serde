package parser

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	events "publisher/pkg/events/v1"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"

	"publisher/pkg/logging"
)

// FullNameFormat enables the generation of the fullly qualified
// name of the protobuf message given its simple name.
const FullNameFormat = "hyp0th3rmi4.protobuf.sample.%s"

// ParseRaw reads the content of the file specified by `sourcePath`
// and interprets it as protobuf binary containing an instance of
// the message whose schema is defined in the location pointed by
// `schemaUri` (the fragment includes the type of the message). It
// then renders a map containing the values and attributes for this
// message. If `isDynamic` is `true` the resolution of the protobuf
// will be done by leveraging the type descriptor associated to the
// emssage specified in the schema URI, otherwise static types that
// are linked to the executable will be used based on the schema
// URI.
func ParseRaw(sourcePath string, schemaUri string, isDynamic bool) (map[string]interface{}, error) {

	data, err := os.ReadFile(sourcePath)
	if err != nil {
		return nil, err
	}

	logging.SugarLog.Infof("Read file (path: %s, size: %d bytes)", sourcePath, len(data))

	return deserialize(data, schemaUri, isDynamic)
}

// ParseCloudEvent reads the content of the file specified by `sourcePath` and
// interprets it as a JSON document containing the definition of a CloudEvent,
// whose payload is a base64 binary of a protobuf. It then deserialises the
// payload based on the type information contained in the event and the supplied
// `schemaUri` and converts it into a map, which is used to replace the original
// payload. The method returns the map representation of the entire CloudEvent,
// whose payload has been exploded into JSON. If `isDynamic` is `true` the
// resolution of the protobuf will be done by leveraging the type descriptor
// associated to the message specified in the schema URI, otherwise static types
// that are linked to the executable will be used based on the schema URI.
func ParseCloudEvent(sourcePath string, schemaUri string, isDynamic bool) (map[string]interface{}, error) {

	ce := cloudevents.Event{}
	data, err := os.ReadFile(sourcePath)
	if err != nil {
		return nil, err
	}

	logging.SugarLog.Infof("Read cloud event (path: %s, size: %d bytes)", sourcePath, len(data))

	err = json.Unmarshal(data, &ce)
	if err != nil {
		return nil, err
	}

	var structure map[string]interface{}
	logging.SugarLog.Infof("Unmarshalled file content into CloudEvent: %v", ce)

	structure, err = deserialize(ce.Data(), ce.DataSchema(), isDynamic)
	if err != nil {
		return nil, err
	}

	logging.SugarLog.Infof("Updated cloud event structure, with deserialised payload: %v", structure)

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
func deserialize(protobuf []byte, schemaUri string, isDynamic bool) (map[string]interface{}, error) {

	descriptor, err := resolveDescriptor(schemaUri, isDynamic)
	if err != nil {
		return nil, err
	}
	logging.SugarLog.Info("Resolved type descriptor for specified schema")

	msg := dynamicpb.NewMessage(descriptor)
	logging.SugarLog.Info("Created dynamic message container with descriptor")

	err = proto.Unmarshal(protobuf, msg)
	if err != nil {
		return nil, err
	}
	logging.SugarLog.Info("Unmarshalled protobuf binary into dynamic message")

	options := protojson.MarshalOptions{
		Multiline:     true,
		UseProtoNames: true,
		Indent:        "  ",
	}
	data, err := options.Marshal(msg)
	if err != nil {
		return nil, err
	}
	logging.SugarLog.Info("Marshalled dynamic message into JSON format")

	structure := map[string]interface{}{}
	err = json.Unmarshal(data, &structure)
	if err != nil {
		return nil, err
	}
	logging.SugarLog.Info("Unmarshalled JSON format into map[string]interface{}")

	return structure, nil

}

// resolveDescriptor examines the given schemaUri and extracts the
// necessary information to resolve the message descriptor pointed
// by the schema. If `isDynamic` is `true`, the a file descritptor
// set is created out of the file pointed by the schema URI and a
// type registry built out of it, which is then queried by using
// the fragment of the schema URI interpreted as type name. If the
// value of `isDynamic` is `false` only the fragment of the URI is
// extracted and mapped to the corresponding statically linked type
// in the executable, from which a message descriptor is resolved.
func resolveDescriptor(schemaUri string, isDynamic bool) (protoreflect.MessageDescriptor, error) {

	schemaUrl, err := url.Parse(schemaUri)
	if err != nil {
		return nil, err
	}

	var descriptor protoreflect.MessageDescriptor

	if isDynamic {

		logging.SugarLog.Infof("Using DYNAMIC type resolution, via type registry")

		var registry *protoregistry.Files
		registry, err = createRegistry(schemaUrl.Path)
		if err != nil {
			return nil, err
		}
		logging.SugarLog.Info("Resolved type registry")

		messageType := schemaUrl.Fragment
		messageTypeFullName := protoreflect.FullName(fmt.Sprintf(FullNameFormat, messageType))

		logging.SugarLog.Infof("Message full type name is: %s", messageTypeFullName)

		var pd protoreflect.Descriptor
		pd, err = registry.FindDescriptorByName(messageTypeFullName)
		if err != nil {
			return nil, err
		}
		logging.SugarLog.Info("Found descritptor")
		descriptor = pd.(protoreflect.MessageDescriptor)

	} else {

		logging.SugarLog.Info("Using STATIC type resolution, via compiled types descriptor")

		var msg protoreflect.ProtoMessage
		switch schemaUrl.Fragment {
		case "SimpleMessage":
			msg = &events.SimpleMessage{}
		case "ComplexMessage":
			msg = &events.ComplexMessage{}
		case "ImportMessage":
			msg = &events.ImportMessage{}
		case "ComposedMessage":
			msg = &events.ComposedMessage{}
		case "EnumtMessage":
			msg = &events.EnumMessage{}
		case "NestedMessage":
			msg = &events.NestedMessage{}
		default:
			return nil, fmt.Errorf("no message matching type: %s", schemaUrl.Fragment)
		}

		descriptor = msg.ProtoReflect().Descriptor()
	}

	return descriptor, nil

}

// createRegistry builds a registry of descriptor out of the protobuf
// file pointed by `pbFilePath`. The content of the file is unmarshalled
// as a `FileDescriptorSet`, which is then used to initialise the registry
// providing lookup capabilities for the descriptors in the set.
func createRegistry(pbFilePath string) (*protoregistry.Files, error) {

	buffer, err := os.ReadFile(pbFilePath)
	if err != nil {
		return nil, err
	}
	logging.SugarLog.Infof("Read file descriptor set metadata (size: %d bytes)", len(buffer))

	fds := descriptorpb.FileDescriptorSet{}
	err = proto.Unmarshal(buffer, &fds)
	if err != nil {
		return nil, err
	}
	logging.SugarLog.Info("Unmarshalled metadata infor file descriptor instance")

	return protodesc.NewFiles(&fds)
}
