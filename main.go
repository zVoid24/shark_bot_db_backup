package main

import (
	"os"
)

func main() {
	InitLogger()
	cfg := LoadConfig()

	os.MkdirAll(cfg.BackupDir, os.ModePerm)

	Logger.Println("Starting backup...")

	file, err := RunBackup(cfg)
	if err != nil {
		Logger.Println("Backup failed:", err)
		return
	}

	Logger.Println("Backup created:", file)

	compressed, err := Compress(file)
	if err != nil {
		Logger.Println("Compression failed:", err)
		return
	}

	os.Remove(file)

	Logger.Println("Compressed:", compressed)

	ownerIDs, err := GetAdminOwnerIDs(cfg)
	if err != nil {
		Logger.Println("Failed to load admin owners_id:", err)
		return
	}

	if len(ownerIDs) == 0 {
		Logger.Println("No admins found in admins table")
		return
	}

	for _, ownerID := range ownerIDs {
		err = SendToTelegram(cfg, compressed, ownerID)
		if err != nil {
			Logger.Printf("Telegram send failed for owners_id %s: %v", ownerID, err)
			continue
		}

		Logger.Printf("Sent to Telegram successfully for owners_id %s", ownerID)
	}

	CleanupOldBackups(cfg.BackupDir, cfg.MaxBackups)

	Logger.Println("Backup process completed")
}
