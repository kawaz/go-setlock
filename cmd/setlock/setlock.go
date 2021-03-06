package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"

	"github.com/moznion/go-setlock"
)

const (
	version = "1.2.0"
)

type opt struct {
	flagndelay     bool
	flagx          bool
	showVer        bool
	showVerVerbose bool
}

func main() {
	o := parseOpt()
	argv := flag.Args()
	run(o, argv...)
}

func run(o *opt, argv ...string) {
	if o.showVerVerbose {
		fmt.Printf("go-setlock (version: %s)\n", version)
		os.Exit(0)
	}

	if o.showVer {
		fmt.Printf("%s\n", version)
		os.Exit(0)
	}

	if len(argv) < 2 {
		// show usage
		fmt.Fprintf(os.Stderr, "setlock: usage: setlock [ -nNxXvV ] file program [ arg ... ]\n")
		os.Exit(100)
	}

	filePath := argv[0]

	locker := setlock.NewLocker(filePath, o.flagndelay)
	err := locker.LockWithErr()
	if err != nil {
		if o.flagx {
			os.Exit(0)
		}
		fmt.Println(err)
		os.Exit(111)
	}
	defer locker.Unlock()

	cmd := exec.Command(argv[1])
	for _, arg := range argv[2:] {
		cmd.Args = append(cmd.Args, arg)
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "setlock: fatal: unable to run %s: file does not exist\n", argv[1])
		os.Exit(111)
	}
}

func parseOpt() *opt {
	var n, N, x, X, showVer, showVerVerbose bool
	flag.BoolVar(&n, "n", false, "No delay. If fn is locked by another process, setlock gives up.")
	flag.BoolVar(&N, "N", false, "(Default.) Delay. If fn is locked by another process, setlock waits until it can obtain a new lock.")
	flag.BoolVar(&x, "x", false, "If fn cannot be opened (or created) or locked, setlock exits zero.")
	flag.BoolVar(&X, "X", false, "(Default.) If fn cannot be opened (or created) or locked, setlock prints an error message and exits nonzero.")
	flag.BoolVar(&showVer, "v", false, "Show version.")
	flag.BoolVar(&showVerVerbose, "V", false, "Show version verbosely.")
	flag.Parse()

	return &opt{
		flagndelay:     n && !N,
		flagx:          x && !X,
		showVer:        showVer,
		showVerVerbose: showVerVerbose,
	}
}
