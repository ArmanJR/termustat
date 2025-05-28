# Termustat

**The online timetabling tool for Iranian university students**  

Termustat sits on top of the [Golestan Integrated University System](https://fa.wikipedia.org/wiki/%DA%AF%D9%84%D8%B3%D8%AA%D8%A7%D9%86_(%D9%86%D8%B1%D9%85_%D8%A7%D9%81%D8%B2%D8%A7%D8%B1)), powering course planning at many Iranian universities.

![Termustat Screenshot](docs/screenshot.png)

## Features

- **Course Planner**: Easily add, remove, and rearrange courses on a visual calendar.
- **Conflict Detection**: Highlights time overlaps and warns of schedule clashes.
- **Dynamic Filtering**: Search and filter by faculty, professor, semester, or keywords.
- **Data Import/Export**: Uses the “engine” module to parse and transform Golestan exports into structured JSON.

## Tech Stack

| Component  | Language / Framework    |
| ---------- |-------------------------|
| **Backend**| Go          |
| **Engine** | Go                      |
| **Frontend**| React |
| **Database**| PostgreSQL              |
| **Containerization** | Docker & Docker Compose |


## Getting Started

### Configuration

Copy the example environment file and edit `.env` with your values:

```bash
cp .env.example .env
```

### Development

Make sure you have **Docker Engine** Installed. Following commands spins up the API, engine, and frontend in live‑reload mode:

| Command        | Purpose                                                          |
| -------------- | ---------------------------------------------------------------- |
| `make up`      | Start the stack (attached).                                      |
| `make build`   | Build / rebuild images.                                          |
| `make stop`    | Gracefully stop running containers without removing them.        |
| `make down`    | Stop containers **and** remove containers, networks and volumes. |
| `make restart` | Convenience alias: `make down` then `make up`.                   |
| `make logs`    | Follow logs for all services (`docker compose logs -f`).         |
| `make help`    | Show an auto-generated summary of all targets.                   |

- **Frontend** on `http://localhost:3000`
- **API** on `http://localhost:8080/api`

### Production

> [!WARNING]
> Termusat is currently undergoing a full refactor and is not yet ready for production use.

Build and run all services:

```bash
docker-compose up -d --build
```

## API Documentation

Auto‑generated Swagger docs are available at:

```
GET /api/docs/swagger.yaml
GET /api/docs/swagger.json
```

Or browse the UI:

```
http://localhost:8080/api/v1/swagger/index.html
```

## Architecture

```
┌──────────┐    ┌──────────────┐    ┌────────────┐
│ Golestan │───▶│  Engine      │───▶│  API       │───▶ PostgreSQL
│ Export   │    │ (Parser &    │    │ (Handlers, │
└──────────┘    │ Transformer) │    │  Services) │
                └──────────────┘    └────────────┘
                                         │
                                         ▼
                                      Frontend
                                      (React)
```

- **Engine**
    - Parses raw HTML/SQL exports from Golestan
    - Converts into normalized Go models
- **API**
    - Exposes REST endpoints for courses, faculties, semesters, users, etc.
    - Handles authentication, authorization, email workflows
- **Frontend**
    - Interactive calendar UI
    - Dynamic course selection & filtering


## Contributing

1. Fork the repo
2. Create a feature branch
3. Commit your changes & push
4. Open a Pull Request

Please ensure all new code is covered by tests and linted.


## License

This project is licensed under the **GNU General Public License v3.0**. See [LICENSE](LICENSE) for details.


## Resources

- [Arman Jafarnezhad](https://linkedin.com/in/ArmanJ), Author
- [Safar-o-Yek Magazine](https://t.me/sefroyekpub/43), “From Hesarak to Abbaspur” article
