package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

// Env is 環境変数
type Env struct {
	AppEnv             string `envconfig:"APP_ENV" default:"dev"`
	DBHost             string `envconfig:"DB_HOST" default:"henpin_rdb"`
	DBPort             string `envconfig:"DB_PORT" default:"5432"`
	DBName             string `envconfig:"DB_NAME" default:"postgres"`
	DBUserName         string `envconfig:"DB_USERNAME" default:"postgres"`
	DBPassword         string `envconfig:"DB_PASSWORD" default:"P@ssw0rd"`
	SSLMode            string `envconfig:"SSL_MODE" default:"disable"`
	LogFormat          string `envconfig:"LOG_FORMAT" default:"json"`
	LogLevel           string `envconfig:"LOG_LEVEL" default:"debug"`
	FirebaseAPIKey     string `envconfig:"FIREBASE_API_Key"`
	FirebaseServiceKey string `envconfig:"FIREBASE_SERVICE_KEY"`
	CustomerBucket     string `envconfig:"CUSTOMER_BUCKET" default:"customer"`
	ManagerBucket      string `envconfig:"MANAGER_BUCKET" default:"manager"`
}

var env Env

func init() {
	if err := envconfig.Process("", &env); err != nil {
		panic(fmt.Errorf("failed to get environment variables: %w", err))
	}
}

func GetEnv() *Env {
	return &env
}

// IsTest : Check if the current environment is test
func IsTest() bool {
	return env.AppEnv == "test"
}

// IsDev : Check if the current environment is development
func IsDev() bool {
	return env.AppEnv == "dev"
}

// IsStg : Check if the current environment is staging
func IsStg() bool {
	return env.AppEnv == "stg"
}

// IsProd : Check if the current environment is production
func IsProd() bool {
	return env.AppEnv == "prod"
}
