package logger

import (
	"fmt"
	"testing"
)

type fakeLogWriter struct {
	output string
}

func (lw *fakeLogWriter) Write(bs []byte) (int, error) {
	lw.output = fmt.Sprintf("%s", string(bs))
	return fmt.Print(lw.output)
}

func TestInitLogger(t *testing.T) {
	InitLogger("TEST", "DEBUG")
	if logger == nil {
		t.Error("No logger initiated")
	}
	if logLevel != 1 {
		t.Error("Invalid logLevel initiated")
	}

}
func TestInitLoggerWithoutDebug(t *testing.T) {
	InitLogger("TEST", "ERROR")
	if logger == nil {
		t.Error("No logger initiated")
	}
	if logLevel != 0 {
		t.Error("Invalid logLevel initiated")
	}

}
func TestInitLoggerWithoutLogLevel(t *testing.T) {
	InitLogger("TEST", "")
	if logger == nil {
		t.Error("No logger initiated")
	}
	if logLevel != 0 {
		t.Error("Invalid logLevel initiated")
	}

}
func TestInitLoggerWithoutName(t *testing.T) {
	InitLogger("", "ERROR")
	writer := new(fakeLogWriter)
	logger.SetOutput(writer)
	expected := "[UNKNOWN Service]: INFO TEST RECORD\n"
	Info("TEST RECORD")
	if writer.output != expected {
		t.Error("Expected", expected, "got", writer.output)
	}

}
func TestDebug(t *testing.T) {
	InitLogger("TEST", "DEBUG")
	writer := new(fakeLogWriter)
	logger.SetOutput(writer)
	expected := "[TEST Service]: DEBUG TEST RECORD\n"
	Debug("TEST RECORD")
	if writer.output != expected {
		t.Error("Expected", expected, "got", writer.output)
	}
}
func TestDebugOff(t *testing.T) {
	InitLogger("TEST", "ERROR")
	writer := new(fakeLogWriter)
	logger.SetOutput(writer)
	expected := ""
	Debug("TEST RECORD")
	if writer.output != expected {
		t.Error("Expected", expected, "got", writer.output)
	}
}

func TestInfo(t *testing.T) {
	InitLogger("TEST", "DEBUG")
	writer := new(fakeLogWriter)
	logger.SetOutput(writer)
	expected := "[TEST Service]: INFO TEST RECORD\n"
	Info("TEST RECORD")
	if writer.output != expected {
		t.Error("Expected", expected, "got", writer.output)
	}
}

func TestError(t *testing.T) {
	InitLogger("TEST", "DEBUG")
	writer := new(fakeLogWriter)
	logger.SetOutput(writer)
	expected := "[TEST Service]: ERROR TEST RECORD\n"
	Error("TEST RECORD")
	if writer.output != expected {
		t.Error("Expected", expected, "got", writer.output)
	}
}
