package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/robfig/cron"
)

func main() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	config, err := LoadConfig(".", "config", "env")
	if err != nil {
		panic(err)
	}
	log.Println("config loaded successfuly")

	if config.RUN_ON_STARTUP || config.SINGLE_SHOT_MODE {
		runBackup(config)

		if config.SINGLE_SHOT_MODE {
			log.Println("Shutting down")

			os.Exit(0)
		}
	}

	c := cron.New()
	if err := c.AddFunc(config.BACKUP_CRON_SCHEDULE, func() { runBackup(config) }); err != nil {
		log.Fatal("error adding cron func: %w", err)
	}

	c.Start()

	<-quit

	c.Stop()

	log.Println("Shutting Down")
}

func runBackup(config Config) {
	log.Println("Starting Backup...")
	ctx := context.Background()
	backup := NewBackUpService(config)
	var filename string

	if config.BACKUP_FILE_PREFIX != "" {
		filename = config.BACKUP_FILE_PREFIX
	}

	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal("error getting working dir: %w", err)
	}

	timestamp := time.Now().Format("2006-01-02_15-04-05")

	filename = filename + "_" + timestamp + ".sql"
	filePath := filepath.Join(currentDir, filename)

	backupFilePath, err := backup.dumpToFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	if err := backup.uploadToR2(ctx, fmt.Sprintf("%s.gz", filename), backupFilePath); err != nil {
		log.Fatal(err)
	}

	backup.deleteFile(backupFilePath)

	log.Println("Backup Complete.")
}
