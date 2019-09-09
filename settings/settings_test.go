package settings

import (
	"os"
	"reflect"
	"testing"
)

func TestNewSettings(t *testing.T) {
	observed := NewSettings()
	observedType := reflect.TypeOf(observed)
	expectedType := reflect.TypeOf(&Settings{})
	if observedType != expectedType {
		t.Error("Expected", expectedType, "got", observedType)
	}
	if observed.RegistryIP != DefaultRegistryIP {
		t.Error("Expected", DefaultRegistryIP, "got", observed.RegistryIP)
	}
	if observed.RegistryPort != DefaultRegistryPort {
		t.Error("Expected", DefaultRegistryPort, "got", observed.RegistryPort)
	}
	if observed.Crontab != DefaultCrontab {
		t.Error("Expected", DefaultCrontab, "got", observed.Crontab)
	}
	if observed.LogLevel != DefaultLogLevel {
		t.Error("Expected", DefaultLogLevel, "got", observed.LogLevel)
	}
}

func TestString(t *testing.T) {
	config := NewSettings()
	expected := "\nREGISTRY_IP=192.168.20.126\nREGISTRY_PORT=5000\nAPP_REFIX=''\nCRONTAB='0 0 0 * * *'\nLOG_LEVEL=ERROR\nPERIOD=60\nIMAGE_AMOUNT=5\nAUTOUPDATE=true"
	observed := config.String()
	if observed != expected {
		t.Error("Expected", expected, "got", observed)
	}
}
func TestToEnvString(t *testing.T) {
	config := NewSettings()
	expected := []string{
		"REGISTRY_IP=192.168.20.126",
		"REGISTRY_PORT=5000",
		"APP_REFIX=",
		"CRONTAB=0 0 0 * * *",
		"LOG_LEVEL=ERROR",
		"PERIOD=60",
		"IMAGE_AMOUNT=5",
		"AUTOUPDATE=true"}
	observed := config.ToEnvString()
	for i, value := range observed {
		if value != expected[i] {
			t.Error("Expected", expected[i], "got", value)
		}
	}

}

func TestGetEnvStr(t *testing.T) {
	expected := "Test"
	os.Setenv("TEST", expected)
	defer os.Unsetenv("TEST")
	observed := getEnvStr("TEST", "Default")
	if observed != expected {
		t.Error("Expected", expected, "got", observed)
	}
}

func TestGetEnvStrFallback(t *testing.T) {
	observed := getEnvStr("TEST", "Default")
	if observed != "Default" {
		t.Error("Expected Default, got", observed)
	}
}

func TestGetEnvInt(t *testing.T) {
	expected := "13"
	os.Setenv("TEST", expected)
	defer os.Unsetenv("TEST")
	observed := getEnvInt("TEST", 42)
	if observed != 13 {
		t.Error("Expected", expected, "got", observed)
	}
}

func TestGetEnvIntFallback(t *testing.T) {
	observed := getEnvInt("TEST", 42)
	if observed != 42 {
		t.Error("Expected 42, got", observed)
	}
}

func TestGetEnvIntInvalid(t *testing.T) {
	os.Setenv("TEST", "TEST")
	defer os.Unsetenv("TEST")
	observed := getEnvInt("TEST", 42)
	if observed != 42 {
		t.Error("Expected 42, got", observed)
	}
}

func TestGetEnvBool(t *testing.T) {
	expected := "1"
	os.Setenv("TEST", expected)
	defer os.Unsetenv("TEST")
	observed := getEnvBool("TEST", false)
	if observed != true {
		t.Error("Expected false, got", observed)
	}
}

func TestGetEnvBoolFallback(t *testing.T) {
	observed := getEnvBool("TEST", false)
	if observed != false {
		t.Error("Expected false, got", observed)
	}
}
