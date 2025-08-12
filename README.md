# Load Tester / Нагрузочный тестер
<div align="center">
  <img src="https://cdn-icons-png.freepik.com/512/18047/18047039.png?ga=GA1.1.1943782784.1755027723" alt="Load Testing" width="150">
  <br>
  <em>Инструмент для тестирования производительности API</em>
</div>

<br>

[![Go](https://img.shields.io/badge/Go-1.20+-blue.svg)](https://golang.org/)
[![HTTP](https://img.shields.io/badge/Protocol-HTTP-orange.svg)](https://developer.mozilla.org/en-US/docs/Web/HTTP)
[![JSON](https://img.shields.io/badge/Data-JSON-yellow.svg)](https://www.json.org/)

## 🇬🇧 English

Lightweight HTTP load tester for API endpoints. Generates configurable JSON payloads and tracks errors under high load.

**Features:**
- 🚀 Parallel request execution
- 📝 JSON templating with `RANDOM_INT`, `RANDOM_STRING` etc.
- 🔍 Detailed error logging (timeouts, 4xx/5xx)
- ⏱️ Configurable timeouts

**Quick Start:**
```bash
go run . -url=http://your-api.com -requests=1000 -workers=20
```


## 🇷🇺 Русский

Инструмент для нагрузочного тестирования API с генерацией JSON-запросов.

**Возможности:**
- 🚀 Параллельная отправка запросов
- 📝 Шаблоны JSON с генерацией случайных данных
- 🔍 Логирование ошибок (таймауты, 4xx/5xx)
- ⏱️ Настройка времени ожидания

**Быстрый старт:**
```bash
go run . -url=http://ваш-сервис.ру -requests=1000 -workers=20
```

## ⚙️ Configuration / Конфигурация

Edit `config.json`:
```json
{
  "template": {
    "id": "RANDOM_INT(1,1000)",
    "status": "RANDOM_STRING(new,pending,done)"
  }
}
```

## 🔴 Error Types / Типы ошибок

| Code/Код | Description/Описание          |
|----------|-------------------------------|
| 4XX      | Client-side issues / Ошибки клиента |
| 5XX      | Server errors / Ошибки сервера |
| Timeout  | Service overload / Перегрузка |


## 🤝 Contributing / Как помочь проекту

### 🇬🇧 English  
We welcome contributions! Here's how to help:  

1. **Fork** the repository  
2. Create a **feature branch** (`git checkout -b feature/your-idea`)  
3. Commit your changes (`git commit -am 'Add some feature'`)  
4. **Push** to the branch (`git push origin feature/your-idea`)  
5. Open a **Pull Request**  

Before submitting:  
- Run tests: `go test ./...`  
- Format code: `gofmt -s -w .`  

### 🇷🇺 Русский  
Приветствуем доработки! Как помочь:  

1. Сделайте **форк** репозитория  
2. Создайте **ветку** (`git checkout -b feature/ваша-фича`)  
3. Закоммитьте изменения (`git commit -am 'Добавил фичу'`)  
4. **Запушьте** ветку (`git push origin feature/ваша-фича`)  
5. Создайте **Pull Request**  

Перед отправкой:  
- Запустите тесты: `go test ./...`  
- Отформатируйте код: `gofmt -s -w .`  

---

### 🐛 Found a bug? / Нашли баг?  
Open an [Issue](https://github.com/VladislavKV-MSK/go-LoadTestHTTP/issues) with:  
/ Создайте [Issue](https://github.com/VladislavKV-MSK/go-LoadTestHTTP/issues) с:  
- Steps to reproduce / Шагами воспроизведения  
- Expected vs actual behavior / Ожидаемым и текущим поведением  
- Screenshots if applicable / Скриншотами (если есть)  
