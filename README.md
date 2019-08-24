# Telegram бот для Redmine

[![Build Status](https://travis-ci.org/alphatroya/redmine-helper-bot.svg?branch=master)](https://travis-ci.org/alphatroya/redmine-helper-bot)
[![codecov](https://codecov.io/gh/alphatroya/redmine-helper-bot/branch/master/graph/badge.svg)](https://codecov.io/gh/alphatroya/redmine-helper-bot)

Бот был создан для упрощения процесса заполнения часов в трекере задач Redmine. После быстрой настройки становится доступна команда `/fillhours` для старта пошагового процесса заполнения формы отчета.

Попробовать бот можно [по ссылке](https://t.me/redmine_buddy_bot)

## Разработка

Для развертывания дебаг сервера используется `docker-compose`. Развернуть его можно с помощью команды `docker compose up --build`. Перед этим необходимо установить переменную `TELEGRAM_BOT_KEY` с ключом телеграм бота.

## Авторы

Алексей Королев
