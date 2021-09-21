# Telegram бот для Redmine

[![Go Report Card](https://goreportcard.com/badge/github.com/alphatroya/redmine-helper-bot)](https://goreportcard.com/report/github.com/alphatroya/redmine-helper-bot)

## Разработка

Для развертывания дебаг сервера используется `docker-compose`. Развернуть его можно с помощью команды `docker compose up --build`. Перед этим необходимо установить следующие переменные окружения: 

- `TELEGRAM_BOT_KEY` с ключом телеграм бота
- `STORAGE_PASSPHRASE` с паролем для шифрования токенов

## Авторы

Алексей Королев
