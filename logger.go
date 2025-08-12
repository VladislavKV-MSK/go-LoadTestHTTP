package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	maxResponseSize = 1024 // 1Kb
)

func Init() {
	if _, err := os.Stat("errors.log"); os.IsNotExist(err) {
		if _, err := os.Create("errors.log"); err != nil {
			log.Fatal("Cannot create errors.log")
		}
	}
}

func LogError(workerID int, message string, request, response []byte,
	err error, status ...int) {

	entry := struct {
		Worker   int    `json:"worker_id"`
		Message  string `json:"message"`
		Error    string `json:"error,omitempty"`
		Status   int    `json:"status,omitempty"`
		Request  string `json:"request,omitempty"`
		Response string `json:"response,omitempty"`
	}{
		Worker:  workerID,
		Message: message,
		Request: string(request),
	}

	if err != nil {
		entry.Error = err.Error()
	}
	if len(status) > 0 {
		entry.Status = status[0]
	}
	if len(response) > 0 {
		entry.Response = string(response)
	}

	// Запись в файл errors.log
	logData, _ := json.Marshal(entry)
	appendToFile("errors.log", logData)

	// Дублирование в консоль
	log.Printf("ERROR: %s", logData)
}

func appendToFile(filename string, data []byte) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Failed to open log file: %v", err)
		return
	}
	defer file.Close()

	file.Write(append(data, '\n'))
}

func SendRequest(ctx context.Context, workerID int, url string, payload []byte) {
	// 1. Создаем запрос с контекстом
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(payload))
	if err != nil {
		LogError(workerID, "create request", payload, nil, err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// 2. Настраиваем клиент с таймаутами
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// 3. Отправка запроса
	resp, err := client.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			LogError(workerID, "request cancelled", payload, nil, ctx.Err())
		} else {
			LogError(workerID, "request failed", payload, nil, err)
		}
		return
	}
	defer resp.Body.Close()

	// 4. Чтение ответа (с ограничением размера)
	limitedReader := io.LimitReader(resp.Body, maxResponseSize)
	responseBody, _ := io.ReadAll(limitedReader)

	// 5. Логирование ошибок
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		LogError(workerID, "non-2xx response", payload, responseBody, nil, resp.StatusCode)
	}
}
