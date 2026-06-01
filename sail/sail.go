package sail

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"text/template"
)

// Config represents the Sail Docker configuration
type Config struct {
	AppName     string
	AppPort     int
	DBPort      int
	RedisPort   int
	QueuePort   int
	SSLPort     int
	MySQL       bool
	Postgres    bool
	Redis       bool
	Memory      string
	NodeVersion string
	PHPVersion  string
}

// DefaultConfig returns the default Sail configuration
func DefaultConfig() *Config {
	return &Config{
		AppName:     "gow-app",
		AppPort:     8000,
		DBPort:      3306,
		RedisPort:   6379,
		QueuePort:   9501,
		SSLPort:     443,
		MySQL:       true,
		Redis:       true,
		Memory:      "1G",
		NodeVersion: "20",
		PHPVersion:  "8.2",
	}
}

// Init creates a new Sail environment in the given directory
func Init(dir string, config *Config) error {
	if config == nil {
		config = DefaultConfig()
	}

	// Create docker-compose.yml
	if err := createDockerCompose(dir, config); err != nil {
		return fmt.Errorf("failed to create docker-compose.yml: %w", err)
	}

	// Create Dockerfile
	if err := createDockerfile(dir, config); err != nil {
		return fmt.Errorf("failed to create Dockerfile: %w", err)
	}

	// Create sail script
	if err := createSailScript(dir, config); err != nil {
		return fmt.Errorf("failed to create sail script: %w", err)
	}

	// Create .env.docker
	if err := createDockerEnv(dir, config); err != nil {
		return fmt.Errorf("failed to create .env.docker: %w", err)
	}

	// Create nginx config
	if err := createNginxConfig(dir, config); err != nil {
		return fmt.Errorf("failed to create nginx config: %w", err)
	}

	return nil
}

func createDockerCompose(dir string, config *Config) error {
	tmpl := `version: '3.8'

services:
  app:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: {{.AppName}}-app
    ports:
      - "{{.AppPort}}:8080"
    volumes:
      - .:/app
      - go-mod:/go/pkg/mod
      - go-build:/root/.cache/go-build
    environment:
      - APP_ENV=local
      - APP_DEBUG=true
      - DB_HOST=db
      - DB_PORT=3306
      - DB_DATABASE=gow
      - DB_USERNAME=sail
      - DB_PASSWORD=password
      - REDIS_HOST=redis
      - REDIS_PORT=6379
      - CACHE_DRIVER=redis
      - SESSION_DRIVER=redis
      - QUEUE_CONNECTION=redis
    depends_on:
      - db
      - redis
    networks:
      - sail
{{if .MySQL}}
  db:
    image: mysql:8.0
    container_name: {{.AppName}}-db
    ports:
      - "{{.DBPort}}:3306"
    environment:
      MYSQL_ROOT_PASSWORD: root
      MYSQL_DATABASE: gow
      MYSQL_USER: sail
      MYSQL_PASSWORD: password
    volumes:
      - mysql-data:/var/lib/mysql
    networks:
      - sail
{{end}}
  redis:
    image: redis:7-alpine
    container_name: {{.AppName}}-redis
    ports:
      - "{{.RedisPort}}:6379"
    volumes:
      - redis-data:/data
    networks:
      - sail

networks:
  sail:
    driver: bridge

volumes:
  go-mod:
  go-build:
{{if .MySQL}}
  mysql-data:
{{end}}
  redis-data:
`
	return writeTemplate(dir, "docker-compose.yml", tmpl, config)
}

func createDockerfile(dir string, config *Config) error {
	tmpl := `FROM golang:1.22-alpine

WORKDIR /app

# Install git and ca-certificates
RUN apk add --no-cache git ca-certificates tzdata

# Copy go mod files first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/gow

EXPOSE 8080

CMD ["/app/server"]
`
	return writeTemplate(dir, "Dockerfile", tmpl, config)
}

