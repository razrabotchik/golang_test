package analyzer

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"os"
)

type TsvReader struct {
	fileName string
}

func NewTsvReader(fileName string) *TsvReader {
	return &TsvReader{
		fileName: fileName,
	}
}

func (f *TsvReader) Process(ctx context.Context, fn func(lineNumber uint, line string)) error {
	slog.InfoContext(ctx, "opening file", "fileName", f.fileName)

	file, err := os.Open(f.fileName)
	if err != nil {
		return fmt.Errorf("file open: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			slog.ErrorContext(ctx, "error closing file", "err", err, "fileName", f.fileName)
		}
	}()

	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 1024*1024), 1024*1024)

	var lineNumber uint
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		lineNumber++
		fn(lineNumber, scanner.Text())
	}

	return scanner.Err()
}
