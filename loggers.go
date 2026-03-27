package main

import (
	"log"
	"os"
)

var Logger *log.Logger

func InitLogger() {
	file, _ := os.OpenFile("backup.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	Logger = log.New(file, "BACKUP: ", log.Ldate|log.Ltime|log.Lshortfile)
}
