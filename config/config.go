package config

import (
	"io"
	"os"
)

import (
	"github.com/pelletier/go-toml/v2"
)

var (
	_GLOBAL_CONFIG *Config
)

func GlobalConfig() *Config {
	return _GLOBAL_CONFIG
}

func InitConfig(config *Config) {
	_GLOBAL_CONFIG = config
}

type Config struct {
	Bind              string `toml:"bind"`
	HostProtocol      string `toml:"protocol"`
	SessionEncryptKey string `toml:"session_encrypt_key"`
	Sessionkey        string `toml:"session_key"`

	TemplatePath string `toml:"template_path"`
	StaticPath   string `toml:"static_path"`

	DbConfig struct {
		DbConnStr  string `toml:"db_conn_str"`
		DbMinCount int    `toml:"db_min_count"`
		DbMaxCount int    `toml:"db_max_count"`
	} `toml:"db"`

	HttpConfig struct {
		Enable bool   `toml:"enable"`
		Bind   string `toml:"bind"`
	} `toml:"http"`

	SiteConfig struct {
		DefaultSiteName string `toml:"default_site_name"`
	} `toml:"site"`

	DbObjConfig struct {
		DbConnStr  string `toml:"db_conn_str"`
		DbMinCount int    `toml:"db_min_count"`
		DbMaxCount int    `toml:"db_max_count"`
	} `toml:"db_obj"`
}

func Load(path string) *Config {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	buf, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}
	var config = new(Config)
	if err := toml.Unmarshal(buf, config); err != nil {
		panic(err)
	}

	return config
}
