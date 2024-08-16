package config

import (
  "fmt"
  "os"

  "github.com/joho/godotenv"
  "github.com/spf13/cast"
)

type Config struct {
  HTTPPort string

  PostgresHost     string
  PostgresPort     int
  PostgresUser     string
  PostgresPassword string
  PostgresDatabase string

  MongoHost     string
  MongoPort     int
  MongoUser     string
  MongoPassword string
  MongoDatabase string


  DefaultOffset string
  DefaultLimit  string

  TokenKey string
}

func Load() Config {
  if err := godotenv.Load(); err != nil {
    fmt.Println("No .env file found")
  }

  config := Config{}

  config.HTTPPort = cast.ToString(getOrReturnDefaultValue("HTTP_PORT", ":7777"))

  config.PostgresHost = cast.ToString(getOrReturnDefaultValue("POSTGRES_HOST", "localhost"))
  config.PostgresPort = cast.ToInt(getOrReturnDefaultValue("POSTGRES_PORT", 5432))
  config.PostgresUser = cast.ToString(getOrReturnDefaultValue("POSTGRES_USER", "postgres"))
  config.PostgresPassword = cast.ToString(getOrReturnDefaultValue("POSTGRES_PASSWORD", "20005"))
  config.PostgresDatabase = cast.ToString(getOrReturnDefaultValue("POSTGRES_DATABASE", "auth_exam"))

  config.MongoHost = cast.ToString(getOrReturnDefaultValue("MONGO_HOST", "mongo-db"))
  config.MongoPort = cast.ToInt(getOrReturnDefaultValue("MONGO_PORT", 27017))
  config.MongoPassword = cast.ToString(getOrReturnDefaultValue("MONGO_PASSWORD", "20005"))
  config.MongoDatabase = cast.ToString(getOrReturnDefaultValue("MONGO_DATABASE", "car_wash"))
  config.MongoUser = cast.ToString(getOrReturnDefaultValue("MONGO_USER", "postgres"))


  config.DefaultOffset = cast.ToString(getOrReturnDefaultValue("DEFAULT_OFFSET", "0"))
  config.DefaultLimit = cast.ToString(getOrReturnDefaultValue("DEFAULT_LIMIT", "10"))
  config.TokenKey = cast.ToString(getOrReturnDefaultValue("TokenKey", "my_secret_key"))
  return config
}

func getOrReturnDefaultValue(key string, defaultValue interface{}) interface{} {
  val, exists := os.LookupEnv(key)

  if exists {
    return val
  }

  return defaultValue
}
