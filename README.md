# Hunt Royale Discord Hunt Helper

A Discord bot originally developed for the official Hunt Royale Discord server to facilitate global matchmaking and cross-server game posting.

> [!NOTE]
> This project is now open source for demonstration purposes only. The bot is no longer online, and all credentials (tokens, database passwords) have been rotated.

## Project Overview

The main purpose of this bot was to help players find games ("Hunts") more easily. It featured a system where players could request a game, and the bot would post this request across multiple Discord servers, including the player's in-game ID and the specific game mode they wanted to play.

## Features

*   **Dungeon Finder**: The core feature that allowed automated group finding and posting across servers.
*   **Player ID Registry**: A system to register and look up player game IDs, making it easier to connect in-game.
*   **Cross-Server Communication**: bridged the gap between different community servers.

## Tech Stack

This project is built using:

*   **Language**: Go (Golang) 1.23
*   **Discord Library**: [DiscordGo](https://github.com/bwmarrin/discordgo)
*   **Database**: MySQL with [GORM](https://gorm.io/)
*   **CLI Framework**: [Cobra](https://github.com/spf13/cobra)
*   **Logging**: Logrus

## Project Structure

*   `cmd/`: Entry point for the application (CLI commands).
*   `commands/`: Contains the logic for the bot's slash commands and interactions.
*   `handlers/`: Event handlers for Discord events (messages, interactions, etc.).
*   `internal/`: Internal packages and utilities.

## Disclaimer

I no longer maintain or run this bot. This repository serves as a portfolio piece to showcase my skills in Go and Discord bot development. It was written entirely from scratch without the use of AI generation.