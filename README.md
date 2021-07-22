# Telegram бот для Redmine

[![Go Report Card](https://goreportcard.com/badge/github.com/alphatroya/redmine-helper-bot)](https://goreportcard.com/report/github.com/alphatroya/redmine-helper-bot)
[![codecov](https://codecov.io/gh/alphatroya/redmine-helper-bot/branch/master/graph/badge.svg)](https://codecov.io/gh/alphatroya/redmine-helper-bot)

Бот был создан для упрощения процесса заполнения часов в трекере задач Redmine. После быстрой настройки становится доступна команда `/fillhours` для старта пошагового процесса заполнения формы отчета.

Попробовать бот можно [по ссылке](https://t.me/redmine_buddy_bot)

## Разработка

Для развертывания дебаг сервера используется `docker-compose`. Развернуть его можно с помощью команды `docker compose up --build`. Перед этим необходимо установить следующие переменные окружения: 

- `TELEGRAM_BOT_KEY` с ключом телеграм бота
- `STORAGE_PASSPHRASE` с паролем для шифрования токенов

## Авторы

Алексей Королев
