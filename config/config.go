package config

import (
	"io/ioutil"
	"os"
)

import (
	"github.com/naoina/toml"
)

var (
	GLOBAL_CONFIG *Config
)

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

	FastcgiConfig struct {
		Enable bool   `toml:"enable"`
		Bind   string `toml:"bind"`
	} `toml:"fastcgi"`
}

func Load(path string) *Config {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	buf, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	var config = new(Config)
	if err := toml.Unmarshal(buf, config); err != nil {
		panic(err)
	}

	return config
}
