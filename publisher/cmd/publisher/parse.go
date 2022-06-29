package publisher

import (
	"encoding/json"
	"fmt"
	"os"
	"publisher/pkg/parser"

	"github.com/spf13/cobra"
)

// definition of the command that parses the content of a given
// file to verify the serialisation of a protobuf message. The
// actual parsing capability is delegated to the `parser` package.
var parseCmd = &cobra.Command{
	Use:   "parse",
	Short: "Parses a protobuf message to a file (optionally wrapped into a CloudEvent)",
	Args:  cobra.OnlyValidArgs,
	Run: func(cmd *cobra.Command, args []string) {

		// parse the content based on the parameters passed to
		// the command.
		var result map[string]interface{}
		var err error
		if isRaw {
			result, err = parser.ParseRaw(sourcePath, schemaURI, isDynamic)
		} else {
			result, err = parser.ParseCloudEvent(sourcePath, schemaURI, isDynamic)
		}
		if err != nil {
			fmt.Println("Error while parsing message:" + err.Error())
			os.Exit(1)
		}
		// dump to file or console based on whether an output path
		// has been specified.
		if len(targetPath) > 0 {

			err = writeToTarget(targetPath, result)
			if err != nil {
				fmt.Println("Error: " + err.Error())
				os.Exit(1)
			}

		} else {

			data, _ := json.Marshal(result)
			fmt.Println(string(data))
		}
	},
}

// writeToTarget marshals the given content to a JSON string and
// then writes it to the specified file.
func writeToTarget(targetPath string, content interface{}) error {

	bytes, err := json.Marshal(content)
	if err != nil {
		return err
	}
	var fp *os.File
	fp, err = os.OpenFile(targetPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	var written int
	written, err = fp.Write(bytes)
	if err != nil {
		return err
	}
	if written != len(bytes) {
		return fmt.Errorf("write to file (path: %s) errored, not all data has been saved", targetPath)
	}

	return nil
}

// init initialises the command with the required flags
// and adds it to the root command.
func init() {
	rootCmd.AddCommand(parseCmd)
	parseCmd.Flags().BoolVarP(&isDynamic, "dynamic", "d", true, "Uses dynamic type resolution to deserialise protobuf binary")
	parseCmd.Flags().BoolVarP(&isRaw, "raw", "r", false, "Determine whether to emit the message as a raw protobuf binary (default) or wrapped in a CloudEvent structure")
	parseCmd.Flags().StringVarP(&sourcePath, "source_path", "s", "", "Path to the file where to read the message or CloudEvent from")
	parseCmd.Flags().StringVarP(&targetPath, "target_path", "t", "", "Path to the file where to store the message (existing files will be overwritten)")
	parseCmd.Flags().StringVarP(&schemaURI, "schema_uri", "u", "", "URI of the protobuf file descriptor providing type information about the message payload")
	parseCmd.Flags().StringVarP(&messageType, "type", "m", "", "Simple name of the protobuf message to parse")
	parseCmd.MarkFlagRequired("source_path")
	parseCmd.MarkFlagRequired("schema_uri")
}
