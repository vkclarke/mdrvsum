package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"os"
)

var progName = func(path string) string {
	for i := len(path); i > 0; i-- {
		if path[i-1] == '/' || path[i-1] == '\\' {
			return path[i:]
		}
	}
	return path
}(os.Args[0])

func init() {
	log.SetFlags(0)
	log.SetPrefix(progName + ": ")
}

func main() {
	// load ROM data
	bin, out, err := func() ([]byte, *os.File, error) {
		var inName, outName string
		switch len(os.Args[1:]) {
		case 0:
			fmt.Println("USAGE:\n" +
				"\t" + progName + " infile [outfile]\n" +
				"NOTES:\n" +
				"\tOutput file is optional. Passing only an input file will overwrite it.")
			return nil, nil, errors.New("no arguments")
		case 1:
			inName = os.Args[1]
			outName = inName
		default:
			log.Println("ignoring arguments:", os.Args[3:])
			fallthrough
		case 2:
			inName = os.Args[1]
			outName = os.Args[2]
		}
		switch bin, err := os.ReadFile(inName); {
		case err != nil:
			return bin, nil, err
		case len(bin) < 512:
			return bin, nil, errors.New("ROM file must be 512 bytes or more")
		default:
			out, err := os.Create(outName)
			return bin, out, err
		}
	}()
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer out.Close()

	// calc/apply checksum
	checksum := func(bin []byte) (checksum uint16) {
		for i, j := 0, 2; j <= len(bin); i, j = i+2, j+2 {
			checksum += binary.BigEndian.Uint16(bin[i:j])
		}
		return
	}(bin[512:len(bin)])
	binary.BigEndian.PutUint16(bin[398:400], checksum)
	log.Printf("applied checksum: 0x%X", checksum)

	// write to file
	switch n, err := out.Write(bin); {
	case err != nil:
		log.Fatalln(err.Error())
	default:
		log.Printf("wrote %d bytes to %s", n, out.Name())
	}
}