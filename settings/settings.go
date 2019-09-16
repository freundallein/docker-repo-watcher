package settings

import (
	"fmt"
	"os"
	"strconv"
)

// Settings - service settings
type Settings struct {
	RegistryIP,
	RegistryPort,
	AppPrefix,
	Crontab,
	LogLevel,
	RegistryPath string
	Period,
	ImageAmount int
	AutoUpdate,
	CleanRegistry bool
}

// ToEnvString - form list of env variables
func (s *Settings) ToEnvString() []string {
	return []string{
		fmt.Sprintf("REGISTRY_IP=%s", s.RegistryIP),
		fmt.Sprintf("REGISTRY_PORT=%s", s.RegistryPort),
		fmt.Sprintf("APP_PREFIX=%s", s.AppPrefix),
		fmt.Sprintf("CRONTAB=%s", s.Crontab),
		fmt.Sprintf("LOG_LEVEL=%s", s.LogLevel),
		fmt.Sprintf("PERIOD=%d", s.Period),
		fmt.Sprintf("IMAGE_AMOUNT=%d", s.ImageAmount),
		fmt.Sprintf("AUTOUPDATE=%t", s.AutoUpdate),
		fmt.Sprintf("CLEAN_REGISTRY=%t", s.CleanRegistry),
		fmt.Sprintf("REGISTRY_PATH=%s", s.RegistryPath),
	}
}

func (s *Settings) String() string {
	return fmt.Sprintf("\nREGISTRY_IP=%s\nREGISTRY_PORT=%s\nAPP_PREFIX='%s'\nCRONTAB='%s'\nLOG_LEVEL=%s\nPERIOD=%d\nIMAGE_AMOUNT=%d\nAUTOUPDATE=%t",
		s.RegistryIP, s.RegistryPort, s.AppPrefix, s.Crontab, s.LogLevel, s.Period, s.ImageAmount, s.AutoUpdate)
}

// NewSettings - create new settings
func NewSettings() *Settings {
	return &Settings{
		RegistryIP:    getEnvStr("REGISTRY_IP", DefaultRegistryIP),
		RegistryPort:  getEnvStr("REGISTRY_PORT", DefaultRegistryPort),
		AppPrefix:     getEnvStr("APP_PREFIX", ""),
		Crontab:       getEnvStr("CRONTAB", DefaultCrontab),
		LogLevel:      getEnvStr("LOG_LEVEL", DefaultLogLevel),
		Period:        getEnvInt("PERIOD", DefaultPeriod),
		ImageAmount:   getEnvInt("IMAGE_AMOUNT", DefaultImageAmount),
		AutoUpdate:    getEnvBool("AUTOUPDATE", DefaultAutoUpdate),
		CleanRegistry: getEnvBool("CLEAN_REGISTRY", DefautlCleanRegistry),
		RegistryPath:  getEnvStr("REGISTRY_PATH", "/var/lib/registry"),
	}
}

func getEnvStr(name, fallback string) string {
	if value, ok := os.LookupEnv(name); ok {
		return value
	}
	return fallback
}

func getEnvInt(name string, fallback int) int {
	if valueStr, ok := os.LookupEnv(name); ok {
		value, err := strconv.Atoi(valueStr)
		if err == nil {
			return value
		}
	}
	return fallback
}

func getEnvBool(name string, fallback bool) bool {
	if value, ok := os.LookupEnv(name); ok {
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return fallback
		}
		return boolValue
	}
	return fallback
}
