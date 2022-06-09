/*
Copyright (c) 2014-2020 CGCL Labs
Container_Migrate is licensed under Mulan PSL v2.
You can use this software according to the terms and conditions of the Mulan PSL v2.
You may obtain a copy of Mulan PSL v2 at:
         http://license.coscl.org.cn/MulanPSL2
THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
See the Mulan PSL v2 for more details.
*/
package logrus

import (
	"bufio"
	"io"
	"runtime"
)

func (logger *Logger) Writer() *io.PipeWriter {
	return logger.WriterLevel(InfoLevel)
}

func (logger *Logger) WriterLevel(level Level) *io.PipeWriter {
	reader, writer := io.Pipe()

	var printFunc func(args ...interface{})
	switch level {
	case DebugLevel:
		printFunc = logger.Debug
	case InfoLevel:
		printFunc = logger.Info
	case WarnLevel:
		printFunc = logger.Warn
	case ErrorLevel:
		printFunc = logger.Error
	case FatalLevel:
		printFunc = logger.Fatal
	case PanicLevel:
		printFunc = logger.Panic
	default:
		printFunc = logger.Print
	}

	go logger.writerScanner(reader, printFunc)
	runtime.SetFinalizer(writer, writerFinalizer)

	return writer
}

func (logger *Logger) writerScanner(reader *io.PipeReader, printFunc func(args ...interface{})) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		printFunc(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		logger.Errorf("Error while reading from Writer: %s", err)
	}
	reader.Close()
}

func writerFinalizer(writer *io.PipeWriter) {
	writer.Close()
}
