package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/trebuhs/asa-cli/internal/output"
	"github.com/trebuhs/asa-cli/internal/services"
)

var geoCmd = &cobra.Command{
	Use:   "geo",
	Short: "Search geographic locations",
}

var geoSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for geo locations",
	RunE:  runGeoSearch,
}

var (
	geoQuery       string
	geoLimit       int
	geoOffset      int
	geoEntity      string
	geoCountryCode string
)

func init() {
	geoSearchCmd.Flags().StringVar(&geoQuery, "query", "", "Search query (required)")
	geoSearchCmd.Flags().IntVar(&geoLimit, "limit", 20, "Number of results")
	geoSearchCmd.Flags().IntVar(&geoOffset, "offset", 0, "Results offset")
	geoSearchCmd.Flags().StringVar(&geoEntity, "entity", "", "Entity type filter")
	geoSearchCmd.Flags().StringVar(&geoCountryCode, "country-code", "", "Country code filter")
	geoSearchCmd.MarkFlagRequired("query")

	geoCmd.AddCommand(geoSearchCmd)
	rootCmd.AddCommand(geoCmd)
}

func runGeoSearch(cmd *cobra.Command, args []string) error {
	client, err := newAPIClient()
	if err != nil {
		return err
	}

	svc := services.NewAppService(client)
	geos, _, err := svc.SearchGeo(geoQuery, geoLimit, geoOffset, geoEntity, geoCountryCode)
	if err != nil {
		return fmt.Errorf("searching geo locations: %w", err)
	}

	output.Print(getFormat(), geos, []output.Column{
		{Header: "ID", Field: "ID", Width: 10},
		{Header: "ENTITY", Field: "Entity", Width: 15},
		{Header: "NAME", Field: "DisplayName", Width: 30},
	})
	return nil
}
