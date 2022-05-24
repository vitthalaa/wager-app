//go:build integration
// +build integration

package integration_tests

import (
	"testing"

	env "github.com/joho/godotenv"
	"github.com/stretchr/testify/require"

	"github.com/vitthalaa/wager-app/internal/config"
)

func Test_EnvFile_Test(t *testing.T) {
	var loadEnv = env.Overload
	err := loadEnv("../.env")

	require.Nil(t, err)

	conf := config.GetAppConfig()

	require.NotEmpty(t, conf.DataBaseConfig.DBUser)
	require.NotEmpty(t, conf.DataBaseConfig.DBName)
	require.NotEmpty(t, conf.DataBaseConfig.DBHost)
	require.NotZero(t, conf.DataBaseConfig.DBPort)
}
