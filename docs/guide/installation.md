# Installation

> **Status**: Stable

To get started with GoW, you will need to have Go version 1.22 or higher installed on your machine.

## Installing the GoW CLI

### One-Line Install (Recommended)

**macOS / Linux**
```bash
curl -sSfL https://raw.githubusercontent.com/mechneerd/gow/main/install.sh | sh
```

**Windows (PowerShell)**
```powershell
iwr -useb https://raw.githubusercontent.com/mechneerd/gow/main/install.ps1 | iex
```

### Go Install

```bash
go install github.com/mechneerd/gow/cmd/gow@latest
```

After installation, verify:
```bash
gow --version
```

## Creating a New Project

The easiest way to scaffold a new GoW project is via the GoW CLI. This command will download the latest framework skeleton and configure your environment.

```bash
gow new my-app
cd my-app
```

## Running the Application

Once your project is scaffolded, use the GoW CLI to start the development server:

```bash
gow serve
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
