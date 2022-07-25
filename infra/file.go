package infra

import (
	"context"
	"io"
	"os"
	"sync"
)

type FileReader struct {
	file     *os.File
	fileName string
	wg       *sync.WaitGroup
}

func NewFileReader(wg *sync.WaitGroup, fileName string) *FileReader {
	return &FileReader{
		wg:       wg,
		fileName: fileName,
	}
}

func (f *FileReader) Open() error {
	file, err := os.Open(f.fileName)

	if err != nil {
		return err
	}

	f.file = file

	return nil
}

func (f *FileReader) Read(ctx context.Context, lineCans []chan string) error {
	lineNumber := 0
	var lineString string

	iteration := 0
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			elem := make([]byte, 1)
			_, err := f.file.Read(elem)

			if err == io.EOF {
				return nil
			}

			if err != nil {
				return err
			}

			char := string(elem)
			if char == "\n" {
				lineNumber += 1

				if lineNumber > 1 {
					iteration += 1
					if iteration > (len(lineCans) - 1) {
						iteration = 0
					}

					f.wg.Add(1)
					lineCans[iteration] <- lineString
				}

				lineString = ""
			} else {
				lineString += char
			}
		}
	}
}

func (f *FileReader) Close() error {
	return f.file.Close()
}
