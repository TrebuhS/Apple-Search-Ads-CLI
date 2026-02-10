package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/trebuhs/asa-cli/internal/models"
	"github.com/trebuhs/asa-cli/internal/output"
	"github.com/trebuhs/asa-cli/internal/services"
)

var adgroupsCmd = &cobra.Command{
	Use:   "adgroups",
	Short: "Manage ad groups",
}

var adgroupsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List ad groups for a campaign",
	RunE:  runAdGroupsList,
}

var adgroupsGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get an ad group by ID",
	Args:  cobra.ExactArgs(1),
	RunE:  runAdGroupsGet,
}

var adgroupsFindCmd = &cobra.Command{
	Use:   "find",
	Short: "Find ad groups with filters",
	RunE:  runAdGroupsFind,
}

var adgroupsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new ad group",
	RunE:  runAdGroupsCreate,
}

var adgroupsUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update an ad group",
	Args:  cobra.ExactArgs(1),
	RunE:  runAdGroupsUpdate,
}

var adgroupsDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete an ad group",
	Args:  cobra.ExactArgs(1),
	RunE:  runAdGroupsDelete,
}

var (
	agCampaignID int64
	agLimit      int
	agOffset     int
	agFilters    []string
	agSorts      []string
	agAll        bool
	agName       string
	agBid        string
	agCpaGoal    string
	agStatus     string
	agAutoKW     string
	agStartTime  string
	agEndTime    string
)

func init() {
	// Common campaign-id flag
	for _, cmd := range []*cobra.Command{adgroupsListCmd, adgroupsGetCmd, adgroupsFindCmd, adgroupsCreateCmd, adgroupsUpdateCmd, adgroupsDeleteCmd} {
		cmd.Flags().Int64Var(&agCampaignID, "campaign-id", 0, "Campaign ID (required)")
		cmd.MarkFlagRequired("campaign-id")
	}

	// list
	adgroupsListCmd.Flags().IntVar(&agLimit, "limit", 20, "Number of results")
	adgroupsListCmd.Flags().IntVar(&agOffset, "offset", 0, "Results offset")

	// find
	adgroupsFindCmd.Flags().StringSliceVar(&agFilters, "filter", nil, `Filter conditions`)
	adgroupsFindCmd.Flags().StringSliceVar(&agSorts, "sort", nil, `Sort order`)
	adgroupsFindCmd.Flags().IntVar(&agLimit, "limit", 20, "Number of results")
	adgroupsFindCmd.Flags().IntVar(&agOffset, "offset", 0, "Results offset")
	adgroupsFindCmd.Flags().BoolVar(&agAll, "all", false, "Fetch all pages")

	// create
	adgroupsCreateCmd.Flags().StringVar(&agName, "name", "", "Ad group name (required)")
	adgroupsCreateCmd.Flags().StringVar(&agBid, "default-bid", "", "Default bid amount (e.g. 1.50)")
	adgroupsCreateCmd.Flags().StringVar(&agCpaGoal, "cpa-goal", "", "CPA goal amount")
	adgroupsCreateCmd.Flags().StringVar(&agStatus, "status", "ENABLED", "Status")
	adgroupsCreateCmd.Flags().StringVar(&agAutoKW, "auto-keywords", "true", "Automated keywords opt-in (true/false)")
	adgroupsCreateCmd.Flags().StringVar(&agStartTime, "start-time", "", "Start time (ISO 8601)")
	adgroupsCreateCmd.Flags().StringVar(&agEndTime, "end-time", "", "End time (ISO 8601)")
	adgroupsCreateCmd.MarkFlagRequired("name")
	adgroupsCreateCmd.MarkFlagRequired("default-bid")

	// update
	adgroupsUpdateCmd.Flags().StringVar(&agName, "name", "", "Ad group name")
	adgroupsUpdateCmd.Flags().StringVar(&agBid, "default-bid", "", "Default bid amount")
	adgroupsUpdateCmd.Flags().StringVar(&agCpaGoal, "cpa-goal", "", "CPA goal amount")
	adgroupsUpdateCmd.Flags().StringVar(&agStatus, "status", "", "Status (ENABLED/PAUSED)")
	adgroupsUpdateCmd.Flags().StringVar(&agAutoKW, "auto-keywords", "", "Automated keywords (true/false)")
	adgroupsUpdateCmd.Flags().StringVar(&agStartTime, "start-time", "", "Start time")
	adgroupsUpdateCmd.Flags().StringVar(&agEndTime, "end-time", "", "End time")

	adgroupsCmd.AddCommand(adgroupsListCmd, adgroupsGetCmd, adgroupsFindCmd, adgroupsCreateCmd, adgroupsUpdateCmd, adgroupsDeleteCmd)
	rootCmd.AddCommand(adgroupsCmd)
}

var adgroupColumns = []output.Column{
	{Header: "ID", Field: "ID", Width: 12},
	{Header: "NAME", Field: "Name", Width: 25},
	{Header: "STATUS", Field: "Status", Width: 10},
	{Header: "SERVING", Field: "ServingStatus", Width: 12},
	{Header: "DEFAULT BID", Field: "DefaultBidAmount", Width: 15},
	{Header: "CPA GOAL", Field: "CpaGoal", Width: 12},
}

