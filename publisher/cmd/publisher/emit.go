package publisher

import "github.com/spf13/cobra"

var emitCmd = &cobra.Command{
	Use:   "emit",
	Short: "Emits a cloudevent to a file",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	rootCmd.AddCommand(emitCmd)
}
