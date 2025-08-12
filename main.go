package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	// Инициализация логирования
	Init()

	// Инициализация контекста с возможностью прерывания
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Устанавливаем общий таймаут для всего теста
	ctx, cancel := context.WithTimeout(ctx, 30*time.Minute)
	defer cancel()

	// Флаги
	url := flag.String("url", "", "Target URL (required)")
	requests := flag.Int("requests", 1000, "Total requests")
	workers := flag.Int("workers", 10, "Concurrent workers")
	configPath := flag.String("config", "config.json", "Path to config")
	flag.Parse()

	if *url == "" {
		flag.PrintDefaults()
		log.Fatal("URL is required")
	}

	// Запуск с передачей контекста
	if err := RunLoadTest(ctx, *url, *requests, *workers, *configPath); err != nil {
		log.Fatal(err)
	}
}
