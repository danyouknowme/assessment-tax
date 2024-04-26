package config

import "os"

type Config struct {
	Port          string
	DatabaseUrl   string
	AdminUsername string
	AdminPassword string
}

func New() *Config {
	return &Config{
		Port:          os.Getenv("PORT"),
		DatabaseUrl:   os.Getenv("DATABASE_URL"),
		AdminUsername: os.Getenv("ADMIN_USERNAME"),
		AdminPassword: os.Getenv("ADMIN_PASSWORD"),
	}
}
