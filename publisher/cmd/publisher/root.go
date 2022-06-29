package publisher

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// messageType stores the specified value for the type
// of message to emit.
var messageType string

// sourcePath points to a location storing the serialised
// data to be consumed.
var sourcePath string

// messagePath stores the specified value for the path
// where to save the file containing the event emitted.
var targetPath string

// schemaURI stores the specified value for the root
// URI that points to the file descriptor associated
// to the message payload.
var schemaURI string

// isRaw determines whether the message should be emitted
// as a raw binary protobuf or not.
var isRaw bool

// rootCmd is the root command for the publisher
// executable.
var rootCmd = &cobra.Command{

	Use:   "publisher",
	Short: "publisher - a simple cloud event publisher",
	Long: `publisher is a simple command line utility used to generate a cloud event and set its 
	data payload to a protobuf binary encoded in base 64 for the purpose of testing dynamic
	deserialisation.
	`,
	Args: cobra.OnlyValidArgs,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

// Execute executes the command.
func Execute() {

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}
