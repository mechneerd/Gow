package artisan

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var MakeRequestCmd = &cobra.Command{
	Use:   "make:request [name]",
	Short: "Create a new Form Request (validation struct)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		if !strings.HasSuffix(name, "Request") {
			name += "Request"
		}

		path := fmt.Sprintf("app/Http/Requests/%s.go", name)

		stub := `package requests

import (
	"gow/http/request"
)

type ` + name + ` struct {
	request.FormRequest
}

// Rules defines the validation rules for this request.
func (r *` + name + `) Rules() map[string]string {
	return map[string]string{
		// "title": "required|min:3|max:255",
		// "email": "required|email",
	}
}

// Messages returns custom validation error messages.
func (r *` + name + `) Messages() map[string]string {
	return map[string]string{
		// "title.required": "The title field is required.",
	}
}

// Authorize determines if the user is authorized to make this request.
func (r *` + name + `) Authorize() bool {
	return true
}
`
		generateFile(path, stub)
		fmt.Printf("Form Request created: %s\n", path)
	},
}
