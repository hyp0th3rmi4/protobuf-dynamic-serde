package publisher

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"

	events "publisher/pkg/events/v1"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	proto "google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

// SerializeSimpleMessage persists to the specified path an `SimpleMessage`.
// The serialisation process can either wrap the serialised protobuf version
// of the message with a `CloudEvent` structure or publishing it as is (raw).
// In case the message is serialised within a cloud event, then the value of
// `schemaURI` is used to embed information about the schema of the message
// to enable dynamic consumption.
func SerializeSimpleMessage(path string, schemaURI string, isRaw bool) error {

	message := newSimpleMessage()
	return SerializeMessage(path, "SimpleMessage", schemaURI, message, isRaw)
}

// SerializeComplexMessage persists to the specified path an `ComplexMessage`.
// The serialisation process can either wrap the serialised protobuf version
// of the message with a `CloudEvent` structure or publishing it as is (raw).
// In case the message is serialised within a cloud event, then the value of
// `schemaURI` is used to embed information about the schema of the message
// to enable dynamic consumption.
func SerializeComplexMessage(path string, schemaURI string, isRaw bool) error {

	message := newComplexMessage()
	return SerializeMessage(path, "ComplexMessage", schemaURI, message, isRaw)
}

// SerializeImportMessage persists to the specified path an `ImportMessage`.
// The serialisation process can either wrap the serialised protobuf version of
// the message with a `CloudEvent` structure or publishing it as it is (raw).
// In case the message is serialised within a cloud event, then the value of
// `schemaURI` is used to embed information about the schema of the message
// to enable dynamic consumption.
func SerializeImportMessage(path string, schemaURI string, isRaw bool) error {

	message := newImportMessage()
	return SerializeMessage(path, "ImportMessage", schemaURI, message, isRaw)
}

// SerializeComposedMessage persists to the specified path an `ComposedMessage`.
// The serialisation process can either wrap the serialised protobuf version of
// the message with a `CloudEvent` structure or publishing it as it is (raw).
// In case the message is serialised within a cloud event, then the value of
// `schemaURI` is used to embed information about the schema of the message
// to enable dynamic consumption.
func SerializeComposedMessage(path string, schemaURI string, isRaw bool) error {

	message := newComposedMessage()
	return SerializeMessage(path, "ComposedMessage", schemaURI, message, isRaw)
}

// SerializeEnumMessage persists to the specified path an `EnumMessage`. The
// serialisation process can either wrap the serialised protobuf version of
// the message with a `CloudEvent`` structure or publishing it as it is (raw).
// In case the message is serialised within a cloud event, then the value of
// `schemaURI` is used to embed information about the schema of the message
// to enable dynamic consumption.
func SerializeEnumMessage(path string, schemaURI string, isRaw bool) error {
	message := newEnumMessage()
	return SerializeMessage(path, "EnumMessage", schemaURI, message, isRaw)
}

// SerializeNestedMessage persists to the specified path an `NestedMessage`.
// The serialisation process can either wrap the serialised protobuf version
// of the message with a `CloudEvent`` structure or publishing it as it is (raw).
// In case the message is serialised within a cloud event, then the value of
// `schemaURI` is used to embed information about the schema of the message
// to enable dynamic consumption.
func SerializeNestedMessage(path string, schemaURI string, isRaw bool) error {
	message := newNestedMessage()
	return SerializeMessage(path, "NestedMessage", schemaURI, message, isRaw)
}

// SerializeMessage implements the heavy-lifting required for emitting a cloud event.
// It generates a cloud even wrapper and configures it to transport the given message
// as payload of the event, serialised in base64 binary. The cloud event isntance is
// then serialised into JSON and persisted to the path specified by the caller. The
// value of schemaURI is used to provide a reference to the Protobuf file descriptor
// that can be used by consumer to deserialise the payload of the event.
func SerializeMessage(path string, messageType string, schemaURI string, message protoreflect.ProtoMessage, isRaw bool) error {

	buffer, err := proto.Marshal(message)
	if err != nil {
		return err
	}

	fp, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}

	if !isRaw {

		ce := cloudevents.NewEvent()
		ce.SetID(uuid.New().String())
		ce.SetSource("http://localhost/publisher")
		ce.SetSubject("publisher")
		ce.SetType(messageType)
		ce.SetDataSchema(fmt.Sprintf("%s#%s", schemaURI, messageType))
		ce.SetData("application/protobuf", buffer)
		ce.SetTime(time.Now())

		bytes, err := json.Marshal(ce)
		if err != nil {
			return err
		}

		buffer = bytes
	}

	written, err := fp.Write(buffer)

	if err != nil {
		return err
	}

	if written != len(buffer) {
		return errors.New("could not write the entire buffer to disk")
	}
	return nil
}

// newSimpleMessage generates a simple message and
// returns a pointer to it to the caller.
func newSimpleMessage() *events.SimpleMessage {

	return &events.SimpleMessage{
		Param_01: "first parameter",
		Param_02: true,
		Param_03: []byte{0x00, 0x01, 0x02},
		Param_04: -32,
		Param_05: -32321323412,
		Param_06: 10,
		Param_07: 2000000,
		Param_08: 12,
		Param_09: -391,
		Param_10: 88888,
		Param_11: 32412141431,
		Param_12: 33224,
		Param_13: -213123,
		Param_14: -0.2,
		Param_15: -0.000002,
	}
}

// newComplexMessage generates a complex message and
// returns a pointer to it to the caller.
func newComplexMessage() *events.ComplexMessage {

	return &events.ComplexMessage{
		Param_01: []string{"one", "two", "three"},
		Param_02: map[string]string{
			"autumn": "red",
			"winter": "blue",
			"spring": "green",
			"summer": "yellow",
		},
		Param_03: &events.ComplexMessage_Param_03String{
			Param_03String: "this is a oneof<string>",
		},
	}
}

// newImportMessage generates a import message
// and returns a pointer to it to the caller.
func newImportMessage() *events.ImportMessage {

	return &events.ImportMessage{
		Param_01: &timestamppb.Timestamp{},
		Param_02: &events.SubMessage{
			Param_01: events.Values_VALUE_1,
			Param_02: "this is nested!",
		},
	}
}

// newComposedMessage generates a composed message
// and returns a pointer to it to the caller.
func newComposedMessage() *events.ComposedMessage {

	return &events.ComposedMessage{
		Param_01: newSimpleMessage(),
		Param_02: newComplexMessage(),
	}
}

// newEnumMessage generates a message that wraps an
// enumeration and uses it to define an attribute.
func newEnumMessage() *events.EnumMessage {

	return &events.EnumMessage{
		PreferredSeason: events.EnumMessage_SPRING,
	}
}

// newNestedMessage generates a message that wraps another
// message definition and uses it to define an attribute
// of the message.
func newNestedMessage() *events.NestedMessage {

	return &events.NestedMessage{
		Users: []*events.NestedMessage_ProfileMessage{
			{
				Name: "Malcolm",
				Age:  32,
				Interests: []string{
					"woodworking",
					"photography",
					"cooking",
				},
			}, {
				Name: "Sally",
				Age:  29,
				Interests: []string{
					"photography",
					"basketball",
					"reading",
					"movies",
				},
			},
		},
	}
}
