package main

import (
	"flag"
	"os"
)

func main() {
	var (
		dryrun  = flag.Bool("d", false, "dryrun flag")
		command = flag.String("c", "help", "command")
	)
	flag.Parse()

	if len(os.Args) == 1 {
		flag.Usage()
		return
	}

	switch {
	case "bulk" == *command:
		Bulk(*dryrun)
	case "endless" == *command:
		Endless()
	default:
		flag.Usage()
	}
}
