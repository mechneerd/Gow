# Installation

To get started with GoW, you will need to have Go version 1.24 or higher installed on your machine.

## Creating a New Project

The easiest way to scaffold a new GoW project is via our npm initializer. This command will download the latest framework skeleton and configure your environment.

```bash
npx gow new my-app
cd my-app
```

## Running the Application

Once your project is scaffolded, you can use the Artisan CLI (or native Go commands) to run your server.

```bash
# Start the HTTP server on port 8080
go run cmd/app/main.go
```

## Environment Configuration

GoW utilizes a `.env` file at the root of your project to manage environment-specific variables like database credentials. 

```env
APP_NAME=GoW_App
APP_ENV=local
APP_KEY=base64:your_generated_key_here
APP_DEBUG=true
APP_URL=http://localhost:8080

DB_CONNECTION=sqlite
DB_DATABASE=database/database.sqlite
```
