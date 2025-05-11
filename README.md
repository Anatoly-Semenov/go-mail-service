# go-mail-service

Сервис для отправки электронных писем, подписанный на Kafka топик.

## Требования

- Go 1.16 или выше
- Kafka
- SMTP сервер (например, Gmail)


## Запуск

```bash
make up
```

## Формат сообщений Kafka

Сервис подписан на топик `mail.send-email` и ожидает сообщения в следующем формате:

```json
{
    "user_id": "string",
    "email_type": "string"
}
```

Поддерживаемые типы писем:
- `registration` - письмо при регистрации
- `payment_reminder` - напоминание об оплате подписки
- `payment_success` - подтверждение успешной оплаты
- `password_change` - уведомление о смене пароля

## Структура проекта

```
.
├── cmd
│   └── app
│       └── main.go
├── configs
│   └── config.go
├── internal
│   ├── domain
│   │   └── email.go
│   ├── repository
│   │   ├── email_repository.go
│   │   └── email_repository_impl.go
│   ├── usecase
│   │   └── email_usecase.go
│   └── delivery
│       └── kafka
│           └── consumer.go
└── README.md
```