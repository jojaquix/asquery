package main

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
)

const (
	name = "asquery"
)

func main() {
	parser := flags.NewNamedParser(name, flags.Default)
	parser.AddCommand("query", "Execute a SQL query a asset.", "", &CmdQuery{})
	parser.AddCommand("shell", "Start an interactive session.", "", &CmdShell{})
	parser.AddCommand("version", "Show the version information.", "", &CmdVersion{})

	_, err := parser.Parse()
	if err != nil {
		if e, ok := err.(*flags.Error); ok && e.Type == flags.ErrCommandRequired {
			parser.WriteHelp(os.Stdout)
		}

		os.Exit(1)
	}
}

type cmd struct {
	Verbose bool `short:"v" description:"Activates the verbose mode"`
}

func (c *cmd) print(format string, a ...interface{}) {
	if !c.Verbose {
		return
	}

	fmt.Printf(format, a...)
}
