package core

import (
	"fmt"
	"os"
	"strconv"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*             - Config -              */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Our most used Configs are... -> 🌍 Globals!~
// -> Just call them like core.Whatever from anywhere and you're good to go.

var AppName = "obreros"
var AppAlias = "GOBR"
var AppEmoji = "👷"

// These are our non-global Configs 🌍❌
// -> The App loads an instance of this on startup and passes it around.
type Config struct {
	DBCfg      // -> DB Credentials and such
	LoggerCfg  // -> Logger settings
	RetrierCfg // -> N° Retries
}

func LoadConfig() *Config {

	// -> 🌍 Globals
	AppName = envVar("OBREROS_APP_NAME", AppName)
	AppAlias = envVar("OBREROS_APP_ALIAS", AppAlias)
	AppEmoji = envVar("OBREROS_APP_EMOJI", AppEmoji)

	return &Config{
		DBCfg:      loadDBConfig(),
		LoggerCfg:  loadLoggerConfig(),
		RetrierCfg: loadRetrierConfig(),
	}
}

func loadDBConfig() DBCfg {
	return DBCfg{
		Username:      envVar("DB_USERNAME", "root"),
		Password:      envVar("DB_PASSWORD", ""),
		Hostname:      envVar("DB_HOSTNAME", "localhost"),
		Port:          envVar("DB_PORT", "3306"),
		Schema:        envVar("DB_SCHEMA", "obreros"),
		Params:        envVar("DB_PARAMS", "?charset=utf8&parseTime=True&loc=Local"),
		Retries:       envVar("DB_RETRIES", 7),
		EraseAllData:  envVar("DB_ERASE_ALL_DATA", false),
		MigrateModels: envVar("DB_MIGRATE_MODELS", true),
		LogLevel:      LogLevels[envVar("DB_LOG_LEVEL", "error")],
	}
}

func loadLoggerConfig() LoggerCfg {
	return LoggerCfg{
		Level:       LogLevels[envVar("LOGGER_LEVEL", "info")],
		LevelStackT: LogLevels[envVar("LOGGER_LEVEL_STACKTRACE", "dpanic")],
		LogCaller:   envVar("LOGGER_LOG_CALLER", false),
	}
}

func loadRetrierConfig() RetrierCfg {
	return RetrierCfg{}
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

type (
	DBCfg struct {
		Username string
		Password string
		Hostname string
		Port     string
		Schema   string
		Params   string
		Retries  int

		EraseAllData  bool
		MigrateModels bool
		LogLevel      int
	}
	LoggerCfg struct {
		Level       int
		LevelStackT int
		LogCaller   bool
	}
	RetrierCfg struct {
	}
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func envVar[T string | bool | int](key string, fallback T) T {
	val, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	switch any(fallback).(type) {
	case string:
		return any(val).(T)
	case bool:
		return any(val == "true" || val == "TRUE" || val == "1").(T)
	case int:
		if intVal, err := strconv.Atoi(val); err == nil {
			return any(intVal).(T)
		}
	}
	return fallback
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func (c *DBCfg) GetSQLConnString() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s%s", c.Username, c.Password, c.Hostname, c.Port, c.Schema, c.Params)
}

// Used on init if the DB we need is not yet created
func (c *DBCfg) GetSQLConnStringNoSchema() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", c.Username, c.Password, c.Hostname, c.Port, c.Params)
}
