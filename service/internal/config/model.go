package config

import "go.uber.org/dig"

type NewsServer struct {
	dig.Out
	DBMS           DatabaseManageSystem `mapstructure:"DatabaseManageSystem"`
	ServiceAddress ServiceAddress       `mapstructure:"service_address"`
}

type ServiceAddress struct {
	News string `mapstructure:"News"`
}

type DatabaseManageSystem struct {
	MongoDBSystem  map[string]MongoDB  `mapstructure:"MongoDB"`
	RedisSystem    map[string]Redis    `mapstructure:"Redis"`
	MariaDBSystem  map[string]MariaDB  `mapstructure:"MariaDB"`
	PostgresSystem map[string]Postgres `mapstructure:"Postgres"`
}

type MongoDB struct {
	Host     string `mapstructure:"Host"`
	Port     string `mapstructure:"Port"`
	User     string `mapstructure:"User"`
	Password string `mapstructure:"Password"`
	Database string `mapstructure:"Database"`
}

type Redis struct {
	Host     string `mapstructure:"Host"`
	Port     string `mapstructure:"Port"`
	Password string `mapstructure:"Password"`
	Database int    `mapstructure:"Database"`
}

type MariaDB struct {
	Host     string `mapstructure:"Host"`
	Port     string `mapstructure:"Port"`
	User     string `mapstructure:"User"`
	Account  string `mapstructure:"Account"`
	Password string `mapstructure:"Password"`
	Database string `mapstructure:"Database"`
	MaxIdle  int    `mapstructure:"MaxIdle"`
	MaxOpen  int    `mapstructure:"MaxOpen"`
}

type Postgres struct {
	Host     string `mapstructure:"Host"`
	Port     string `mapstructure:"Port"`
	User     string `mapstructure:"User"`
	Account  string `mapstructure:"Account"`
	Password string `mapstructure:"Password"`
	Database string `mapstructure:"Database"`
	MaxIdle  int    `mapstructure:"MaxIdle"`
	MaxOpen  int    `mapstructure:"MaxOpen"`
}
