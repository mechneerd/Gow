package artisan

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/spf13/cobra"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"

	"github.com/mechneerd/gow/bootstrap"
	"github.com/mechneerd/gow/cmd/tinker"
	"github.com/mechneerd/gow/config"
	"github.com/mechneerd/gow/container"
	"github.com/mechneerd/gow/database/orm"
)

var TinkerCmd = &cobra.Command{
	Use:   "tinker",
	Short: "Start an interactive Go REPL with the GoW application preloaded",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("GoW Tinker")
		fmt.Println("Type 'exit' or 'quit' to leave.")
		fmt.Println("Preloaded variables: App, DB, Config")
		fmt.Println("-------------------------------------------")

		// Boot the application
		app := bootstrap.NewApplication(".")

		// Resolve useful services
		cfg, _ := container.Make[*config.Repository](app.Container)

		var db *orm.DB
		if d, err := container.Make[*orm.DB](app.Container); err == nil {
			db = d
		}

		// Create interpreter
		i := interp.New(interp.Options{})
		i.Use(stdlib.Symbols)

		// Load our custom symbols (App, DB, Config)
		symbols := tinker.Symbols(app, db, cfg)
		i.Use(symbols)

		reader := bufio.NewReader(os.Stdin)

		for {
			fmt.Print(">>> ")
			line, err := reader.ReadString('\n')
			if err != nil {
				break
			}

			input := strings.TrimSpace(line)
			if input == "exit" || input == "quit" {
				fmt.Println("Goodbye!")
				break
			}
			if input == "" {
				continue
			}

			// Execute
			_, err = i.Eval(input)
			if err != nil {
				fmt.Println("Error:", err)
			}
		}
	},
}

func init() {
	// This will be picked up when make_commands.go runs
	_ = reflect.TypeOf(TinkerCmd)
}

