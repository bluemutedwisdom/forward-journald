// Main package for forward-journald
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/coreos/go-systemd/journal"
)

var object_exe = flag.String("tag", "", "Add OBJECT_EXE=<tag> to journald entries")
var pri_info = flag.Bool("1", true, "forward stdin to journald as Priority Informational")
var pri_error = flag.Bool("2", false, "forward stdin to journald as Priority Error")

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %v:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %v -tag TAG\n\n", path.Base(os.Args[0]))
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	var priority journal.Priority
	var priority_name string
	if *pri_error {
		priority = journal.PriErr
		priority_name = "Priority Error"
	} else {
		priority = journal.PriInfo
		priority_name = "Priority Informational"
	}

	if journal.Enabled() {
		journal.Send(
			fmt.Sprintf("Forwarding stdin to journald using %v and tag %v", priority_name, *object_exe),
			priority,
			nil)
	} else {
		fmt.Fprintln(os.Stderr, "forward-journald: Unable to connect to journald")
		os.Exit(1)
	}

	var fields map[string]string = nil

	if len(*object_exe) > 0 {
		fields = map[string]string{
			"OBJECT_EXE": *object_exe,
		}
	}

	reader := bufio.NewReader(os.Stdin)

	var line string = ""
	var err error = nil
	for err == nil {
		line, err = reader.ReadString('\n')

		if journal.Enabled() {
			if err == nil {
				line = strings.TrimSpace(line)
				journal.Send(line, priority, fields)
			} else {
				// log any partial lines
				if len(line) > 0 {
					journal.Send(line, priority, fields)
				}
			}
		}
	}
}
