# Load Tester / Нагрузочный тестер
<div align="center">
  <img src="https://cdn-icons-png.freepik.com/512/10941/10941323.png?ga=GA1.1.1943782784.1755027723" alt="Load Testing" width="600">
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

