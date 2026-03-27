package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func SendToTelegram(cfg Config, filePath string, ownerID string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	part, err := writer.CreateFormFile("document", file.Name())
	if err != nil {
		return err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}

	_ = writer.WriteField("chat_id", ownerID)
	writer.Close()

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendDocument", cfg.TelegramToken)

	req, err := http.NewRequest("POST", url, &body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func GetAdminOwnerIDs(cfg Config) ([]string, error) {
	psqlBin, err := resolvePostgresTool(cfg.PsqlPath, "psql")
	if err != nil {
		return nil, err
	}

	os.Setenv("PGPASSWORD", cfg.DBPassword)

	columnCmd := exec.Command(
		psqlBin,
		"-U", cfg.DBUser,
		"-h", cfg.DBHost,
		"-p", cfg.DBPort,
		"-d", cfg.DBName,
		"-tA",
		"-c", "SELECT column_name FROM information_schema.columns WHERE table_name = 'admins' AND column_name IN ('owners_id', 'user_id') ORDER BY CASE column_name WHEN 'owners_id' THEN 0 ELSE 1 END LIMIT 1;",
	)

	columnOut, err := columnCmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to detect admin id column: %v: %s", err, strings.TrimSpace(string(columnOut)))
	}

	idColumn := strings.TrimSpace(string(columnOut))
	if idColumn == "" {
		return nil, fmt.Errorf("admins table has neither owners_id nor user_id column")
	}

	cmd := exec.Command(
		psqlBin,
		"-U", cfg.DBUser,
		"-h", cfg.DBHost,
		"-p", cfg.DBPort,
		"-d", cfg.DBName,
		"-tA",
		"-c", "SELECT "+idColumn+" FROM admins;",
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to query admins: %v: %s", err, strings.TrimSpace(string(out)))
	}

	lines := strings.Split(string(out), "\n")
	ownerIDs := make([]string, 0, len(lines))
	for _, line := range lines {
		id := strings.TrimSpace(line)
		if id != "" {
			ownerIDs = append(ownerIDs, id)
		}
	}

	return ownerIDs, nil
}
