package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/joho/godotenv"
)

type ApplicationEnvEnum string

var AppConfig *AppVars

const (
	Development ApplicationEnvEnum = "DEVELOPMENT"
	Test        ApplicationEnvEnum = "TEST"
	Production  ApplicationEnvEnum = "PRODUCTION"
	Staging     ApplicationEnvEnum = "STAGING"
)

type AppVars struct {
	AppName       string `env:"APP_NAME"`
	AppEnv        string `env:"APP_ENV"`
	GinMode       string `env:"GIN_MODE"`
	Port          string `env:"PORT"`
	RootPath      string
	DB            DB
	RabbitMQ      RabbitMQConfig
	DebugMode     string
	SecretKey     string
	EncryptionKey string
	Issuer        string
}

type RabbitMQConfig struct {
	Username string
	Password string
	Host     string
	Port     int
}

type Parameters struct {
	Host            string
	Name            string
	Username        string
	Password        string
	Charset         string
	Collation       string
	ParseTime       bool
	MultiStatements bool
	Timeout         int
	MaxOpenConns    int
	MaxIdleConns    int
	MaxLifetime     time.Duration
}

type DB struct {
	Driver     string
	DBParams   Parameters
	DBRoParams Parameters
}

func RootPath() string {
	_, b, _, _ := runtime.Caller(0)
	root := filepath.Join(filepath.Dir(b), "../")
	return root
}

func Load() (*AppVars, error) {
	os.Setenv("TZ", "America/Sao_Paulo")

	path := RootPath()

	envFile := path + "/.env"

	err := godotenv.Load(envFile)
	if err != nil {
		fmt.Printf("Error loading .env file: %v", err)
		exception := errors.New(fmt.Sprint("Error loading env file: ", path+envFile))
		return nil, exception
	}

	// APP
	config := &AppVars{}
	config.AppName = GetEnvString("APP_NAME")
	config.AppEnv = strings.ToLower(GetEnvString("APP_ENV"))
	config.Port = viper.GetString("PORT")
	config.RootPath = path
	config.GinMode = GetEnvString("GIN_MODE")

	// KEYS

	config.SecretKey = GetEnvString("SECRET_KEY")
	config.Issuer = GetEnvString("ISSUER")
	config.EncryptionKey = GetEnvString("ENCRYPTION_KEY")

	// DB
	config.DB.Driver = "mysql"

	config.DB.DBParams.Host = GetEnvString("DB_HOST")
	config.DB.DBParams.Name = GetEnvString("DB_NAME")
	config.DB.DBParams.Username = GetEnvString("DB_USERNAME")
	config.DB.DBParams.Password = GetEnvString("DB_PASSWORD")
	config.DB.DBParams.Charset = "utf8mb4"
	config.DB.DBParams.Collation = "utf8mb4_unicode_ci"
	config.DB.DBParams.ParseTime = true
	config.DB.DBParams.MultiStatements = true
	config.DB.DBParams.MaxOpenConns = 5
	config.DB.DBParams.MaxIdleConns = 5
	config.DB.DBParams.MaxLifetime = time.Duration(60)

	config.RabbitMQ.Host = os.Getenv("RABBITMQ_HOST")
	config.RabbitMQ.Password = os.Getenv("RABBITMQ_PASSWORD")
	config.RabbitMQ.Port = mustAtoi(os.Getenv("RABBITMQ_PORT"))
	config.RabbitMQ.Username = os.Getenv("RABBITMQ_USERNAME")

	config.DebugMode = viper.GetString("DEBUG_MODE")

	AppConfig = config

	return config, config.Validate()
}

func (config AppVars) Validate() error {
	return validation.ValidateStruct(
		&config,
		validation.Field(&config.AppEnv, validation.Required),
		// validation.Field(&config.RoDBHost, validation.Required),
		// validation.Field(&config.RoDBName, validation.Required),
		// validation.Field(&config.RoDBUsername, validation.Required),
		// validation.Field(&config.RoDBPassword, validation.Required),
		// validation.Field(&config.RwDBHost, validation.Required),
		// validation.Field(&config.RwDBName, validation.Required),
		// validation.Field(&config.RwDBUsername, validation.Required),
		// validation.Field(&config.RwDBPassword, validation.Required),
		// validation.Field(&config.DBMaxOpenConns, validation.Required),
		// validation.Field(&config.DBMaxIdleConns, validation.Required),
		// validation.Field(&config.DBMaxDbLifetime, validation.Required),
		// validation.Field(&config.NewRelicLicense, validation.Required),
	)
}

func GetEnvString(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return ""
}

func GetEnvBool(key string) bool {
	if value, ok := os.LookupEnv(key); ok {
		if value, err := strconv.ParseBool(value); err == nil {
			return value
		}
	}
	return false
}

func mustAtoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("Error converting string to integer: %v", err)
	}
	return i
}
