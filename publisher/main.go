package main

import "publisher/cmd/publisher"

// main invokes the execution of the root command
// defined in the cmd/publisher package. This defines
// the command line tool for emitting events.
func main() {
	publisher.Execute()
}
