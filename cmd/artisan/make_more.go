package artisan

import (
	"fmt"
	"github.com/spf13/cobra"
)

var MakeMailCmd = &cobra.Command{
	Use:   "make:mail [name]",
	Short: "Create a new mailable class",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		stub := `package mail

import "github.com/mechneerd/gow/mail"

type ` + name + ` struct {
	mail.Mailable
}

func New` + name + `() *` + name + ` {
	return &` + name + `{}
}
`
		path := "app/Mail/" + name + ".go"
		generateFile(path, stub)
		fmt.Println("Mail class created. Remember to register it if needed.")
	},
}

var MakeEventCmd = &cobra.Command{
	Use:   "make:event [name]",
	Short: "Create a new event class",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		stub := `package events

type ` + name + ` struct {
	// Add your payload here
}

func (e ` + name + `) BroadcastOn() []string { return []string{} }
func (e ` + name + `) BroadcastAs() string   { return "` + name + `" }
func (e ` + name + `) BroadcastWith() map[string]any { return map[string]any{} }
`
		path := "app/Events/" + name + ".go"
		generateFile(path, stub)
	},
}

var MakeListenerCmd = &cobra.Command{
	Use:   "make:listener [name]",
	Short: "Create a new event listener",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		stub := `package listeners

import "fmt"

type ` + name + ` struct{}

func (l *` + name + `) Handle(event any) {
	fmt.Printf("Handling event: %T\n", event)
}
`
		path := "app/Listeners/" + name + ".go"
		generateFile(path, stub)
	},
}

var MakePolicyCmd = &cobra.Command{
	Use:   "make:policy [name]",
	Short: "Create a new authorization policy",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		stub := `package policies

type ` + name + ` struct{}

func (p *` + name + `) Before(user any, ability string) bool { return false }

func (p *` + name + `) ` + name + `(user any, model any) bool {
	return true
}
`
		path := "app/Policies/" + name + ".go"
		generateFile(path, stub)
	},
}

var MakeResourceCmd = &cobra.Command{
	Use:   "make:resource [name]",
	Short: "Create a new API resource",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		stub := `package resources

import "github.com/mechneerd/gow/http/resources"

type ` + name + ` struct {
	resources.Resource
}

func (r ` + name + `) ToArray() map[string]any {
	return map[string]any{
		"id": r.Get("ID"),
	}
}
`
		path := "app/Http/Resources/" + name + ".go"
		generateFile(path, stub)
	},
}

// MakeJobCmd is defined in make_commands.go (authoritative location).
// Duplicate removed here to fix redeclaration.

var MakeNotificationCmd = &cobra.Command{
	Use:   "make:notification [name]",
	Short: "Create a new notification",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		stub := `package notifications

import "github.com/mechneerd/gow/notifications"

type ` + name + ` struct {
	notifications.Notification
}

func (n *` + name + `) ToMail(notifiable any) *notifications.MailMessage {
	return notifications.NewMailMessage().Subject("Notification").Line("Hello")
}
`
		path := "app/Notifications/" + name + ".go"
		generateFile(path, stub)
	},
}

var MakeTestCmd = &cobra.Command{
	Use:   "make:test [name]",
	Short: "Create a new test file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		stub := `package tests

import "testing"

func Test` + name + `(t *testing.T) {
	// Write your test here
}
`
		path := "tests/" + name + "_test.go"
		generateFile(path, stub)
	},
}

