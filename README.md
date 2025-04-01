# JWT Authentication Practice Project

This project is a practice implementation of JSON Web Token (JWT) authentication using Go. It demonstrates user
registration, login, token generation (access and refresh tokens), and protected routes, integrated with a PostgreSQL
database. The project includes a Makefile for task automation and a Docker Compose setup for running a local PostgreSQL
instance.

## Features

- User registration and login with JWT-based authentication.
- Generation of access (short-lived) and refresh (long-lived) tokens.
- Protected endpoints accessible only with a valid access token.
- Token refresh mechanism using refresh tokens.
- Database migrations using `goose`.
- Mock generation for testing with `mockgen`.
- Dockerized PostgreSQL for local development.

## Prerequisites

- Go (version 1.21 or later recommended).
- Docker and Docker Compose (for running PostgreSQL locally).
- PostgreSQL (if running without Docker).
- Make (for using the Makefile).

## Setup

1. **Clone the repository:**
   ```bash
   git clone github.com.sanchey92/jwt-example
   cd <project_folder>

2. **Create a `.env` file in the root directory with the following environment variables:**
    ```env
   PORT=8080
   MIGRATION_DIR=./migrations
   POSTGRES_DB=jwt
   POSTGRES_USER=jwt-user
   POSTGRES_PASSWORD=jwt-password
   PG_DSN="host=localhost port=5432 dbname=jwt user=jwt-user password=jwt-password sslmode=disable"
   JWT_ACCESS_SECRET=your_access_secret
   JWT_REFRESH_SECRET=your_refresh_secret
   JWT_ACCESS_TTL=15
   JWT_REFRESH_TTL=7
   ``` 

   **Environment Variables Description**

   - PORT: Port for the application server (default: 8080).
   - MIGRATION_DIR: Directory containing migration files.
   - POSTGRES_DB: PostgreSQL database name.
   - POSTGRES_USER: PostgreSQL user.
   - POSTGRES_PASSWORD: PostgreSQL password.
   - PG_DSN: PostgreSQL connection string.
   - JWT_ACCESS_SECRET: Secret key for signing access tokens (replace with a secure value).
   - JWT_REFRESH_SECRET: Secret key for signing refresh tokens (replace with a secure value).
   - JWT_ACCESS_TTL: Access token TTL in minutes (15 minutes).
   - JWT_REFRESH_TTL: Refresh token TTL in days (7 days).

3. **Install dependencies:**
   ```bash
      make init-deps
   
4. **Set up PostgreSQL using docker container:**
   ```bash 
      docker-compose up -d

5. **Apply migrations:**
   ```bash
      make local-migrations-up

6. **Build and run the application:**
   ```bash
      make run

## Notes

- Replace JWT_ACCESS_SECRET and JWT_REFRESH_SECRET with strong, unique values for security.
- Adjust JWT_ACCESS_TTL and JWT_REFRESH_TTL as needed.
- The project structure assumes mocks are generated in internal/service/mocks.

## License
This project is licensed under the MIT License (or specify another if applicable).