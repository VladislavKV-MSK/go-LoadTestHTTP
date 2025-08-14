package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"go.uber.org/zap"
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
				logger.Error("Failed to generate payload",
					zap.Int("worker", workerID),
					zap.Error(err))
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

func SendRequest(ctx context.Context, workerID int, url string, payload []byte) ([]byte, error) {
	startTime := time.Now()

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(payload))
	if err != nil {
		logger.Error("Failed to create request",
			zap.Int("worker", workerID),
			zap.String("url", url),
			zap.Error(err))
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			ResponseHeaderTimeout: 2 * time.Second,
			TLSHandshakeTimeout:   3 * time.Second,
			IdleConnTimeout:       30 * time.Second,
		},
	}

	resp, err := client.Do(req)
	duration := time.Since(startTime)

	// Всегда логируем основные параметры запроса
	logFields := []zap.Field{
		zap.Int("worker", workerID),
		zap.String("url", url),
		zap.Duration("duration", duration),
		zap.ByteString("payload", payload),
	}

	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()

	if err != nil {
		// Логируем ошибки сети/таймаута
		logger.Error("Request failed",
			append(logFields, zap.Error(err))...)
		return nil, err
	}

	// Чтение ответа с обработкой ошибок
	responseBody, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseSize))
	if err != nil {
		logger.Error("Failed to read response",
			append(logFields, zap.Error(err))...)
		return nil, err
	}

	// Добавляем статус и тело ответа в логи
	logFields = append(logFields,
		zap.Int("status", resp.StatusCode),
		zap.ByteString("response", responseBody),
	)

	// Логирование в зависимости от статуса
	switch {
	case resp.StatusCode >= 500:
		logger.Error("Server error", logFields...)
	case resp.StatusCode >= 400:
		logger.Warn("Client error", logFields...)
	case duration > 3*time.Second:
		logger.Warn("Slow response", logFields...)
	default:
		logger.Info("Request succeeded", logFields...)
	}

	if resp.StatusCode >= 400 {
		return responseBody, fmt.Errorf("HTTP %d: %s", resp.StatusCode, responseBody)
	}

	return responseBody, nil
}
