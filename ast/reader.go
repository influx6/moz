package ast

import (
	"io/ioutil"
	"os"
)

// readSource attempts to read giving sourcelines from the provided file.
func readSourceIn(path string, offset int64, length int) ([]byte, error) {
	inFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer inFile.Close()

	content := make([]byte, length)
	_, err = inFile.ReadAt(content, offset)
	return content, err
}

// readSource attempts to read giving sourcelines from the provided file.
func readSource(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}
