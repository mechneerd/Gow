# Phase 5 Implementation Plan (Advanced Features)

This document outlines the detailed implementation of Phase 5 within the GoW framework.

## 1. Queue System & Worker [P5.1, P5.2]
- **`queue/manager.go`**: Resolves queue drivers and exposes the central `Push(job)` API.
- **`queue/job.go`**: Defines the fundamental `Job` interface with `Handle()` and `Failed(error)` methods.
- **`queue/driver_sync.go`**: Provides an immediate, Go-native channel-based queue driver (`SyncDriver`) for high-velocity local development.
- **`queue/worker.go`**: A scalable worker daemon that pops jobs from the queue and executes them reliably.

## 2. Event Dispatcher [P5.3]
- **`events/dispatcher.go`**: Utilizes `reflect.Type` to create a robust pub/sub event dispatcher. 
- Features `Listen()` for specific struct events and `ListenAny()` for global wildcard listeners.

## 3. Mail & Notifications [P5.4, P5.5]
- **`mail/mailer.go` & `mail/message.go`**: Introduces a fluent `Mailable` interface to build email payloads (`From()`, `To()`, `HTML()`). Includes `LogDriver` and `SmtpDriver` integrations.
- **`notifications/manager.go` & `notifications/channel.go`**: Allows classes to implement the `Notifiable` interface. The system inspects the `Via()` method to dispatch notifications concurrently across Mail, Database, or Custom channels.

## 4. Task Scheduling [P5.6]
- **`console/schedule.go`**: Wraps the battle-tested `robfig/cron/v3` parser into a fluent Laravel-style DSL (e.g., `Schedule().Command("sync:db").EveryMinute().WithoutOverlapping()`).
- **`cmd/artisan/schedule_run.go`**: Registers the `artisan schedule:run` command to act as the master blocking cron daemon for the framework.

## 5. Broadcasting [P5.7]
- **`broadcasting/manager.go` & `broadcasting/channel.go`**: Establishes the real-time websocket foundation, mapping out `PublicChannel`, `PrivateChannel`, and `PresenceChannel` structures and abstracting driver logic for easy future integration with Pusher or Redis.
