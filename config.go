package main

type Config struct {
	DBUser     string
	DBPassword string
	DBName     string
	DBHost     string
	DBPort     string
	PgDumpPath string
	PsqlPath   string

	BackupDir  string
	MaxBackups int

	TelegramToken string
}

func LoadConfig() Config {
	return Config{
		DBUser:     "postgres",
		DBPassword: "8135",
		DBName:     "sharkbot",
		DBHost:     "localhost",
		DBPort:     "5432",
		PgDumpPath: "",
		PsqlPath:   "",

		BackupDir:  "backups",
		MaxBackups: 5,

		TelegramToken: "8450644464:AAEYs9rMC_UmWcV9t9bC6WNNJdVSxpJbEDY",
	}
}
