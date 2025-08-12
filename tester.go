package main

import (
	"context"
	"log"
	"sync"
	"time"
)

func RunLoadTest(ctx context.Context, url string, totalRequests, workers int, configPath string) error {
	config, err := LoadConfig(configPath)
	if err != nil {
		return err
	}

	startTime := time.Now()
	var wg sync.WaitGroup
	requestsPerWorker := distributeRequests(totalRequests, workers)

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(workerID, requests int) {
			defer wg.Done()
			runWorker(ctx, workerID, url, requests, config)
		}(i, requestsPerWorker[i])
	}

	// Ожидаем завершения или отмены контекста
	select {
	case <-ctx.Done():
		log.Println("Test interrupted by context")
		wg.Wait() // Дожидаемся завершения текущих запросов
		return ctx.Err()
	default:
		wg.Wait()
		log.Printf("Test completed in %v", time.Since(startTime))
		return nil
	}
}

func runWorker(ctx context.Context, workerID int, url string, requests int, config *Config) {
	for i := 0; i < requests; i++ {
		select {
		case <-ctx.Done():
			return // Прерываем работу при отмене контекста
		default:
			payload, err := GeneratePayload(config)
			if err != nil {
				LogError(workerID, "generate payload", nil, nil, err)
				continue
			}
			SendRequest(ctx, workerID, url, payload)
		}
	}
}

// Распределяет запросы между worker'ами с учетом остатка
func distributeRequests(total, workers int) []int {
	base := total / workers
	remainder := total % workers

	distribution := make([]int, workers)
	for i := 0; i < workers; i++ {
		distribution[i] = base
		if i < remainder {
			distribution[i]++
		}
	}
	return distribution
}
