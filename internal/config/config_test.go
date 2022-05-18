package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDatabaseConfigConfig(t *testing.T) {
	// Reset to previous value
	dbName, _ := os.LookupEnv("DB_NAME")
	dbuser, _ := os.LookupEnv("DB_USER")
	defer func() {
		os.Setenv("DB_NAME", dbName)
		os.Setenv("DB_USER", dbuser)
	}()

	os.Setenv("DB_NAME", "test_db_env_name")
	os.Setenv("DB_USER", "test_db_env_user")

	conf := GetDatabaseConfig()

	assert.Equal(t, "test_db_env_name", conf.DBName)
	assert.Equal(t, "test_db_env_user", conf.DBUser)
}
