package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Config struct {
	Template map[string]interface{} `json:"template"`
}

var (
	counter      int64 = 1
	counterMutex sync.Mutex
)

func LoadConfig(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(file, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func GeneratePayload(config *Config) ([]byte, error) {
	processed := make(map[string]interface{})

	for key, value := range config.Template {
		switch v := value.(type) {
		case string:
			parsed, err := parsePlaceholder(v)
			if err != nil {
				return nil, fmt.Errorf("field '%s': %v", key, err)
			}
			processed[key] = parsed
		default:
			processed[key] = v
		}
	}

	return json.MarshalIndent(processed, "", "  ")
}

func parsePlaceholder(s string) (interface{}, error) {
	switch {
	// Счетчик + строка
	case strings.HasPrefix(s, "COUNTER_"):
		counterMutex.Lock()
		defer counterMutex.Unlock()
		baseValue := strings.TrimPrefix(s, "COUNTER_")
		parsedValue, err := parsePlaceholder(baseValue)
		if err != nil {
			return nil, err
		}
		if strValue, ok := parsedValue.(string); ok {
			result := fmt.Sprintf("%s_%d", strValue, counter)
			counter++
			return result, nil
		}
		return nil, fmt.Errorf("COUNTER_ supports only strings, got %T", parsedValue)

	// Случайная комбинация строк
	case strings.HasPrefix(s, "RANDOM_COMBO"):
		re := regexp.MustCompile(`RANDOM_COMBO\(\(([^)]+)\),\(([^)]+)\)\)`)
		matches := re.FindStringSubmatch(s)
		if len(matches) == 3 {
			part1 := strings.Split(matches[1], ",")
			part2 := strings.Split(matches[2], ",")
			if len(part1) == 0 || len(part2) == 0 {
				return nil, fmt.Errorf("empty parts in RANDOM_COMBO")
			}
			return fmt.Sprintf("%s %s",
				part1[rand.Intn(len(part1))],
				part2[rand.Intn(len(part2))]), nil
		}
		return nil, fmt.Errorf("invalid RANDOM_COMBO format")

	// Случайное целое число
	case strings.HasPrefix(s, "RANDOM_INT"):
		re := regexp.MustCompile(`RANDOM_INT\((\d+),(\d+)\)`)
		matches := re.FindStringSubmatch(s)
		if len(matches) == 3 {
			min, _ := strconv.Atoi(matches[1])
			max, _ := strconv.Atoi(matches[2])
			return rand.Intn(max-min+1) + min, nil
		}
		return nil, fmt.Errorf("invalid RANDOM_INT format")

	// Случайная строка из списка
	case strings.HasPrefix(s, "RANDOM_STRING"):
		re := regexp.MustCompile(`RANDOM_STRING\(([^)]+)\)`)
		matches := re.FindStringSubmatch(s)
		if len(matches) == 2 {
			options := strings.Split(matches[1], ",")
			if len(options) == 0 {
				return nil, fmt.Errorf("empty RANDOM_STRING options")
			}
			return options[rand.Intn(len(options))], nil
		}
		return nil, fmt.Errorf("invalid RANDOM_STRING format")

	// Случайное булево значение
	case s == "RANDOM_BOOL":
		return rand.Intn(2) == 1, nil

	// Случайное дробное число
	case strings.HasPrefix(s, "RANDOM_FLOAT"):
		re := regexp.MustCompile(`RANDOM_FLOAT\(([\d.]+),([\d.]+)\)`)
		matches := re.FindStringSubmatch(s)
		if len(matches) == 3 {
			min, _ := strconv.ParseFloat(matches[1], 64)
			max, _ := strconv.ParseFloat(matches[2], 64)
			return min + rand.Float64()*(max-min), nil
		}
		return nil, fmt.Errorf("invalid RANDOM_FLOAT format")

	// Случайная дата в формате YYYY-MM-DD
	case strings.HasPrefix(s, "RANDOM_DATE"):
		re := regexp.MustCompile(`RANDOM_DATE\(([^,]+),([^)]+)\)`)
		matches := re.FindStringSubmatch(s)
		if len(matches) == 3 {
			start, err1 := time.Parse("2006-01-02", matches[1])
			end, err2 := time.Parse("2006-01-02", matches[2])
			if err1 != nil || err2 != nil {
				return nil, fmt.Errorf("invalid date format, use YYYY-MM-DD")
			}

			days := int(end.Sub(start).Hours() / 24)
			randomDays := rand.Intn(days + 1)
			return start.AddDate(0, 0, randomDays).Format("2006-01-02"), nil
		}
		return nil, fmt.Errorf("invalid RANDOM_DATE format")

	// Обычная строка (без обработки)
	default:
		return s, nil
	}
}
