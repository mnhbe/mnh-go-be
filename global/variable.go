package global

import (
  "fmt"
  "os"

  "github.com/spf13/viper"
)

const (
  Local       = "local"
  Development = "development"
  Sandbox     = "sandbox"
  Production  = "production"
)

type TodoAppConf struct{}

type Config struct {
  PostgresConnectionString string
  Mongo                    struct {
    ConnectionString string
    DatabaseName     string
  }
  JwtSecret string
}

var ServiceConfig *Config

func FetchEnvironmentVariables() {
  EnvType := os.Getenv("MNHBE_ENV")
  if EnvType == "" {
    EnvType = Local
  }
  ServiceConfig = NewConfig(EnvType)
}

func GetEnv() string {
  EnvType := os.Getenv("MNHBE_ENV")
  if EnvType == "" {
    EnvType = Local
  }

  return EnvType
}

func IsProductionEnv() bool {
  return GetEnv() == Production
}

func NewConfig(env string) *Config {
  cf := Config{}
  // Get the current working directory
  dir, err1 := os.Getwd()
  if err1 != nil {
    fmt.Println(err1)
  }

  // Print the current working directory
  fmt.Println("Current working directory:", dir)
  fileConfig := cf.GetConfigFile(env)
  fmt.Printf("Load Config File: %s \n", fileConfig)

  viper.SetConfigFile(fileConfig)
  err := viper.ReadInConfig()
  if err != nil {
    panic(err)
  }

  cf.PostgresConnectionString = viper.GetString("postgres_connection_string")
  cf.Mongo.ConnectionString = viper.GetString("mongo.connection_string")
  cf.Mongo.DatabaseName = viper.GetString("mongo.database_name")
  cf.JwtSecret = viper.GetString("jwt_secret")

  return &cf
}

func (config *Config) GetConfigFile(env string) string {
  fileF := "config/mnhbe_%s_config.json"
  switch env {
  case Local, Development, Sandbox, Production:
    return fmt.Sprintf(fileF, env)
  default:
    return fmt.Sprintf(fileF, Local)
  }
}
