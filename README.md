# Conversense - AI-Powered Meeting Intelligence

Conversense is a cutting-edge platform designed to revolutionize how teams conduct and document meetings. By integrating real-time video conferencing with advanced AI capabilities, it provides automated transcription, intelligent summarization, and actionable insights.

---

## Features

- **Real-time Video Conferencing**: Seamless, high-quality video calls powered by **LiveKit**.
- **AI-Driven Insights**: Automated transcription and intelligent summarization using **Google Gemini**.
- **Interactive Dashboard**: Manage meetings, view history, and access AI-generated reports.
- **Smart Agents**: AI agents that participate in meetings to assist, record, and provide real-time support.
- **Secure Authentication**: Robust user management and authentication.

---

## Technology Stack

### Backend

- **Go** 1.25+: Core backend language.
- **Gin**: High-performance HTTP web framework.
- **PostgreSQL**: Primary relational database.
- **SQLC**: Type-safe Go code generation from SQL.
- **LiveKit Server SDK**: For managing real-time video and audio.
- **Google GenAI SDK**: For AI-powered transcription and summarization.
- **Inngest**: For reliable background job processing and event-driven workflows.
- **AWS S3**: For secure storage of meeting recordings and artifacts.

### Frontend

- **React**: UI library for building interactive interfaces.
- **Vite**: Next-generation frontend tooling.
- **Tailwind CSS**: Utility-first CSS framework for rapid UI development.
- **Shadcn/UI**: Reusable components built with Radix UI and Tailwind.
- **TanStack Router**: Type-safe routing for React applications.
- **TanStack Query**: Powerful asynchronous state management.
- **BetterAuth**: Comprehensive authentication solution.

---

## Getting Started

### Prerequisites

- **Go** 1.25+
- **Node.js** 18+
- **Docker** & **Docker Compose**
- **mprocs** (optional, for running multiple processes in one terminal)

### Installation & Setup

1.  **Clone the repository**

    ```bash
    git clone https://github.com/rahulSailesh-shah/converSense.git
    cd converSense
    ```

2.  **Environment Configuration**
    Copy the example environment file and configure your credentials:

    ```bash
    cp .env.example .env
    ```

    _Note: You will need to add credentials for LiveKit, Google Gemini, and AWS to fully enable all features._

3.  **Start Infrastructure**
    Spin up the database and other services using Docker:

    ```bash
    make docker-up
    ```

4.  **Run Migrations**
    Apply database migrations:

    ```bash
    make migrate-up
    ```

5.  **Run the Application**
    The recommended way to run the full stack (backend, frontend, auth, studio) is using `mprocs` via the make command:
    ```bash
    make dev
    ```
    Alternatively, you can run services individually:
    - **Backend**: `make run-backend`
    - **Frontend**: `make run-frontend`
    - **Inngest**: `make run-inngest`

---

## Environment Variables

Ensure your `.env` file includes the following configurations:

```env
# App
PORT=8080
APP_ENV=development

# Database
DB_URL="postgresql://admin:admin@localhost:5432/dbname"

# Authentication
JWKS_URL=http://localhost:3000/api/auth/jwks

# LiveKit (Required for video)
LIVEKIT_API_KEY=your_key
LIVEKIT_API_SECRET=your_secret
LIVEKIT_URL=your_url

# Google Gemini (Required for AI)
GEMINI_API_KEY=your_key

# AWS S3 (Required for storage)
AWS_REGION=us-east-1
AWS_ACCESS_KEY_ID=your_key
AWS_SECRET_ACCESS_KEY=your_secret
AWS_BUCKET_NAME=your_bucket
```

---

## Development Tools

- **Makefile**: Centralized commands for building, running, and maintaining the project.
- **SQLC**: Generates type-safe Go code from SQL queries. Configured in `sqlc.yml`.
- **Goose**: Handles database migrations.
- **Air**: Live reload for Go development (optional, used in `make watch`).

## License

MIT
