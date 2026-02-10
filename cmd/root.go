package cmd

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/trebuhs/asa-cli/internal/api"
	"github.com/trebuhs/asa-cli/internal/auth"
	"github.com/trebuhs/asa-cli/internal/config"
	"github.com/trebuhs/asa-cli/internal/models"
	"github.com/trebuhs/asa-cli/internal/output"
)

var (
	outputFormat string
	profileName  string
	verbose      bool
	noColor      bool
)

var rootCmd = &cobra.Command{
	Use:   "asa-cli",
	Short: "Apple Search Ads CLI",
	Long:  "A command-line interface for the Apple Search Ads Campaign Management API v5.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if noColor {
			color.NoColor = true
		}
		config.SetProfile(profileName)
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "table", "Output format: json or table")
	rootCmd.PersistentFlags().StringVarP(&profileName, "profile", "p", "", "Config profile name")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "Disable color output")
}

func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return err
	}
	return nil
}

// getFormat returns the output format.
func getFormat() output.Format {
	switch strings.ToLower(outputFormat) {
	case "json":
		return output.FormatJSON
	default:
		return output.FormatTable
	}
}

// newAPIClient creates an authenticated API client from config.
func newAPIClient() (*api.Client, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("loading config: %w", err)
	}

	if err := auth.ValidateConfig(cfg); err != nil {
		return nil, err
	}

	tokenProvider := auth.NewTokenProvider(cfg)
	transport := &auth.Transport{
		Token:   tokenProvider,
		OrgID:   cfg.OrgID,
		Verbose: verbose,
	}

	httpClient := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}

	client := api.NewClient(httpClient)
	client.Verbose = verbose
	return client, nil
}

// parseFilters parses filter strings like "status=ENABLED" into Conditions.
func parseFilters(filters []string) []models.Condition {
	var conditions []models.Condition
	for _, f := range filters {
		// Find operator (check multi-char operators first)
		for _, op := range []string{">=", "<=", "!~", "=", "~", "@", ">", "<"} {
			idx := strings.Index(f, op)
			if idx > 0 {
				field := f[:idx]
				value := f[idx+len(op):]
				apiOp := models.ParseFilterOperator(op)

				var values []string
				if op == "@" {
					values = strings.Split(value, ",")
				} else {
					values = []string{value}
				}

				conditions = append(conditions, models.Condition{
					Field:    field,
					Operator: apiOp,
					Values:   values,
				})
				break
			}
		}
	}
	return conditions
}

// parseSorts parses sort strings like "name:asc" into OrderByItems.
func parseSorts(sorts []string) []models.OrderByItem {
	var items []models.OrderByItem
	for _, s := range sorts {
		parts := strings.SplitN(s, ":", 2)
		field := parts[0]
		order := "ASCENDING"
		if len(parts) > 1 {
			order = models.ParseSortOrder(parts[1])
		}
		items = append(items, models.OrderByItem{
			Field:     field,
			SortOrder: order,
		})
	}
	return items
}

// exitWithError prints an error and exits with the given code.
func exitWithError(msg string, code int) {
	fmt.Fprintf(os.Stderr, "Error: %s\n", msg)
	os.Exit(code)
}
