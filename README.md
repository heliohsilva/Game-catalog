# 🎮 Game Catalog

A centralized, distraction-free game collection catalog designed to unify all your games across every launcher, console, and retro platform.

> **Objective:** Stop wondering what to play next. Game Catalog brings together your Steam, GOG, Epic, PlayStation, Switch, Xbox, emulation, and physical collections into one cohesive, fast dashboard with HowLongToBeat estimates and customizable platforms.

---

## ✨ Features

- **Unified Collection Hub**: Track all your games across PC (Steam, GOG, Epic, Ubisoft Connect, etc.), consoles (PlayStation, Switch, Xbox), retro emulators, and custom platforms.
- **Dynamic Platforms & Subplatforms**: Add or remove platforms on demand (e.g. Nintendo DS, PSP, Dreamcast) and configure custom subcategories / launchers (e.g. Ubisoft Connect, EA App, Cartridge, Homebrew).
- **HowLongToBeat (HLTB) Integration**: Fetch story and completion times automatically so you can pick a game that fits your schedule.
- **Status & Backlog Management**: Organize games by status (**Backlog**, **Playing**, **Completed**, **Abandoned**), genre, completion hours, and personal notes.
- **Hyprland / Waybar Aesthetic**: Clean, minimalist frosted-glass UI with static tactile micro-texture, high-legibility typography, and no distracting animations.
- **Omnibar Quick Search**: Instant search and filtering with `Ctrl+K` or `/`.
- **Theme Engine**: Switch between 10+ themes including Catppuccin Mocha, Tokyo Night, Nord, Gruvbox, Rose Pine, Dracula, and Cyberpunk.
- **Homelab & Self-Hosting Ready**: Single-port deployment (`3000`), automatic environment initialization, internal API proxy, and persistent SQLite storage.

---

## 🏗️ Architecture & Tech Stack

```
game-catalog/
├── api/          # Backend REST API (Go + Gin + Modernc SQLite)
├── fe/           # Frontend SPA (Next.js 15 App Router + React + Tailwind CSS + Lucide)
├── env/          # Docker compose and container definitions
├── data/         # Persistent SQLite database storage (ignored by git)
└── Makefile      # Simple three-target automation (make, make clean, make test)
```

- **Frontend**: [Next.js](https://nextjs.org/) (App Router), React 19, TypeScript, Tailwind CSS, Lucide Icons, Vitest.
- **Backend API**: [Go](https://golang.org/) (Gin-Gonic), SQLite with foreign keys and WAL mode.
- **Reverse Proxy**: Built-in Next.js App Router proxy (`/api/v1/*` ➔ `http://api:8080/api/v1/*`), eliminating CORS issues and keeping the raw database API isolated inside the Docker network.
- **Containerization**: Docker & Docker Compose with multi-stage builds.

---

## 🚀 Quick Start & Homelab Deployment

### Prerequisites

- [Docker](https://docs.docker.com/get-docker/) & Docker Compose
- *Optional for local development:* Go 1.22+, Node.js 20+

### 1. Clone & Start

```bash
git clone https://github.com/<your-username>/game-catalog.git
cd game-catalog
make
```

`make` will automatically initialize `env/.env` from `env/.env.example` if it doesn't exist, build both containers, and start them in the background.

### 2. Access the Application

- **Locally on the host:** [http://localhost:3000](http://localhost:3000)
- **From another PC or phone on your LAN:** `http://<your-homelab-ip>:3000`
- **Behind a Reverse Proxy (Nginx, Caddy, Traefik, Tailscale, Cloudflare Tunnel):** Point your proxy directly to port `3000`.

### Stop & Clean Up

```bash
make clean
```

### Run Tests

```bash
make test
```

---

## 🕹️ Managing Platforms & Subplatforms

You can configure platforms and launchers either via the UI or directly through the REST API.

### Via the Web Interface

1. Click the **"Platforms"** button on the top Waybar header.
2. **Add New Platform**: Enter a name (e.g., `Nintendo DS`) and optional comma-separated subplatforms (e.g., `Cartridge, R4, Digital`).
3. **Add Subplatform**: Click `+ Add Subplatform` on any platform card (e.g., add `Ubisoft Connect` or `EA App` under `PC`).
4. **Remove**: Click the remove icon on subplatform chips or platform headers (requires confirmation if games are assigned).

### Via REST API

```bash
# List all platforms and subcategories with game counts
curl http://localhost:3000/api/v1/platforms

# Add a new platform
curl -X POST http://localhost:3000/api/v1/platforms \
  -H "Content-Type: application/json" \
  -d '{"name": "Nintendo DS", "subcategories": ["Cartridge", "R4"]}'

# Add a subcategory to an existing platform
curl -X POST http://localhost:3000/api/v1/platforms/PC/subcategories \
  -H "Content-Type: application/json" \
  -d '{"name": "Ubisoft Connect"}'

# Delete a platform (use ?force=true if games are assigned)
curl -X DELETE "http://localhost:3000/api/v1/platforms/Nintendo%20DS?force=true"
```

---

## 📡 REST API Reference

All endpoints are accessible through the unified frontend port (`3000`):

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/health` | API health check and uptime (or direct on `8080`) |
| `GET` | `/api/v1/games` | List games with search/filter/pagination |
| `POST` | `/api/v1/games` | Add a new game to the catalog |
| `GET` | `/api/v1/games/:id` | Get details for a single game |
| `PUT` | `/api/v1/games/:id` | Update game metadata, status, or times |
| `DELETE` | `/api/v1/games/:id` | Remove a game from the catalog |
| `GET` | `/api/v1/platforms` | Get all platforms and subcategories |
| `POST` | `/api/v1/platforms` | Create a new platform |
| `DELETE` | `/api/v1/platforms/:name` | Remove a platform (`?force=true` supported) |
| `POST` | `/api/v1/platforms/:name/subcategories` | Add a subplatform to a platform |
| `DELETE` | `/api/v1/platforms/:name/subcategories/:sub` | Remove a subplatform |
| `GET` | `/api/v1/hltb?title=:title` | Query HowLongToBeat completion times |

---

## 💾 Data & Persistence

- The database file is stored at `data/game_catalog.sqlite`.
- The `data/` directory is mounted into the API container (`/app/data`), preserving all games, platforms, and personal notes across container updates and builds.
- The `data/` directory is excluded from version control (`.gitignore`), keeping your personal catalog private.

---

## ⌨️ Shortcuts

- **`/` or `Ctrl + K`**: Open Omnibar quick search
- **`Esc`**: Close open modals
