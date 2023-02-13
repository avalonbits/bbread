package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"math/bits"
	"os"
)

var basic = flag.String("basic", "", "path to bbc  basic tokenized file")

func main() {
	flag.Parse()

	if *basic == "" {
		log.Fatal("No file to read")
	}

	f, err := os.Open(*basic)
	if err != nil {
		log.Fatalf("unable to open file %q: %v", *basic, err)
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		log.Fatalf("error reading file %q: %v", *basic, err)
	}

	printListing(data)
}

func printListing(data []byte) {
	curr := 0
	for !bytes.Equal(data[curr:curr+3], []byte{0, 0xFF, 0xFF}) {
		length := int(data[curr])
		line := data[curr : curr+length]
		line = line[1 : len(line)-1]
		curr += length
		printLine(line)
	}
}

func printLine(line []byte) {
	num := (int(line[0]) | int(line[1])<<8)
	fmt.Printf("%5d ", num)

	for i := 2; i < len(line); i++ {
		b := line[i]
		if b == 0x20 {
			fmt.Print(" ")
			continue
		}
		if b == 0x8d {
			i += 3
			b0, b1, b2 := line[i-2], line[i-1], line[i]
			L := (bits.RotateLeft8(b0, 2) & 0xC0) ^ b1
			H := (bits.RotateLeft8(b0, 4) & 0xC0) ^ b2
			target := int(L) | int(H)<<8
			fmt.Printf("%d", target)
		} else {

			if b < 0x7F {
				fmt.Printf("%c", b)
			} else {
				tok := tokens[b]
				if tok == "" {
					fmt.Printf("[__0x%X__]", b)
				} else {
					fmt.Print(tok)
				}
			}
		}
	}
	fmt.Print("\n")
}

var tokens = map[byte]string{
	0x80: "AND",
	0x84: "OR",
	0x85: "ERROR",
	0x86: "LINE",
	0x87: "OFF",
	0x88: "STEP",
	0x8B: "ELSE",
	0x8C: "THEN",
	0x91: "TIME",
	0x9B: "COS",
	0xA3: "FALSE",
	0xA5: "GET",
	0xA8: "INT",
	0xAB: "LOG",
	0xAF: "PI",
	0xB3: "RND",
	0xB5: "SIN",
	0xB8: "TO",
	0xC3: "STR$",
	0xC6: "CMD",
	0xCE: "PUT",
	0xD4: "SOUND",
	0xD6: "CALL",
	0xDB: "CLS",
	0xDC: "DATA",
	0xDD: "DEF",
	0xDE: "DIM",
	0xDF: "DRAW",
	0xE0: "END",
	0xE1: "ENDPROC",
	0xE3: "FOR",
	0xE4: "GOSUB",
	0xE5: "GOTO",
	0xE6: "GCOL",
	0xE7: "IF",
	0xE8: "INPUT",
	0xE9: "LET",
	0xEB: "MODE",
	0xEA: "LOCAL",
	0xEC: "MOVE",
	0xED: "NEXT",
	0xEE: "ON",
	0xEF: "VDU",
	0xF0: "PLOT",
	0xF1: "PRINT",
	0xF2: "PROC",
	0xF3: "READ",
	0xF4: "REM",
	0xF5: "REPEAT",
	0xF7: "RESTORE",
	0xF8: "RETURN",
	0xFA: "STOP",
	0xFD: "UNTIL",
	0xFF: "OSCLI",
}
