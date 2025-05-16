# Test News Application

A simple news application built with Go, MongoDB, and HTMX.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

## Prerequisites

- Go 1.23 or higher
- MongoDB 4.4 or higher
- Templ engine
- Docker and Docker Compose

## Make sure .env file is created with the following variables

```bash
PORT=8080
APP_ENV=local


BLUEPRINT_DB_HOST=localhost //default host for MongoDB
BLUEPRINT_DB_PORT=27017 //default port for MongoDB
BLUEPRINT_DB_USERNAME=your_username
BLUEPRINT_DB_ROOT_PASSWORD=your_password

```

## MakeFile

Note: If templ is not installing through the script, you should manually add the GOPATH.
You can do this by adding the following line to your `.zshrc` or `.bashrc` file: `export PATH="$PATH:$HOME/go/bin"`
Then run `source ~/.zshrc` or `source ~/.bashrc` to apply the changes.

Run build make command with tests

```bash
make all
```

Build the application

```bash
make build
```

Run the application

```bash
make run
```

Create DB container

```bash
make docker-run
```

Shutdown DB Container

```bash
make docker-down
```

DB Integrations Test:

```bash
make itest
```

Live reload the application:

```bash
make watch
```

Run the test suite:

```bash
make test
```

Clean up binary from the last build:

```bash
make clean
```

## API Endpoints

### Posts

- `GET /api/posts` - Get all posts
- `GET /api/posts/:id` - Get a specific post
- `POST /api/posts` - Create a new post
- `PUT /api/posts/:id` - Update a post
- `DELETE /api/posts/:id` - Delete a post

## Web Interface

The web interface is available at the following routes:

- `/web` - Home page
- `/web/posts` - List all posts
- `/web/posts/:id` - View a specific post
- `/web/upload` - Create a new post
- `/web/update` - Update a post
- `/web/delete` - Delete a post

## Key Technologies

- **Go**: Backend language
- **Gin**: Web framework
- **MongoDB**: Database
- **HTMX**: Frontend interactivity
- **Templ**: HTML templating
- **Tailwind CSS**: Styling

## Deployment

For notes on how to deploy the project on a live system, consider the following approaches:

### Using Docker

A Docker-compose is provided for containerized deployment. Build and run the Docker container:

```bash
docker compose up
```

### Manual Deployment

1. Build the application:

```bash
make build
```

2. Configure environment variables:

```bash
export MONGODB_URI="mongodb://your-production-mongodb-uri"
export PORT="8080"  # Or your preferred port
```

3. Run the application:

```bash
./test-news
```
