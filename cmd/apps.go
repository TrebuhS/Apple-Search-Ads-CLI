package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/trebuhs/asa-cli/internal/output"
	"github.com/trebuhs/asa-cli/internal/services"
)

var appsCmd = &cobra.Command{
	Use:   "apps",
	Short: "Search App Store apps",
}

var appsSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for apps by name",
	RunE:  runAppsSearch,
}

var (
	appQuery    string
	appLimit    int
	appOffset   int
	appOwnedOnly bool
)

func init() {
	appsSearchCmd.Flags().StringVar(&appQuery, "query", "", "Search query (required)")
	appsSearchCmd.Flags().IntVar(&appLimit, "limit", 20, "Number of results")
	appsSearchCmd.Flags().IntVar(&appOffset, "offset", 0, "Results offset")
	appsSearchCmd.Flags().BoolVar(&appOwnedOnly, "owned", false, "Return only owned apps")
	appsSearchCmd.MarkFlagRequired("query")

	appsCmd.AddCommand(appsSearchCmd)
	rootCmd.AddCommand(appsCmd)
}

func runAppsSearch(cmd *cobra.Command, args []string) error {
	client, err := newAPIClient()
	if err != nil {
		return err
	}

	svc := services.NewAppService(client)
	apps, _, err := svc.Search(appQuery, appLimit, appOffset, appOwnedOnly)
	if err != nil {
		return fmt.Errorf("searching apps: %w", err)
	}

	output.Print(getFormat(), apps, []output.Column{
		{Header: "ADAM ID", Field: "AdamID", Width: 12},
		{Header: "APP NAME", Field: "AppName", Width: 30},
		{Header: "DEVELOPER", Field: "DeveloperName", Width: 25},
	})
	return nil
}
