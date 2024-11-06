package config

var conf RuntimeConfig

type RuntimeConfig struct {
	JWT    JWTConfig    `yaml:"jwt"`
	Data   DataConfig   `yaml:"data"`
	Outing OutingConfig `yaml:"outing"`
	Email  EmailConfig  `yaml:"email"`
}

func JWT() JWTConfig {
	return conf.JWT
}

type JWTConfig struct {
	AccessSecret  string `mapstructure:"access_secret"`
	RefreshSecret string `mapstructure:"refresh_secret"`
	AccessExp     int64  `mapstructure:"access_exp"`
	RefreshExp    int64  `mapstructure:"refresh_exp"`
}

func Data() DataConfig {
	return conf.Data
}

type DataConfig struct {
	Mysql MysqlConfig `yaml:"mysql"`
	Redis RedisConfig `yaml:"redis`
}

type MysqlConfig struct {
	User string `yaml:"user"`
	Pass string `yaml:"pass"`
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Db   string `yaml:"db"`
}

type RedisConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

func Outing() OutingConfig {
	return conf.Outing
}

type OutingConfig struct {
	OutingExp          int `mapstructure:"outing_exp"`
	OutingBlacklistExp int `mapstructure:"outing_blacklist_exp"`
}

func Email() EmailConfig {
	return conf.Email
}

type EmailConfig struct {
	Id   string `mapstructure:"id"`
	Pass string `mapstructure:"pass"`
}
