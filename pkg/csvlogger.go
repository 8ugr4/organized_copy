package pkg

import (
	"encoding/csv"
	"os"
	"sync"
)

// CSVLogger writes log entries into a CSV file with three columns:
// sourceFileName, destinationFileName, SUCCESS/FAILURE.
type CSVLogger struct {
	mu     sync.Mutex
	writer *csv.Writer
	file   *os.File
}

// NewCSVLogger creates or truncates a CSV file and writes the header row.
func NewCSVLogger(path string) (*CSVLogger, error) {
	if path == "" {
		return nil, nil
	}
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	w := csv.NewWriter(f)

	// header
	if err := w.Write([]string{"sourceFilePath", "destinationFilePath", "fileName", "status"}); err != nil {
		f.Close()
		return nil, err
	}
	w.Flush()
	return &CSVLogger{writer: w, file: f}, nil
}

// Log writes single entry into the CSV file.
func (l *CSVLogger) Log(status, source, fileName, destination string) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	record := []string{source, destination, fileName, status}
	if err := l.writer.Write(record); err != nil {
		return err
	}
	l.writer.Flush()
	return l.writer.Error()
}

func (l *CSVLogger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.writer.Flush()
	return l.file.Close()
}
