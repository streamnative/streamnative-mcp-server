// Copyright (c) 2025 StreamNative, Inc.. All Rights Reserved.

package log

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

// IOLogger is a wrapper around io.Reader and io.Writer that can be used
// to log the data being read and written from the underlying streams
type IOLogger struct {
	reader io.Reader
	writer io.Writer
	logger *logrus.Logger
}

// NewIOLogger creates a new IOLogger instance
func NewIOLogger(r io.Reader, w io.Writer, logger *logrus.Logger) *IOLogger {
	return &IOLogger{
		reader: r,
		writer: w,
		logger: logger,
	}
}

// Read reads data from the underlying io.Reader and logs it.
func (l *IOLogger) Read(p []byte) (n int, err error) {
	if l.reader == nil {
		return 0, io.EOF
	}
	n, err = l.reader.Read(p)
	if n > 0 {
		l.logger.Infof("[stdin]: received %d bytes: %s", n, string(p[:n]))
	}
	return n, err
}

// Write writes data to the underlying io.Writer and logs it.
func (l *IOLogger) Write(p []byte) (n int, err error) {
	if l.writer == nil {
		return 0, io.ErrClosedPipe
	}
	l.logger.Infof("[stdout]: sending %d bytes: %s", len(p), string(p))
	return l.writer.Write(p)
}

// InitLogger initializes and returns a new logger
func InitLogger(outPath string) (*logrus.Logger, error) {
	if outPath == "" {
		return logrus.New(), nil
	}

	file, err := os.OpenFile(outPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	logger.SetOutput(file)

	return logger, nil
}
