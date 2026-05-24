package artisan

import (
	"fmt"

	"github.com/mechneerd/gow/foundation"

	"github.com/spf13/cobra"
)

// VendorPublishCmd publishes assets from PublishableProviders registered in the Application.
var VendorPublishCmd = &cobra.Command{
	Use:   "vendor:publish",
	Short: "Publish any publishable assets from service providers",
	Run: func(cmd *cobra.Command, args []string) {
		providerFlag, _ := cmd.Flags().GetString("provider")

		// We expect the Application to be passed from the console kernel
		appIface := cmd.Context().Value("app")
		app, ok := appIface.(*foundation.Application)
		if !ok || app == nil {
			fmt.Println("Error: Application not available in context")
			return
		}

		registry := app.ProviderRegistry()
		publishables := registry.Publishables()

		if len(publishables) == 0 {
			fmt.Println("No publishable providers registered.")
			return
		}

		published := 0
		for _, p := range publishables {
			name := fmt.Sprintf("%T", p)
			if providerFlag != "" && name != providerFlag {
				continue
			}

			fmt.Printf("Publishing from %s...\n", name)
			if err := foundation.PublishAssets(p, app.BasePath()); err != nil {
				fmt.Printf("  Error publishing %s: %v\n", name, err)
				continue
			}
			published++
		}

		if published == 0 {
			fmt.Println("Nothing was published.")
		} else {
			fmt.Printf("Published %d provider(s).\n", published)
		}
	},
}

func init() {
	VendorPublishCmd.Flags().StringP("provider", "p", "", "Publish assets from a specific provider only")
}

