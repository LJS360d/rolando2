# Rolando#7135

![Status](https://img.shields.io/website?url=http%3A%2F%2F188.245.46.41%3A25576)
![Discord](https://img.shields.io/discord/1122938014637756486)

Discord Bot implementation of a MarkovChain to learn how to type like a user of the guild.<br>
Created with: [discordgo](https://github.com/bwmarrin/discordgo)

## Overview

Rolando is a Discord bot that leverages Markov Chains to mimic the speech patterns of users within the server. The bot guesses the next word in a sentence based on the patterns it has learned.

## Credits

The concept is inspired by [`Fioriktos`](https://github.com/FiorixF1/fioriktos-bot), a Telegram bot using similar principles.

## How to Startup Your Own Rolando

- Prerequisites:
  - Go 1.23.0+
- Reccomended:
  - Make 4.4.1+

### Development:

1. Rename the file `example.env` in `.env` and follow the instructions written on it for the value of each needed enviroment variable
2. Run `go mod tidy` to download dependencies
2. Run `air` (`make dev` or `make` also work) to start a development server with hot reloading
   - Run `make build` to build the project
   - Run `make run` to start the application in development mode

### Troubleshooting:

- If you get an error like `Error: failed to create session: 40001: Unknown account`, make sure you have the `bot` scope enabled in your Discord application.

## Contact

If you have any questions or issues with Rolando, you can join the official Discord Server: [Join Here](https://discord.gg/tyrj7wte5b) or DM the creator directly, username: `zlejon`

## Support the Project

Consider supporting the project to help with hosting costs.
[![Buy Me a Coffee](https://img.shields.io/badge/Buy%20Me%20a%20Coffee-Support%20the%20Project-brightgreen)](https://www.buymeacoffee.com/rolandobot)

Thanks for your support!
