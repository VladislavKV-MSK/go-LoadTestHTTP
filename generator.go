package main

import (
	"encoding/json"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Template map[string]interface{} `json:"template"`
}

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

var randSrc = rand.NewSource(time.Now().UnixNano())

func GeneratePayload(config *Config) ([]byte, error) {
	processed := make(map[string]interface{})

	for key, value := range config.Template {
		switch v := value.(type) {
		case string:
			processed[key] = parsePlaceholder(v)
		default:
			processed[key] = v // Оставляем как есть (числа, bool, null)
		}
	}

	return json.Marshal(processed)
}

func parsePlaceholder(s string) interface{} {
	switch {
	// Целые числа: RANDOM_INT(min,max)
	case strings.HasPrefix(s, "RANDOM_INT"):
		re := regexp.MustCompile(`RANDOM_INT\((\d+),(\d+)\)`)
		matches := re.FindStringSubmatch(s)
		if len(matches) == 3 {
			min, _ := strconv.Atoi(matches[1])
			max, _ := strconv.Atoi(matches[2])
			return rand.Intn(max-min+1) + min
		}
		return rand.Intn(1000) // Дефолтное значение

	// Строки: RANDOM_STRING(val1,val2,...)
	case strings.HasPrefix(s, "RANDOM_STRING"):
		re := regexp.MustCompile(`RANDOM_STRING\(([^)]+)\)`)
		matches := re.FindStringSubmatch(s)
		if len(matches) == 2 {
			options := strings.Split(matches[1], ",")
			return options[rand.Intn(len(options))]
		}
		return "default_" + strconv.Itoa(rand.Intn(100))

	// Булево: RANDOM_BOOL
	case s == "RANDOM_BOOL":
		return rand.Intn(2) == 1

	// Дробные числа: RANDOM_FLOAT(min,max)
	case strings.HasPrefix(s, "RANDOM_FLOAT"):
		re := regexp.MustCompile(`RANDOM_FLOAT\(([\d.]+),([\d.]+)\)`)
		matches := re.FindStringSubmatch(s)
		if len(matches) == 3 {
			min, _ := strconv.ParseFloat(matches[1], 64)
			max, _ := strconv.ParseFloat(matches[2], 64)
			return min + rand.Float64()*(max-min)
		}
		return rand.Float64() * 100

	default:
		return s // Если не плейсхолдер, возвращаем строку как есть
	}
}
