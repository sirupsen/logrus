package logrus

import (
	"bufio"
	"io"
	"runtime"
)

func (logger *Logger) InputWriter() (*io.PipeWriter) {
	inputReader, inputWriter := io.Pipe()

	go logger.inputWriterScanner(inputReader)
	runtime.SetFinalizer(inputWriter, inputWriterFinalizer)

	return inputWriter
}

func (logger *Logger) inputWriterScanner(inputReader *io.PipeReader) {
	scanner := bufio.NewScanner(inputReader)
	for scanner.Scan() {
		logger.Print(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		logger.Errorf("Error while reading from InputWriter: %s", err)
	}
	inputReader.Close()
}

func inputWriterFinalizer(inputWriter *io.PipeWriter) {
	inputWriter.Close()
}
