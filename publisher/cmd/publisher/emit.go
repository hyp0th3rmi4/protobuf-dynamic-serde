package publisher

import (
	"fmt"
	"os"
	emitter "publisher/pkg/emitter"

	"github.com/spf13/cobra"
)

// emitters is a map of callbacks that is used to map
// the specified type of message to emit to the
// corresponding function that will emit and persist
// the message according to the given parameters.
var emitters = map[string]func(string, string, bool) error{
	"SimpleMessage":   emitter.SerializeSimpleMessage,
	"ComplexMessage":  emitter.SerializeComplexMessage,
	"ComposedMessage": emitter.SerializeComposedMessage,
	"ImportMessage":   emitter.SerializeImportMessage,
	"EnumMessage":     emitter.SerializeEnumMessage,
}

// definition of the command that emits the cloud event
// based on the given parameters. The implementation of
// the event generation and persistence to file is then
// delegated to the `publisher` packages.
var emitCmd = &cobra.Command{
	Use:   "emit",
	Short: "Emits a protobuf message to a file (optionally wrapped into a CloudEvent)",
	Args:  cobra.OnlyValidArgs,
	Run: func(cmd *cobra.Command, args []string) {

		var err error = nil
		emitter, isPresent := emitters[messageType]
		if !isPresent {
			err = fmt.Errorf("unknown message type: '%s'", messageType)
		} else {
			err = emitter(targetPath, schemaURI, isRaw)
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
	emitCmd.Flags().BoolVarP(&isRaw, "raw", "r", false, "Determine whether to emit the message as a raw protobuf binary (default) or wrapped in a CloudEvent structure")
	emitCmd.Flags().StringVarP(&messageType, "type", "m", "", "Type of the message to emit (SimpleMessage, ComplexMessage, ComposedMessage, ImportMessage)")
	emitCmd.Flags().StringVarP(&targetPath, "target_path", "t", "", "Path to the file where to store the message (existing files will be overwritten)")
	emitCmd.Flags().StringVarP(&schemaURI, "schema_uri", "u", "", "URI of the protobuf file descriptor providing type information about the message payload")
	emitCmd.MarkFlagRequired("type")
	emitCmd.MarkFlagRequired("target_path")
	emitCmd.MarkFlagRequired("schema_uri")
}
