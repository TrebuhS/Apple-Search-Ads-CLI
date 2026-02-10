package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/trebuhs/asa-cli/internal/models"
	"github.com/trebuhs/asa-cli/internal/output"
	"github.com/trebuhs/asa-cli/internal/services"
)

var campaignsCmd = &cobra.Command{
	Use:   "campaigns",
	Short: "Manage campaigns",
}

var campaignsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all campaigns",
	RunE:  runCampaignsList,
}

var campaignsGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get a campaign by ID",
	Args:  cobra.ExactArgs(1),
	RunE:  runCampaignsGet,
}

var campaignsFindCmd = &cobra.Command{
	Use:   "find",
	Short: "Find campaigns with filters",
	RunE:  runCampaignsFind,
}

var campaignsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new campaign",
	RunE:  runCampaignsCreate,
}

var campaignsUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a campaign",
	Args:  cobra.ExactArgs(1),
	RunE:  runCampaignsUpdate,
}

var campaignsDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a campaign",
	Args:  cobra.ExactArgs(1),
	RunE:  runCampaignsDelete,
}

var (
	campLimit     int
	campOffset    int
	campFilters   []string
	campSorts     []string
	campAll       bool
	campName      string
	campBudget    string
	campDaily     string
	campCountries string
	campAppID     int64
	campStatus    string
)

func init() {
	// list
	campaignsListCmd.Flags().IntVar(&campLimit, "limit", 20, "Number of results")
	campaignsListCmd.Flags().IntVar(&campOffset, "offset", 0, "Results offset")

	// find
	campaignsFindCmd.Flags().StringSliceVar(&campFilters, "filter", nil, `Filter conditions (e.g. "status=ENABLED", "name~MyApp")`)
	campaignsFindCmd.Flags().StringSliceVar(&campSorts, "sort", nil, `Sort order (e.g. "name:asc", "id:desc")`)
	campaignsFindCmd.Flags().IntVar(&campLimit, "limit", 20, "Number of results")
	campaignsFindCmd.Flags().IntVar(&campOffset, "offset", 0, "Results offset")
	campaignsFindCmd.Flags().BoolVar(&campAll, "all", false, "Fetch all pages")

	// create
	campaignsCreateCmd.Flags().StringVar(&campName, "name", "", "Campaign name (required)")
	campaignsCreateCmd.Flags().StringVar(&campBudget, "budget", "", "Total budget (e.g. 1000.00)")
	campaignsCreateCmd.Flags().StringVar(&campDaily, "daily-budget", "", "Daily budget (e.g. 50.00)")
	campaignsCreateCmd.Flags().StringVar(&campCountries, "countries", "", "Comma-separated country codes (e.g. US,GB)")
	campaignsCreateCmd.Flags().Int64Var(&campAppID, "app-id", 0, "App Adam ID (required)")
	campaignsCreateCmd.Flags().StringVar(&campStatus, "status", "ENABLED", "Campaign status")
	campaignsCreateCmd.MarkFlagRequired("name")
	campaignsCreateCmd.MarkFlagRequired("app-id")
	campaignsCreateCmd.MarkFlagRequired("countries")
	campaignsCreateCmd.MarkFlagRequired("budget")
	campaignsCreateCmd.MarkFlagRequired("daily-budget")

	// update
	campaignsUpdateCmd.Flags().StringVar(&campName, "name", "", "Campaign name")
	campaignsUpdateCmd.Flags().StringVar(&campBudget, "budget", "", "Total budget")
	campaignsUpdateCmd.Flags().StringVar(&campDaily, "daily-budget", "", "Daily budget")
	campaignsUpdateCmd.Flags().StringVar(&campStatus, "status", "", "Campaign status (ENABLED/PAUSED)")

	campaignsCmd.AddCommand(campaignsListCmd, campaignsGetCmd, campaignsFindCmd, campaignsCreateCmd, campaignsUpdateCmd, campaignsDeleteCmd)
	rootCmd.AddCommand(campaignsCmd)
}

var campaignColumns = []output.Column{
	{Header: "ID", Field: "ID", Width: 12},
	{Header: "NAME", Field: "Name", Width: 30},
	{Header: "STATUS", Field: "Status", Width: 10},
	{Header: "SERVING", Field: "ServingStatus", Width: 12},
	{Header: "BUDGET", Field: "BudgetAmount", Width: 15},
	{Header: "DAILY BUDGET", Field: "DailyBudgetAmount", Width: 15},
	{Header: "COUNTRIES", Field: "CountriesOrRegions", Width: 15},
}

