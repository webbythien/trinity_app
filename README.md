# trinity_app


# Database Setup with Docker Compose

This project uses Docker Compose to set up a PostgreSQL database. Follow the steps below to get started.

---

## Prerequisites

- [Docker](https://www.docker.com/) installed on your system.
- Basic knowledge of terminal/command-line usage.

---

## Folder Structure

project-root/ │ ├── db/ │ ├── docker-compose.yml │ ├── postgres/ │ ├── Dockerfile # Optional Dockerfile if additional configurations are required │ ├── init.sql # SQL file for initializing the database │


---

## How to Use

### 1. Clone the Repository

```bash
git clone https://github.com/webbythien/trinity_app.git
cd project-root/db

docker compose up -d --build

docker ps
```