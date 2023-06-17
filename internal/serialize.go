package internal

import (
	"io"
	"os"
)

// Serialize gnark object to given file
func Serialize(gnarkObject io.WriterTo, fileName string) {
	f, err := os.Create(fileName)
	if err != nil {
		return
	}

	_, err = gnarkObject.WriteTo(f)
	if err != nil {
		return
	}
}

// Deserialize gnark object from given file
func Deserialize(gnarkObject io.ReaderFrom, fileName string) {
	f, err := os.Open(fileName)
	if err != nil {
		return
	}

	_, err = gnarkObject.ReadFrom(f)
	if err != nil {
		return
	}
}
