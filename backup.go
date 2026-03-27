package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

func RunBackup(cfg Config) (string, error) {
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	file := fmt.Sprintf("%s/backup_%s.dump", cfg.BackupDir, timestamp)

	pgDumpBin, err := resolvePostgresTool(cfg.PgDumpPath, "pg_dump")
	if err != nil {
		return "", err
	}

	os.Setenv("PGPASSWORD", cfg.DBPassword)

	cmd := exec.Command(
		pgDumpBin,
		"-U", cfg.DBUser,
		"-h", cfg.DBHost,
		"-p", cfg.DBPort,
		"-F", "c",
		"-f", file,
		cfg.DBName,
	)

	err = cmd.Run()
	if err != nil {
		return "", err
	}

	return file, nil
}
