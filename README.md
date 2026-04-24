# Glavredusgo - Поиск постов ВКонтакте

Приложение для индексации и поиска постов в группах ВКонтакте.

## Запуск с Docker

### Предварительные требования

- Docker
- Docker Compose

### Быстрый старт

1. Скопируйте файл переменных окружения:

```bash
cp env.example .env
```

2. Отредактируйте `.env` файл, добавив ваш токен VK API:

```bash
VK_API_TOKEN=your_vk_api_token_here
```

Получить токен VK API можно [тут](https://vk.com/apps?act=manage).

3. Соберите и запустите контейнеры:

```bash
docker-compose up -d --build
```

4. Приложение будет доступно по адресу:
   - Web-интерфейс: http://localhost:8080
   - API: http://localhost:8080/api/search
   - Swagger UI: http://localhost:8080/swagger/

### Остановка

```bash
docker-compose down
```

### Просмотр логов

```bash
docker-compose logs -f app
```

## Структура директорий

```
.
├── data/                  # Директория для SQLite БД
├── history.bleve/         # Индекс Bleve для поиска
├── docker-compose.yml     # Конфигурация Docker Compose
├── Dockerfile             # Сборка Docker образа
└── ...
```

## API Endpoints

- `POST /api/search` - Поиск постов (JSON)
- `GET /` - Web-интерфейс поиска
- `GET /swagger/` - Swagger UI