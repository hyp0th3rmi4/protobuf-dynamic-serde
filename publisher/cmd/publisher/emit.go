package publisher

import (
	"fmt"
	"os"
	"publisher/pkg/publisher"

	"github.com/spf13/cobra"
)

// messageType stores the specified value for the type
// of message to emit.
var messageType string

// messagePath stores the specified value for the path
// where to save the file containing the event emitted.
var messagePath string

// schemaURI stores the specified value for the root
// URI that points to the file descriptor associated
// to the message payload.
var schemaURI string

// emitters is a map of callbacks that is used to map
// the specified type of message to emit to the
// corresponding function that will emit and persist
// the message according to the given parameters.
var emitters = map[string]func(string, string) error{
	"SimpleMessage":   publisher.SerializeSimpleMessage,
	"ComplexMessage":  publisher.SerializeComplexMessage,
	"ComposedMessage": publisher.SerializeComposedMessage,
	"ImportMessage":   publisher.SerializeImportMessage,
}

// definition of the command that emits the cloud event
// based on the given parameters. The implementation of
// the event generation and persistence to file is then
// delegated to the `publisher` packages.
var emitCmd = &cobra.Command{
	Use:   "emit",
	Short: "Emits a cloudevent to a file",
	Args:  cobra.OnlyValidArgs,
	Run: func(cmd *cobra.Command, args []string) {

		var err error = nil
		emitter, isPresent := emitters[messageType]
		if !isPresent {
			err = fmt.Errorf("unknown message type: '%s'", messageType)
		} else {
			err = emitter(messagePath, schemaURI)
		}

		if err != nil {
			os.Exit(1)
		}
	},
}

// init initialises the command with the required flags
// and adds it to the root command.
func init() {
	rootCmd.AddCommand(emitCmd)
	emitCmd.Flags().StringVarP(&messageType, "type", "t", "", "Type of the message to emit (SimpleMessage, ComplexMessage, ComposedMessage, ImportMessage)")
	emitCmd.Flags().StringVarP(&messagePath, "path", "p", "", "Path to the file where to store the message (existing files will be overwritten)")
	emitCmd.Flags().StringVarP(&schemaURI, "schema_uri", "s", "", "URI of the protobuf file descriptor providing type information about the message payload")
	emitCmd.MarkFlagRequired("type")
	emitCmd.MarkFlagRequired("path")
	emitCmd.MarkFlagRequired("schema_uri")
}
