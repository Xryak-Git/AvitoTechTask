package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"strings"
	"time"
)

type Config struct {
	Env        string `yaml:"env"  envDefault:"development"`
	LogLevel   string `yaml:"log_level" env-default:"debug"`
	HTTPServer `yaml:"http_server"`
	Postgres   `yaml:"postgres"`
}

type HTTPServer struct {
	Adress      string        `yaml:"port" envDefault:":8080"`
	Timeout     time.Duration `yaml:"timeout"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
}

type Postgres struct {
	Url     string `env:"POSTGRES_CONN"`
	JdbcUrl string `env:"POSTGRES_JDBC_URL"`

	User     string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASSWORD"`
	Host     string `env:"POSTGRES_HOST"`
	Port     int    `env:"POSTGRES_PORT"`
	DbName   string `env:"POSTGRES_DB"`
}

func jdbcToPgx(jdbcURL, username, password string) string {
	// Удаляем префикс "jdbc:"
	pgxURL := strings.Replace(jdbcURL, "jdbc:", "", 1)

	// Добавляем информацию о пользователе и пароле
	return fmt.Sprintf("postgres://%s:%s@%s", username, password, pgxURL[len("postgresql://"):])
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatalf("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file was not found by this path %s: ", configPath)
	}

	cfg := &Config{}

	if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	err := cleanenv.UpdateEnv(cfg)
	if err != nil {
		log.Fatalf("cannot update config: %s", err)
	}

	if cfg.Url == "" {
		if cfg.JdbcUrl != "" {
			cfg.Url = jdbcToPgx(cfg.Url, cfg.User, cfg.Password)
		}
		if cfg.User == "" || cfg.Password == "" || cfg.Host == "" || cfg.Port == 0 || cfg.DbName == "" {
			log.Fatalf("no data for for postgres url")
		}

		cfg.Url = fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DbName)
	}
	return cfg

}
