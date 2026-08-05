package config

import (
	"log"
	"strconv"

	"github.com/spf13/viper"
)

var (
	config      *AppConfig
	JWT_SECRRET = ""
)

type AppConfig struct {
	DBUsername    string
	DBPassword    string
	DBHost        string
	DBPort        int
	DBName        string
	JWTKey        string
	KEYAPI        string
	KEYAPISecret  string
	CloudName     string
	RedisHost     string
	RedisPassword string
	ServiceEnv    string
	ServiceMode   string
	ServerPort    int // Tambahkan field ini
}

func InitConfig() *AppConfig {
	if config != nil {
		return config
	}

	config = ReadENV()
	log.Printf("✅ Config loaded: DBHost=%s, DBPort=%d, DBName=%s, ServiceEnv=%s",
		config.DBHost, config.DBPort, config.DBName, config.ServiceEnv)

	return config
}

func Get() *AppConfig {
	if config == nil {
		config = InitConfig()
	}
	return config
}

func ReadENV() *AppConfig {
	v := viper.New()

	// Prioritas 1: Environment Variables
	v.AutomaticEnv()

	// Determine environment
	mode := v.GetString("SERVICE_ENV")
	if mode == "" {
		mode = "LOCAL"
	}

	app := &AppConfig{
		ServiceEnv: mode,
	}

	// Load file config
	v.SetConfigName("local")
	v.SetConfigType("env")
	v.AddConfigPath(".")
	v.AddConfigPath("./config")

	if err := v.ReadInConfig(); err != nil {
		log.Printf("⚠️  Warning: Could not read local.env file: %v", err)
	}

	// === Mapping Config ===
	app.JWTKey = getString(v, "JWT_KEY", "")
	app.DBUsername = getString(v, "DBUSER", "root")
	app.DBPassword = getString(v, "DBPASS", "")
	app.DBHost = getString(v, "DBHOST", "localhost")
	app.DBName = getString(v, "DBNAME", "mintegra_local")

	// Parse DBPort
	if portStr := getString(v, "DBPORT", "3306"); portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil {
			app.DBPort = p
		} else {
			app.DBPort = 3306
		}
	}

	// Parse ServerPort - TAMBAHKAN INI
	if portStr := getString(v, "SERVER_PORT", "8080"); portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil {
			app.ServerPort = p
		} else {
			app.ServerPort = 8080
		}
	}

	app.KEYAPI = getString(v, "KEY_API", "")
	app.KEYAPISecret = getString(v, "KEY_API_SECRET", "")
	app.CloudName = getString(v, "CLOUD_NAME", "")
	app.RedisHost = getString(v, "IP_REDIS", "localhost")
	app.RedisPassword = getString(v, "PASS_REDIS", "")

	JWT_SECRRET = app.JWTKey

	return app
}

func getString(v *viper.Viper, key, defaultValue string) string {
	if val := v.GetString(key); val != "" {
		return val
	}
	return defaultValue
}
