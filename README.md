# Volcan

A highly configurable and powerful bot for Discord moderation on large guilds.

## Features

- Common moderation commands and features.
- Functionality-replicating API - anything available as a command can be done as an API call.

## Hosting

Volcan is provided as a self-hostable bot only. Due to the single-server nature of the bot, it is not available as a public bot.

### Setup

1. Clone the repository.
2. Copy `.env.example` to `.env`
3. Fill in the values in `.env` with your own.
4. Create a file named `config.yml` and fill it out with your desired config.
5. Run `docker-compose up -d` to start the bot.

### Environment Variables

Environment variables are required unless marked as optional with a `?` suffix in the constraint column.

| Name        | Constraint | Description                                    |
|-------------|------------|------------------------------------------------|
| `BOT_TOKEN` | string     | Your bot's token.                              |
| `BASE_URL`  | string     | The base URL for the API.                      |
| `API_TOKEN` | string     | The token to use to authenticate API requests. |