func runCampaignsList(cmd *cobra.Command, args []string) error {
	client, err := newAPIClient()
	if err != nil {
		return err
	}

	svc := services.NewCampaignService(client)
	campaigns, _, err := svc.List(campLimit, campOffset)
	if err != nil {
		return fmt.Errorf("listing campaigns: %w", err)
	}

	output.Print(getFormat(), campaigns, campaignColumns)
	return nil
}

func runCampaignsGet(cmd *cobra.Command, args []string) error {
	id, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid campaign ID: %s", args[0])
	}

	client, err := newAPIClient()
	if err != nil {
		return err
	}

	svc := services.NewCampaignService(client)
	campaign, err := svc.Get(id)
	if err != nil {
		return fmt.Errorf("getting campaign: %w", err)
	}

	output.Print(getFormat(), campaign, campaignColumns)
	return nil
}

func runCampaignsFind(cmd *cobra.Command, args []string) error {
	client, err := newAPIClient()
	if err != nil {
		return err
	}

	selector := models.NewSelector(campLimit, campOffset)
	selector.Conditions = parseFilters(campFilters)
	selector.OrderBy = parseSorts(campSorts)

	svc := services.NewCampaignService(client)

	if campAll {
		campaigns, err := svc.FindAll(selector)
		if err != nil {
			return fmt.Errorf("finding campaigns: %w", err)
		}
		output.Print(getFormat(), campaigns, campaignColumns)
	} else {
		campaigns, _, err := svc.Find(selector)
		if err != nil {
			return fmt.Errorf("finding campaigns: %w", err)
		}
		output.Print(getFormat(), campaigns, campaignColumns)
	}
	return nil
}

func runCampaignsCreate(cmd *cobra.Command, args []string) error {
	client, err := newAPIClient()
	if err != nil {
		return err
	}

	currency, err := resolveOrgCurrency(client)
	if err != nil {
		return err
	}

	if err := checkBudgetLimit(campDaily); err != nil {
		return err
	}

	campaign := &models.Campaign{
		Name:               campName,
		AdamID:             campAppID,
		Status:             campStatus,
		CountriesOrRegions: strings.Split(campCountries, ","),
		BudgetAmount:       &models.Money{Amount: campBudget, Currency: currency},
		DailyBudgetAmount:  &models.Money{Amount: campDaily, Currency: currency},
		AdChannelType:      "SEARCH",
		SupplySources:      []string{"APPSTORE_SEARCH_RESULTS"},
		BillingEvent:       "TAPS",
	}

	svc := services.NewCampaignService(client)
	created, err := svc.Create(campaign)
	if err != nil {
		return fmt.Errorf("creating campaign: %w", err)
	}

	output.Print(getFormat(), created, campaignColumns)
	return nil
}

func runCampaignsUpdate(cmd *cobra.Command, args []string) error {
	id, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid campaign ID: %s", args[0])
	}

	client, err := newAPIClient()
	if err != nil {
		return err
	}

	update := &models.CampaignUpdate{}
	hasUpdate := false

	if cmd.Flags().Changed("name") {
		update.Name = campName
		hasUpdate = true
	}
	if cmd.Flags().Changed("budget") || cmd.Flags().Changed("daily-budget") {
		currency, err := resolveOrgCurrency(client)
		if err != nil {
			return err
		}
		if cmd.Flags().Changed("budget") {
			update.BudgetAmount = &models.Money{Amount: campBudget, Currency: currency}
			hasUpdate = true
		}
		if cmd.Flags().Changed("daily-budget") {
			if err := checkBudgetLimit(campDaily); err != nil {
				return err
			}
			update.DailyBudgetAmount = &models.Money{Amount: campDaily, Currency: currency}
			hasUpdate = true
		}
	}
	if cmd.Flags().Changed("status") {
		update.Status = campStatus
		hasUpdate = true
	}

	if !hasUpdate {
		return fmt.Errorf("no update flags provided")
	}

	svc := services.NewCampaignService(client)
	updated, err := svc.Update(id, update)
	if err != nil {
		return fmt.Errorf("updating campaign: %w", err)
	}

	output.Print(getFormat(), updated, campaignColumns)
	return nil
}

func runCampaignsDelete(cmd *cobra.Command, args []string) error {
	id, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid campaign ID: %s", args[0])
	}

	client, err := newAPIClient()
	if err != nil {
		return err
	}

	svc := services.NewCampaignService(client)
	if err := svc.Delete(id); err != nil {
		return fmt.Errorf("deleting campaign: %w", err)
	}

	fmt.Printf("Campaign %d deleted.\n", id)
	return nil
}
