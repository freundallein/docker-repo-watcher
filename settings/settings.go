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
	LogLevel string
	Period,
	ImageAmount int
}

func (s *Settings) String() string {
	return fmt.Sprintf("\nREGISTRY_IP %s\nREGISTRY_PORT %s\nAPP_REFIX '%s'\nCRONTAB '%s'\nLOG_LEVEL %s\nPERIOD %d\nIMAGE_AMOUNT %d",
		s.RegistryIP, s.RegistryPort, s.AppPrefix, s.Crontab, s.LogLevel, s.Period, s.ImageAmount)
}

// NewSettings - create new settings
func NewSettings() *Settings {
	return &Settings{
		RegistryIP:   getEnvStr("REGISTRY_IP", DefaultRegistryIP),
		RegistryPort: getEnvStr("REGISTRY_PORT", DefaultRegistryPort),
		AppPrefix:    getEnvStr("APP_PREFIX", ""),
		Crontab:      getEnvStr("CRONTAB", DefaultCrontab),
		LogLevel:     getEnvStr("LOG_LEVEL", DefaultLogLevel),
		Period:       getEnvInt("PERIOD", DefaultPeriod),
		ImageAmount:  getEnvInt("IMAGE_AMOUNT", DefaultImageAmount),
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
