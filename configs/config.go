package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv                       string
	FrontendURL                  string
	AppPort                      string
	DBHost                       string
	DBPort                       string
	DBUser                       string
	DBPass                       string
	DBName                       string
	RedisAddr                    string
	RedisPassword                string
	JWTAccessSecret              string
	JWTRefreshSecret             string
	MailHost                     string
	MailPort                     string
	MailUser                     string
	MailPass                     string
	VerifyURL                    string
	CORSAllowedOrigins           string
	PayOSClientID                string
	PayOSAPIKey                  string
	PayOSChecksumKey             string
	PayOSReturnURL               string
	PayOSCancelURL               string
	PaymentUpdateStatusCancelURL string
}

// LoadConfig biến môi trường
func LoadConfig() *Config {
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		appEnv = "local"
	}

	if appEnv == "local" {
		_ = godotenv.Load(".env.local")
	}

	cfg := &Config{
		AppEnv:                       appEnv,
		FrontendURL:                  os.Getenv("FRONTEND_URL"),
		AppPort:                      os.Getenv("APP_PORT"),
		DBHost:                       os.Getenv("DB_HOST"),
		DBPort:                       os.Getenv("DB_PORT"),
		DBUser:                       os.Getenv("DB_USER"),
		DBPass:                       os.Getenv("DB_PASS"),
		DBName:                       os.Getenv("DB_NAME"),
		RedisAddr:                    os.Getenv("REDIS_ADDR"),
		RedisPassword:                os.Getenv("REDIS_PASSWORD"),
		JWTAccessSecret:              os.Getenv("JWT_ACCESS_SECRET"),
		JWTRefreshSecret:             os.Getenv("JWT_REFRESH_SECRET"),
		MailHost:                     os.Getenv("MAIL_HOST"),
		MailPort:                     os.Getenv("MAIL_PORT"),
		MailUser:                     os.Getenv("MAIL_USER"),
		MailPass:                     os.Getenv("MAIL_PASS"),
		VerifyURL:                    os.Getenv("VERIFY_URL"),
		CORSAllowedOrigins:           os.Getenv("CORS_ALLOWED_ORIGINS"),
		PayOSClientID:                os.Getenv("PAYOS_CLIENT_ID"),
		PayOSAPIKey:                  os.Getenv("PAYOS_API_KEY"),
		PayOSChecksumKey:             os.Getenv("PAYOS_CHECKSUM_KEY"),
		PayOSReturnURL:               os.Getenv("PAYOS_RETURN_URL"),
		PayOSCancelURL:               os.Getenv("PAYOS_CANCEL_URL"),
		PaymentUpdateStatusCancelURL: os.Getenv("PAYMENT_UPDATE_STATUS_CANCEL_URL"),
	}

	if cfg.AppPort == "" {
		log.Fatal("APP_PORT is required")
	}

	if cfg.VerifyURL == "" {
		log.Fatal("VERIFY_URL is required")
	}

	if cfg.PaymentUpdateStatusCancelURL == "" {
		log.Fatal("PAYMENT_UPDATE_STATUS_CANCEL_URL is required")
	}

	return cfg
}
