package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func Get() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error when load envi %s", err.Error())
	}

	return &Config{
		Server{
			Host: os.Getenv("SERVER_HOST"),
			Port: os.Getenv("SERVER_PORT"),
		},
		Database{
			Host:     os.Getenv("DATABASE_HOST"),
			Port:     os.Getenv("DATABASE_PORT"),
			User:     os.Getenv("DATABASE_USER"),
			Password: os.Getenv("DATABASE_PASSWORD"),
			Name:     os.Getenv("DATABASE_NAME"),
		},
		Email{
			Host:     os.Getenv("MAIL_HOST"),
			Port:     os.Getenv("MAIL_PORT"),
			User:     os.Getenv("MAIL_USER"),
			Password: os.Getenv("MAIL_PASS"),
		},
		Redis{
			Addr: os.Getenv("REDIS_ADDR"),
			Pass: os.Getenv("REDIS_PASS"),
		},
		Midtrans{
			Key:    os.Getenv("MIDTRANS_KEY"),
			IsProd: os.Getenv("MIDTRANS_ENV") == "production",
		},
	}
}
