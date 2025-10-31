# Go-React Starter Kit

A modern full-stack starter template combining Go (backend) and React (frontend) for rapid web application development.

---

## Project Structure

```
├── cmd/api/                # Go backend entrypoint
│   └── main.go
├── internal/               # Backend application logic
│   ├── app/                # App initialization
│   ├── db/                 # Database repo & SQL queries
│   ├── server/             # HTTP server setup
│   ├── service/            # Business logic/services
│   └── transport/http/     # HTTP routing, handlers, middleware
├── pkg/                    # Shared backend packages
│   ├── auth/               # JWT & authentication
│   ├── config/             # App configuration
│   └── database/           # DB connection & migrations
├── frontend/               # React frontend (Vite, TS, Tailwind)
│   ├── src/                # App source code
│   │   ├── components/     # UI components (shadcn/ui)
│   │   ├── hooks/          # Custom React hooks
│   │   ├── features/       # Feature modules
│   │   ├── lib/            # Utility functions
│   │   ├── routes/         # TanStack Router routes
│   │   └── assets/         # Static assets
│   ├── public/             # Static files
│   ├── index.html          # App entry HTML
│   ├── package.json        # Frontend dependencies
│   ├── vite.config.ts      # Vite config
│   ├── tsconfig*.json      # TypeScript configs
│   └── components.json     # shadcn/ui config
├── Makefile                # Common dev commands
├── docker-compose.yml      # Local dev services (PostgreSQL)
├── sqlc.yml                # SQLC codegen config
├── go.mod                  # Go dependencies
└── README.md               # Project documentation
```

---

## Technologies Used

### Backend

- **Go** 1.25+
- **Gin** (HTTP framework)
- **PostgreSQL** (via Docker)
- **SQLC** (type-safe DB access)
- **JWT Auth** (jwx)
- **Goose** (DB migrations)
- **Makefile** (build, run, migrate)

### Frontend

- **React** (TypeScript)
- **Vite** (bundler)
- **TanStack Router** (routing)
- **TanStack Query** (data fetching)
- **Tailwind CSS** (utility-first styling)
- **shadcn/ui** (component library)
- **Radix UI** (accessible primitives)
- **ESLint** (linting)
- **BetterAuth** (auth management)

---

## Getting Started

### Prerequisites

- Go 1.25+
- Node.js 18+
- Docker

### Backend Setup

1. Copy `.env.example` to `.env` and fill in DB credentials.
2. Start PostgreSQL:
   ```sh
   make docker-up
   ```
3. Run DB migrations:
   ```sh
   make migrate-up
   ```
4. Build & run backend:
   ```sh
   make build
   ./main
   # or
   make run-backend
   ```

### Frontend Setup

1. Install dependencies:
   ```sh
   cd frontend
   npm install
   ```
2. Start dev server:
   ```sh
   npm run dev
   ```
3. Access at [http://localhost:5173](http://localhost:5173)

---

## Development Tools

- **Makefile**: Build, run, clean, migrations, Docker
- **Docker Compose**: Local PostgreSQL
- **SQLC**: Go codegen from SQL
- **Goose**: DB migrations
- **ESLint**: Linting (frontend)
- **Tailwind**: Styling
- **shadcn/ui**: UI components

---

## Useful Commands

Backend:

- `make build` — Build Go backend
- `make run-backend` — Run backend
- `make docker-up` / `make docker-down` — Start/stop DB
- `make migrate-up` / `make migrate-down` — DB migrations

Frontend:

- `npm run dev` — Start frontend dev server
- `npm run build` — Build frontend
- `npm run lint` — Lint code

---

## License

MIT
