package publisher

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

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
