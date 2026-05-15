// Command brotherql is a CLI for printing to Brother QL-700 label printers.
package main

import (
	"errors"
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/kurtinge/brotherql"
)

const usage = `brotherql - print labels to Brother QL-700

Usage:
  brotherql <command> [options]

Commands:
  list      List connected printers
  status    Show printer status
  print     Print an image to a label
  help      Show this help

Run 'brotherql <command> --help' for command-specific options.
`

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(4)
	}

	switch os.Args[1] {
	case "list":
		os.Exit(cmdList(os.Args[2:]))
	case "status":
		os.Exit(cmdStatus(os.Args[2:]))
	case "print":
		os.Exit(cmdPrint(os.Args[2:]))
	case "help", "--help", "-h":
		fmt.Print(usage)
		os.Exit(0)
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n%s", os.Args[1], usage)
		os.Exit(4)
	}
}

func cmdList(args []string) int {
	fs := flag.NewFlagSet("list", flag.ExitOnError)
	_ = fs.Parse(args)

	infos, err := brotherql.List()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}
	if len(infos) == 0 {
		fmt.Fprintln(os.Stderr, "no printers found")
		return 2
	}
	for _, i := range infos {
		fmt.Println(i.String())
	}
	return 0
}

func cmdStatus(args []string) int {
	fs := flag.NewFlagSet("status", flag.ExitOnError)
	serial := fs.String("serial", "", "Specific printer serial number")
	fs.StringVar(serial, "s", "", "Specific printer serial number (shorthand)")
	_ = fs.Parse(args)

	p, err := openPrinter(*serial)
	if err != nil {
		return errExit(err)
	}
	defer func() { _ = p.Close() }()

	s, err := p.Status()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 3
	}

	if s.Ready {
		fmt.Println("Ready: yes")
	} else {
		fmt.Println("Ready: no")
	}
	if s.MediaLength == 0 {
		fmt.Printf("Media: %dmm continuous\n", s.MediaWidth)
	} else {
		fmt.Printf("Media: %dmm x %dmm\n", s.MediaWidth, s.MediaLength)
	}
	if s.Error != "" {
		fmt.Printf("Error: %s\n", s.Error)
	} else {
		fmt.Println("Error: none")
	}
	return 0
}

func cmdPrint(args []string) int {
	fs := flag.NewFlagSet("print", flag.ExitOnError)
	label := fs.String("label", "", "Label type (e.g. 62, 62x29)")
	fs.StringVar(label, "l", "", "Label type (shorthand)")
	serial := fs.String("serial", "", "Specific printer serial number")
	fs.StringVar(serial, "s", "", "Specific printer serial number (shorthand)")
	copies := fs.Int("copies", 1, "Number of copies")
	fs.IntVar(copies, "c", 1, "Number of copies (shorthand)")
	noCut := fs.Bool("no-cut", false, "Don't auto-cut after print")
	highDPI := fs.Bool("high-dpi", false, "Use 600 DPI instead of 300")
	_ = fs.Parse(args)

	if fs.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "usage: brotherql print --label <type> <image-file>")
		return 4
	}
	if *label == "" {
		fmt.Fprintln(os.Stderr, "--label is required")
		return 4
	}
	lbl, ok := labelByName(*label)
	if !ok {
		fmt.Fprintf(os.Stderr, "unknown label type: %s (try: 62, 62x29)\n", *label)
		return 4
	}

	f, err := os.Open(fs.Arg(0))
	if err != nil {
		fmt.Fprintf(os.Stderr, "open image: %v\n", err)
		return 1
	}
	defer func() { _ = f.Close() }()
	img, _, err := image.Decode(f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "decode image: %v\n", err)
		return 1
	}

	p, err := openPrinter(*serial)
	if err != nil {
		return errExit(err)
	}
	defer func() { _ = p.Close() }()

	opts := brotherql.PrintOptions{
		Label:   lbl,
		Copies:  *copies,
		AutoCut: !*noCut,
		HighDPI: *highDPI,
	}
	for i := 0; i < *copies; i++ {
		if err := p.Print(img, opts); err != nil {
			fmt.Fprintf(os.Stderr, "print: %v\n", err)
			return 3
		}
	}
	return 0
}

func openPrinter(serial string) (*brotherql.Printer, error) {
	if serial != "" {
		return brotherql.OpenBySerial(serial)
	}
	return brotherql.Open()
}

func labelByName(name string) (brotherql.LabelType, bool) {
	switch name {
	case "62":
		return brotherql.Label62, true
	case "62x29":
		return brotherql.Label62x29, true
	default:
		return brotherql.LabelType{}, false
	}
}

func errExit(err error) int {
	if errors.Is(err, brotherql.ErrPrinterNotFound) {
		fmt.Fprintln(os.Stderr, "no printer found")
		return 2
	}
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	return 1
}