func createSailScript(dir string, config *Config) error {
	tmpl := `#!/bin/bash

# GoW Sail - Docker development environment helper

set -e

DOCKER_COMPOSE="docker compose"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

usage() {
    echo -e "${BLUE}GoW Sail${NC} - Docker development environment"
    echo ""
    echo "Usage: ./sail [command]"
    echo ""
    echo "Commands:"
    echo "  up              Start all containers"
    echo "  down            Stop all containers"
    echo "  restart         Restart all containers"
    echo "  logs            View container logs"
    echo "  shell           Open a shell in the app container"
    echo "  migrate         Run database migrations"
    echo "  seed            Run database seeders"
    echo "  fresh           Fresh migration (drop and recreate)"
    echo "  test            Run tests"
    echo "  artisan [cmd]   Run artisan command"
    echo "  build           Rebuild containers"
    echo "  tinker          Start Go REPL in app container"
    echo ""
}

case "$1" in
    up)
        $DOCKER_COMPOSE up -d
        echo -e "${GREEN}Containers started!${NC}"
        ;;
    down)
        $DOCKER_COMPOSE down
        echo -e "${YELLOW}Containers stopped!${NC}"
        ;;
    restart)
        $DOCKER_COMPOSE restart
        echo -e "${GREEN}Containers restarted!${NC}"
        ;;
    logs)
        $DOCKER_COMPOSE logs -f ${@:2}
        ;;
    shell)
        $DOCKER_COMPOSE exec app sh
        ;;
    migrate)
        $DOCKER_COMPOSE exec app ./server migrate ${@:2}
        ;;
    seed)
        $DOCKER_COMPOSE exec app ./server db:seed ${@:2}
        ;;
    fresh)
        $DOCKER_COMPOSE exec app ./server migrate:fresh ${@:2}
        ;;
    test)
        $DOCKER_COMPOSE exec app go test ./... ${@:2}
        ;;
    artisan)
        $DOCKER_COMPOSE exec app ./server ${@:2}
        ;;
    build)
        $DOCKER_COMPOSE build --no-cache
        echo -e "${GREEN}Containers rebuilt!${NC}"
        ;;
    tinker)
        $DOCKER_COMPOSE exec app go run main.go
        ;;
    *)
        usage
        ;;
esac
`
	path := filepath.Join(dir, "sail")
	if err := os.WriteFile(path, []byte(tmpl), 0755); err != nil {
		return err
	}
	return nil
}

func createDockerEnv(dir string, config *Config) error {
	tmpl := `APP_ENV=local
APP_DEBUG=true
APP_PORT={{.AppPort}}

DB_CONNECTION=mysql
DB_HOST=db
DB_PORT=3306
DB_DATABASE=gow
DB_USERNAME=sail
DB_PASSWORD=password

REDIS_HOST=redis
REDIS_PORT={{.RedisPort}}

CACHE_DRIVER=redis
SESSION_DRIVER=redis
QUEUE_CONNECTION=redis
`
	return writeTemplate(dir, ".env.docker", tmpl, config)
}

func createNginxConfig(dir string, config *Config) error {
	tmpl := `server {
    listen 80;
    server_name localhost;
    root /app/public;
    index index.html index.htm;

    charset utf-8;
    client_max_body_size 100M;

    location / {
        try_files $uri $uri/ /index.php?$query_string;
    }

    location ~ \.go$ {
        fastcgi_pass app:9000;
        fastcgi_index index.go;
        fastcgi_param SCRIPT_FILENAME $document_root$fastcgi_script_name;
        include fastcgi_params;
    }

    location ~ /\.ht {
        deny all;
    }
}
`
	path := filepath.Join(dir, "nginx.conf")
	return os.WriteFile(path, []byte(tmpl), 0644)
}

func writeTemplate(dir, filename, tmpl string, data any) error {
	t, err := template.New(filename).Parse(tmpl)
	if err != nil {
		return err
	}

	path := filepath.Join(dir, filename)
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return t.Execute(f, data)
}

// DetectSail checks if running inside Sail container
func DetectSail() bool {
	_, exists := os.LookupEnv("SAIL_CONTAINER")
	return exists
}

// GetPlatform returns the platform-specific docker-compose command
func GetPlatform() string {
	if runtime.GOOS == "windows" {
		return "docker compose"
	}
	return "docker compose"
}
