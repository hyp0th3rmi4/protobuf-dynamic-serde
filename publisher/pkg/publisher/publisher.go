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

func SerializeSimpleMessage(path string) error {

	message := newSimpleMessage()
	return SerializeMessage(path, "SinmpleMessage", message)
}

func SerializeComplexMessage(path string) error {

	message := newComplexMessage()
	return SerializeMessage(path, "ComplexMessage", message)
}

func SerializeImportedMessage(path string) error {

	message := newImportMessage()
	return SerializeMessage(path, "ImportMessage", message)
}

func SerializeComposedMessage(path string) error {

	message := newComposedMessage()
	return SerializeMessage(path, "ComposedMessage", message)
}

func SerializeMessage(path string, messageType string, message protoreflect.ProtoMessage) error {

	buffer, err := proto.Marshal(message)
	if err != nil {
		return err
	}

	fp, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}

	ce := cloudevents.NewEvent()
	ce.SetID(uuid.New().String())
	ce.SetSource("http://localhost/publisher")
	ce.SetSubject("publisher")
	ce.SetType(messageType)
	ce.SetDataSchema(fmt.Sprintf("file://protobuf-dynamic-serde/build/root.pb#%s", messageType))
	ce.SetData("application/protobuf", buffer)
	ce.SetTime(time.Now())

	bytes, err := json.Marshal(ce)

	written, err := fp.Write(bytes)

	if err != nil {
		return err
	}

	if written != len(bytes) {
		return errors.New("Could not write the entire buffer to disk.")
	}
	return nil
}

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

func newImportMessage() *events.ImportMessage {

	return &events.ImportMessage{
		Param_01: &timestamppb.Timestamp{},
		Param_02: &events.SubMessage{
			Param_01: events.Values_VALUE_1,
			Param_02: "this is nested!",
		},
	}
}

func newComposedMessage() *events.ComposedMessage {

	return &events.ComposedMessage{
		Param_01: newSimpleMessage(),
		Param_02: newComplexMessage(),
	}
}
