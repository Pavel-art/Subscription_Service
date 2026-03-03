# Subscription Service

Сервис для управления подписками пользователей с возможностью CRUD-операций и расчета суммарной стоимости подписок.

## 📌 О проекте

REST API сервис для агрегации данных о подписках пользователей с возможностью:
- Создания, чтения, обновления и удаления подписок
- Фильтрации подписок по пользователям и сервисам
- Расчета суммарной стоимости подписок за период

## 🛠 Технологии

- **Язык**: Go 1.24
- **Фреймворк**: Gin
- **База данных**: PostgreSQL
- **Документация**: Swagger
- **Логирование**: Zap
- **Конфигурация**: .env файлы
- **Запуск**: Docker Compose

## 📦 Установка и запуск

1. Клонируйте репозиторий:
```bash
git clone https://github.com/yourusername/subscription-service.git
cd subscription-service
```

2. Настройте окружение:
```bash
cp .env.example .env
# Отредактируйте .env файл под ваши настройки
```

3. Запустите сервис с помощью Docker:
```bash
docker-compose up -d
```

Сервис будет доступен по адресу: `http://localhost:8081`

## 📚 API Документация

После запуска сервиса документация Swagger будет доступна по адресу:
`http://localhost:8080/swagger/index.html`

## 🚀 Основные endpoints

- `POST /api/v1/subscriptions` - Создание подписки
- `GET /api/v1/subscriptions` - Получение списка подписок
- `GET /api/v1/subscriptions/:id` - Получение подписки по ID
- `PUT /api/v1/subscriptions/:id` - Обновление подписки
- `DELETE /api/v1/subscriptions/:id` - Удаление подписки
- `GET /api/v1/subscriptions/cost` - Расчет стоимости подписок

## ⚙️ Конфигурация

Основные настройки в `.env` файле:
- `PORT` - Порт сервиса
- `DB_URL` - URL подключения к PostgreSQL
- `LOG_LEVEL` - Уровень логирования

## 📜 Лицензия

MIT License
