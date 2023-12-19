package main

import (
	"flag"
	"fmt"
	"math"
	"os"
)

const APPVERSION = "v1.7.0/2023"
const LINEBYTES = 16 // number of bytes to a line
const HELPFILENAME = "name of file to be dumped"
const HELPLINES = "number of 16-byte lines to dump (0 dumps the whole file)"
const HELPVERSIONSTRING = "print the current version"

func main() {

	fileName := flag.String("f", "", HELPFILENAME)
	lines := flag.Int("c", 0, HELPLINES)
	version := flag.Bool("v", false, HELPVERSIONSTRING)
	flag.Parse()

	if *version {
		fmt.Println(APPVERSION)
		os.Exit(0)
	}

	if *fileName == "" {
		usage()
		os.Exit(1)
	}

	if *lines == 0 {
		*lines = math.MaxInt64
	} else {
		if *lines < 0 {
			*lines = *lines * (-1)
		}
	}
	dumpFile(*fileName, *lines)
}

func dumpFile(fileName string, lines int) {
	f, err := os.Open(fileName)
	check(err)

	buff := make([]byte, LINEBYTES)

	for i := 0; i < lines; i++ {
		numread, err := f.Read(buff)
		if err == nil {

			offset := i * LINEBYTES
			dumpLine(offset, numread, buff)

			if numread < LINEBYTES {
				lines = 0
			}
		} else {
			break
		}
	}
	f.Close()
}

func dumpLine(offset int, numread int, data []byte) {
	line := fmt.Sprintf("%08X  ", offset)
	dataLine := "  |"
	for i := 0; i < LINEBYTES; i++ {
		if i < numread {
			line = fmt.Sprintf("%s %02X", line, data[i])
			if i == 7 {
				line = fmt.Sprintf("%s ", line)
			}
			dataLine = fmt.Sprintf("%s%c", dataLine, printableChar(data[i]))
		} else {
			line = fmt.Sprintf("%s   ", line)
			if i == 7 {
				line = fmt.Sprintf("%s ", line)
			}
			dataLine = fmt.Sprintf("%s ", dataLine)
		}
	}
	dataLine = fmt.Sprintf("%s|", dataLine)
	fmt.Println(fmt.Sprintf("%s%s", line, dataLine))
}

func printableChar(mychar byte) byte {
	if mychar > 31 && mychar < 127 {
		return mychar
	}
	//	if mychar > 127 && mychar < 255 {
	//		return mychar
	//	}
	return 46
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func usage() {
	fmt.Println("Usage: whexdump -f filename -c lines-to-dump")
}