func runAdGroupsList(cmd *cobra.Command, args []string) error {
	client, err := newAPIClient()
	if err != nil {
		return err
	}

	svc := services.NewAdGroupService(client)
	adgroups, _, err := svc.List(agCampaignID, agLimit, agOffset)
	if err != nil {
		return fmt.Errorf("listing ad groups: %w", err)
	}

	output.Print(getFormat(), adgroups, adgroupColumns)
	return nil
}

func runAdGroupsGet(cmd *cobra.Command, args []string) error {
	id, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid ad group ID: %s", args[0])
	}

	client, err := newAPIClient()
	if err != nil {
		return err
	}

	svc := services.NewAdGroupService(client)
	adgroup, err := svc.Get(agCampaignID, id)
	if err != nil {
		return fmt.Errorf("getting ad group: %w", err)
	}

	output.Print(getFormat(), adgroup, adgroupColumns)
	return nil
}

func runAdGroupsFind(cmd *cobra.Command, args []string) error {
	client, err := newAPIClient()
	if err != nil {
		return err
	}

	selector := models.NewSelector(agLimit, agOffset)
	selector.Conditions = parseFilters(agFilters)
	selector.OrderBy = parseSorts(agSorts)

	svc := services.NewAdGroupService(client)

	if agAll {
		adgroups, err := svc.FindAll(agCampaignID, selector)
		if err != nil {
			return fmt.Errorf("finding ad groups: %w", err)
		}
		output.Print(getFormat(), adgroups, adgroupColumns)
	} else {
		adgroups, _, err := svc.Find(agCampaignID, selector)
		if err != nil {
			return fmt.Errorf("finding ad groups: %w", err)
		}
		output.Print(getFormat(), adgroups, adgroupColumns)
	}
	return nil
}

func runAdGroupsCreate(cmd *cobra.Command, args []string) error {
	client, err := newAPIClient()
	if err != nil {
		return err
	}

	autoKW := agAutoKW == "true"
	adgroup := &models.AdGroup{
		Name:                   agName,
		Status:                 agStatus,
		DefaultBidAmount:       &models.Money{Amount: agBid, Currency: "USD"},
		AutomatedKeywordsOptIn: autoKW,
	}

	if agCpaGoal != "" {
		adgroup.CpaGoal = &models.Money{Amount: agCpaGoal, Currency: "USD"}
	}
	if agStartTime != "" {
		adgroup.StartTime = agStartTime
	}
	if agEndTime != "" {
		adgroup.EndTime = agEndTime
	}

	svc := services.NewAdGroupService(client)
	created, err := svc.Create(agCampaignID, adgroup)
	if err != nil {
		return fmt.Errorf("creating ad group: %w", err)
	}

	output.Print(getFormat(), created, adgroupColumns)
	return nil
}

func runAdGroupsUpdate(cmd *cobra.Command, args []string) error {
	id, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid ad group ID: %s", args[0])
	}

	client, err := newAPIClient()
	if err != nil {
		return err
	}

	update := &models.AdGroupUpdate{}
	hasUpdate := false

	if cmd.Flags().Changed("name") {
		update.Name = agName
		hasUpdate = true
	}
	if cmd.Flags().Changed("default-bid") {
		update.DefaultBidAmount = &models.Money{Amount: agBid, Currency: "USD"}
		hasUpdate = true
	}
	if cmd.Flags().Changed("cpa-goal") {
		update.CpaGoal = &models.Money{Amount: agCpaGoal, Currency: "USD"}
		hasUpdate = true
	}
	if cmd.Flags().Changed("status") {
		update.Status = agStatus
		hasUpdate = true
	}
	if cmd.Flags().Changed("auto-keywords") {
		val := agAutoKW == "true"
		update.AutomatedKeywordsOptIn = &val
		hasUpdate = true
	}
	if cmd.Flags().Changed("start-time") {
		update.StartTime = agStartTime
		hasUpdate = true
	}
	if cmd.Flags().Changed("end-time") {
		update.EndTime = agEndTime
		hasUpdate = true
	}

	if !hasUpdate {
		return fmt.Errorf("no update flags provided")
	}

	svc := services.NewAdGroupService(client)
	updated, err := svc.Update(agCampaignID, id, update)
	if err != nil {
		return fmt.Errorf("updating ad group: %w", err)
	}

	output.Print(getFormat(), updated, adgroupColumns)
	return nil
}

func runAdGroupsDelete(cmd *cobra.Command, args []string) error {
	id, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid ad group ID: %s", args[0])
	}

	client, err := newAPIClient()
	if err != nil {
		return err
	}

	svc := services.NewAdGroupService(client)
	if err := svc.Delete(agCampaignID, id); err != nil {
		return fmt.Errorf("deleting ad group: %w", err)
	}

	fmt.Printf("Ad group %d deleted.\n", id)
	return nil
}
