> ✅ Implemented · 🚧 In Progress · 📋 Planned

# Database & Tooling

> **Status**: ✅ Implemented


GoW ships with a powerful set of database tools beyond just the ORM, allowing you to manage schemas, profile queries, and prune obsolete records seamlessly.

## Schema Dumping

During development, your database schema can grow quite large. Instead of running hundreds of migration files for local setup or CI testing, you can squash your schema into a single SQL dump file.

```bash
artisan schema:dump
```

This command will generate a `.sql` file in your `database/schema` directory. When running fresh migrations in the future, GoW will execute this SQL dump first before running any new, outstanding migrations.

## Model Pruning

Many applications accumulate obsolete records over time (e.g., expired tokens, old logs). GoW allows you to define a `Prunable` interface on your models.

```go
type LogEntry struct {
    ID        int
    Message   string
    CreatedAt time.Time
}

func (l *LogEntry) PrunableQuery() any {
    return orm.Table("log_entries").Where("created_at", "<", time.Now().AddDate(0, -1, 0))
}
```

Once defined, you can routinely run the Artisan command (via cron or scheduler) to automatically delete these records:

```bash
artisan model:prune
```

## Query Profiling

GoW's Query Builder natively dispatches `QueryEvent` hooks upon every SQL execution. This automatically captures the raw SQL string, execution time in milliseconds, and the caller stack trace. 

This enables automated N+1 query detection during development and allows you to dump slow queries into your application logs.
