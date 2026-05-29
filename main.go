package main

import (
	"flag"
	"fmt"
	"io"
	"math"
	"os"
)

const APPVERSION = "v1.9.20/2026"
const LINEBYTES = 16 // number of bytes to a line
const HELPLINES = "number of 16-byte lines to dump (0 dumps the whole file)"
const HELPVERSIONSTRING = "print the current version"

const (
	colorReset     = "\033[0m"
	colorOffset    = "\033[33m" // yellow
	colorSep       = "\033[90m" // dark gray
	colorPrintable = "\033[96m" // bright cyan
	colorDot       = "\033[90m" // dark gray
)

func main() {

	lines := flag.Int("c", 0, HELPLINES)
	version := flag.Bool("v", false, HELPVERSIONSTRING)
	flag.Parse()

	if *version {
		fmt.Println(APPVERSION)
		os.Exit(0)
	}

	if *lines == 0 {
		*lines = math.MaxInt
	} else {
		if *lines < 0 {
			*lines = -*lines
		}
	}

	args := flag.Args()
	if len(args) > 0 {
		f, err := os.Open(args[0])
		check(err)
		defer f.Close()
		dump(f, *lines)
	} else {
		stat, err := os.Stdin.Stat()
		check(err)
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			dump(os.Stdin, *lines)
		} else {
			usage()
			os.Exit(1)
		}
	}
}

func isTerminal() bool {
	if _, set := os.LookupEnv("NO_COLOR"); set {
		return false
	}
	stat, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (stat.Mode() & os.ModeCharDevice) != 0
}

func dump(r io.Reader, lines int) {
	useColor := isTerminal()
	buff := make([]byte, LINEBYTES)

	for i := 0; i < lines; i++ {
		numread, err := io.ReadFull(r, buff)
		if numread > 0 {
			dumpLine(i*LINEBYTES, numread, buff, useColor)
		}
		if err != nil {
			break
		}
	}
}

func dumpLine(offset int, numread int, data []byte, useColor bool) {
	var line, dataLine string
	if useColor {
		line = fmt.Sprintf("%s%08X%s  ", colorOffset, offset, colorReset)
		dataLine = fmt.Sprintf("  %s|%s", colorSep, colorReset)
	} else {
		line = fmt.Sprintf("%08X  ", offset)
		dataLine = "  |"
	}

	for i := 0; i < LINEBYTES; i++ {
		if i < numread {
			line = fmt.Sprintf("%s %02X", line, data[i])
			if i == 7 {
				line = fmt.Sprintf("%s ", line)
			}
			ch := printableChar(data[i])
			if useColor {
				if isPrintable(data[i]) {
					dataLine = fmt.Sprintf("%s%s%c%s", dataLine, colorPrintable, ch, colorReset)
				} else {
					dataLine = fmt.Sprintf("%s%s%c%s", dataLine, colorDot, ch, colorReset)
				}
			} else {
				dataLine = fmt.Sprintf("%s%c", dataLine, ch)
			}
		} else {
			line = fmt.Sprintf("%s   ", line)
			if i == 7 {
				line = fmt.Sprintf("%s ", line)
			}
			dataLine = fmt.Sprintf("%s ", dataLine)
		}
	}

	if useColor {
		dataLine = fmt.Sprintf("%s%s|%s", dataLine, colorSep, colorReset)
	} else {
		dataLine = fmt.Sprintf("%s|", dataLine)
	}

	fmt.Printf("%s%s\n", line, dataLine)
}

func isPrintable(mychar byte) bool {
	return mychar > 31 && mychar < 127
}

func printableChar(mychar byte) byte {
	if isPrintable(mychar) {
		return mychar
	}
	return '.'
}

func check(e error) {
	if e != nil {
		fmt.Fprintf(os.Stderr, "whexdump: %v\n", e)
		os.Exit(1)
	}
}

func usage() {
	fmt.Println("Usage: whexdump [-c lines-to-dump] [filename]")
	fmt.Println("       echo data | whexdump [-c lines-to-dump]")
}
