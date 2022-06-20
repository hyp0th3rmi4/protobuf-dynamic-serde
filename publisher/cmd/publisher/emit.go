package publisher

import (
	"os"
	"publisher/pkg/publisher"

	"github.com/spf13/cobra"
)

var messageType string
var messagePath string

var emitCmd = &cobra.Command{
	Use:   "emit",
	Short: "Emits a cloudevent to a file",
	Run: func(cmd *cobra.Command, args []string) {

		switch messageType {
		case "SimpleMessage":
			publisher.SerializeSimpleMessage(messagePath)
		case "ImportMessage":
			publisher.SerializeImportedMessage(messagePath)
		case "ComplexMessage":
			publisher.SerializeComplexMessage(messagePath)
		case "ComposedMessage":
			publisher.SerializeComposedMessage(messagePath)
		default:
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(emitCmd)
	emitCmd.Flags().StringVarP(&messageType, "type", "t", "", "Type of the message to emit (SimpleMessage, ComplexMessage, ComposedMessage, ImportMessage)")
	emitCmd.Flags().StringVarP(&messagePath, "path", "p", "", "Path to the file where to store the message (existing files will be overwritten)")
	emitCmd.MarkFlagRequired("type")
	emitCmd.MarkFlagRequired("path")
}
