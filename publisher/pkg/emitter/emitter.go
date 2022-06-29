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

// SerializeSimpleMessage persists to the specified path a cloud event
// that contains as payload a random SimpleMessage encoded as a protobuf
// binary in base64. The cloud event generated contains attributes to
// enable consumers to retrieve the file descriptor whose URI is specified
// by the given `schemaURI`.
func SerializeSimpleMessage(path string, schemaURI string, isRaw bool) error {

	message := newSimpleMessage()
	return SerializeMessage(path, "SimpleMessage", schemaURI, message, isRaw)
}

// SerializeComplexMessage persists to the specified path a cloud event
// that contains as payload a random ComplexMessage encoded as a protobuf
// binary in base64. The cloud event generated contains attributes to
// enable consumers to retrieve the file descriptor whose URI is specified
// by the given `schemaURI`.
func SerializeComplexMessage(path string, schemaURI string, isRaw bool) error {

	message := newComplexMessage()
	return SerializeMessage(path, "ComplexMessage", schemaURI, message, isRaw)
}

// SerializeImportMessage persists to the specified path a cloud event
// that contains as payload a random ImportMessage encoded as a protobuf
// binary in base64. The cloud event generated contains attributes to
// enable consumers to retrieve the file descriptor whose URI is specified
// by the given `schemaURI`.
func SerializeImportMessage(path string, schemaURI string, isRaw bool) error {

	message := newImportMessage()
	return SerializeMessage(path, "ImportMessage", schemaURI, message, isRaw)
}

// SerializeComposedMessage persists to the specified path a cloud event
// that contains as payload a random ComposedMessage encoded as a protobuf
// binary in base64. The cloud event generated contains attributes to
// enable consumers to retrieve the file descriptor whose URI is specified
// by the given `schemaURI`.
func SerializeComposedMessage(path string, schemaURI string, isRaw bool) error {

	message := newComposedMessage()
	return SerializeMessage(path, "ComposedMessage", schemaURI, message, isRaw)
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
