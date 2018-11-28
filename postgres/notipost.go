package postgres

import (
	"database/sql"
	"github.com/spf13/viper"
)

type Notipost struct {
	NotificationService
}

type Config struct {
	URL string
}

func New(cfg Config) (*Notipost, error) {
	db, err := sql.Open("postgres", cfg.URL)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	np := &Notipost{}

	err = np.NotificationService.initializeSchema(db)
	if err != nil {
		return nil, err
	}

	return np, nil
}

func NewConfig(v *viper.Viper) *Config {
	cfg := &Config{
		URL: "postgres://localhost/notification?sslmode=verify-full",
	}
	v.SetDefault("postgres.url", cfg.URL)
	return cfg
}