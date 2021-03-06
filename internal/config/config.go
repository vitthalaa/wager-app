package config

import (
	"os"
	"strconv"
	"strings"
)

type AppConfig struct {
	Port           int
	DataBaseConfig DataBaseConfig
}

type DataBaseConfig struct {
	DBName        string
	DBHost        string
	DBPort        int
	DBUser        string
	DBPass        string
	DBMaxOpenConn int
	DBMaxIdleConn int
}

func GetAppConfig() AppConfig {
	return AppConfig{
		Port:           osValToInt("PORT", 8080),
		DataBaseConfig: GetDatabaseConfig(),
	}
}

func GetDatabaseConfig() DataBaseConfig {
	return DataBaseConfig{
		DBName:        osVal("POSTGRES_DB", ""),
		DBHost:        osVal("POSTGRES_HOST", "localhost"),
		DBPort:        osValToInt("POSTGRES_PORT", 5432),
		DBUser:        osVal("POSTGRES_USER", ""),
		DBPass:        osVal("POSTGRES_PASSWORD", ""),
		DBMaxOpenConn: osValToInt("POSTGRES_MAX_OPEN_CONN", 20),
		DBMaxIdleConn: osValToInt("POSTGRES_MAX_IDLE_CONN", 5),
	}
}

func osVal(key, defaultVal string) string {
	val, exist := os.LookupEnv(key)
	if !exist {
		return defaultVal
	}

	return strings.TrimSpace(val)
}

func osValToInt(key string, defaultVal int) int {
	val := osVal(key, "")
	if val == "" {
		return defaultVal
	}

	i64, err := strconv.ParseInt(val, 10, 32)
	if err != nil {
		return defaultVal
	}

	return int(i64)
}

func osValToBool(key string, defaultVal bool) bool {
	val := osVal(key, "")
	if val == "" {
		return defaultVal
	}

	return strings.ToUpper(val) == "TRUE"
}

func osValToArray(key, sep string, defaultVal []string) []string {
	val := osVal(key, "")
	if val == "" {
		return defaultVal
	}

	res := strings.Split(val, sep)
	for i := 0; i < len(res); i++ {
		res[i] = strings.TrimSpace(res[i])
	}

	return res
}
