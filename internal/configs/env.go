package configs

import (
	"log"
	"os"

	"github.com/alexedwards/argon2id"
	"github.com/joho/godotenv"
)

type env struct {
	PORT            string
	DBPort          string
	DBUser          string
	DBPassword      string
	DBName          string
	DBHost          string
	DBSchema        string
	JwtSecret       string
	ServerENV       string
	DefaultUser     string
	DefaultEmail    string
	DefaultPassword string
}

var Env *env

func SetupENV() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println(err.Error())
	}

	defaultPassword, _ := argon2id.CreateHash(os.Getenv("DEFAULT_PASSWORD"), argon2id.DefaultParams)

	Env = &env{
		PORT:            os.Getenv("PORT"),
		DBPort:          os.Getenv("DB_PORT"),
		DBUser:          os.Getenv("DB_USER"),
		DBPassword:      os.Getenv("DB_PASSWORD"),
		DBName:          os.Getenv("DB_NAME"),
		DBHost:          os.Getenv("DB_HOST"),
		DBSchema:        os.Getenv("DB_SCHEMA"),
		JwtSecret:       os.Getenv("JWT_SECRET"),
		ServerENV:       os.Getenv("SERVER_ENV"),
		DefaultUser:     os.Getenv("DEFAULT_USER"),
		DefaultEmail:    os.Getenv("DEFAULT_EMAIL"),
		DefaultPassword: defaultPassword,
	}
}
