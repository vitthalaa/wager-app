package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDatabaseConfigConfig(t *testing.T) {
	// Reset to previous value
	dbName, _ := os.LookupEnv("POSTGRES_DB")
	dbuser, _ := os.LookupEnv("POSTGRES_USER")
	defer func() {
		os.Setenv("POSTGRES_DB", dbName)
		os.Setenv("POSTGRES_USER", dbuser)
	}()

	os.Setenv("POSTGRES_DB", "test_db_env_name")
	os.Setenv("POSTGRES_USER", "test_db_env_user")

	conf := GetDatabaseConfig()

	assert.Equal(t, "test_db_env_name", conf.DBName)
	assert.Equal(t, "test_db_env_user", conf.DBUser)
}
