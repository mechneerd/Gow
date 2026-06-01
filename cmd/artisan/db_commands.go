package artisan

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var DbShowCmd = &cobra.Command{
	Use:   "db:show",
	Short: "Show all databases",
	Run: func(cmd *cobra.Command, args []string) {
		dsn := os.Getenv("DB_CONNECTION")
		if dsn == "" {
			dsn = os.Getenv("DATABASE_URL")
		}
		if dsn == "" {
			fmt.Println("No database connection configured.")
			fmt.Println("Set DB_CONNECTION or DATABASE_URL in your .env file.")
			return
		}

		db, err := sql.Open(getDriverName(), dsn)
		if err != nil {
			fmt.Println("Error connecting to database:", err)
			return
		}
		defer db.Close()

		rows, err := db.Query(getShowDatabasesQuery())
		if err != nil {
			fmt.Println("Error listing databases:", err)
			return
		}
		defer rows.Close()

		fmt.Println("Databases:")
		fmt.Println("----------")
		for rows.Next() {
			var name string
			if err := rows.Scan(&name); err != nil {
				continue
			}
			fmt.Println("  " + name)
		}
	},
}

var DbTableCmd = &cobra.Command{
	Use:   "db:table [table]",
	Short: "Show information about a table",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		table := args[0]
		dsn := os.Getenv("DB_CONNECTION")
		if dsn == "" {
			dsn = os.Getenv("DATABASE_URL")
		}
		if dsn == "" {
			fmt.Println("No database connection configured.")
			return
		}

		db, err := sql.Open(getDriverName(), dsn)
		if err != nil {
			fmt.Println("Error connecting to database:", err)
			return
		}
		defer db.Close()

		rows, err := db.Query(getDescribeTableQuery(table))
		if err != nil {
			fmt.Printf("Error describing table '%s': %v\n", table, err)
			return
		}
		defer rows.Close()

		cols, _ := rows.Columns()
		fmt.Printf("Table: %s\n", table)
		fmt.Println(strings.Repeat("-", 60))
		fmt.Printf("%-30s %-15s %-10s\n", "Column", "Type", "Nullable")
		fmt.Println(strings.Repeat("-", 60))

		for rows.Next() {
			values := make([]sql.NullString, len(cols))
			ptrs := make([]any, len(cols))
			for i := range values {
				ptrs[i] = &values[i]
			}
			if err := rows.Scan(ptrs...); err != nil {
				continue
			}
			colName := values[0].String
			colType := values[1].String
			nullable := "YES"
			if len(values) > 2 && values[2].String == "NO" {
				nullable = "NO"
			}
			fmt.Printf("%-30s %-15s %-10s\n", colName, colType, nullable)
		}
	},
}

func getDriverName() string {
	conn := os.Getenv("DB_CONNECTION")
	switch strings.ToLower(conn) {
	case "mysql":
		return "mysql"
	case "postgres", "pgsql":
		return "postgres"
	case "sqlite":
		return "sqlite"
	default:
		return "sqlite"
	}
}

func getShowDatabasesQuery() string {
	conn := os.Getenv("DB_CONNECTION")
	switch strings.ToLower(conn) {
	case "mysql":
		return "SHOW DATABASES"
	case "postgres", "pgsql":
		return "SELECT datname FROM pg_database WHERE datistemplate = false"
	default:
		return "SELECT name FROM pragma_database_list"
	}
}

func getDescribeTableQuery(table string) string {
	conn := os.Getenv("DB_CONNECTION")
	switch strings.ToLower(conn) {
	case "mysql":
		return fmt.Sprintf("DESCRIBE `%s`", table)
	case "postgres", "pgsql":
		return fmt.Sprintf("SELECT column_name, data_type, is_nullable FROM information_schema.columns WHERE table_name = '%s'", table)
	default:
		return fmt.Sprintf("PRAGMA table_info('%s')", table)
	}
}
