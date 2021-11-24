# tfs-trading-bot

![go workflow](https://github.com/khodand/tfs-trading-bot/actions/workflows/go.yml/badge.svg)

# Инструкция по запуску
## Заполнить файл конфигурации
- Создать в корне файл `config.json` и заполнить его
```json
{
  "dsn": "<postgres dsn>",
  "telegram": "<telegram bot token>",
  "algoPeriod": "<int selling period for EMA algo>",
  "krakenWebsocket": "<kraken websocket url>",
  "krakenREST": "<kraken restapi url>",
  "krakenPublicKey": "<kraken public key>",
  "krakenSecretKey": "<kraken secret key>"
}
```

## Запустить бота
C помощью консоли 
`go run .\cmd\api\main.go` или IDE

## Познакомиться с телеграм ботом
Начать общение с вашим ботом и написать ему ключевое сообщение `/start`

## Управляйте ботом при помощи Rest API
- Чтобы начать торговать какой-то бумагой используйте `http://localhost:5000/trade/{tickerSymbol}`
```http request
http://localhost:5000/trade/pi_ethusd
```
- Для смены алгоритма торговли: `http://localhost:5000/algo/{algorithm}/{period}`
```http request
http://localhost:5000/EMA/20
```
- Для остановки торговли
```http request
http://localhost:5000/stop
```
