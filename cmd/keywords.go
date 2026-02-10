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

var keywordsCmd = &cobra.Command{
	Use:   "keywords",
	Short: "Manage targeting keywords",
}

var kwListCmd = &cobra.Command{
	Use:   "list",
	Short: "List targeting keywords",
	RunE:  runKWList,
}

var kwGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get a keyword by ID",
	Args:  cobra.ExactArgs(1),
	RunE:  runKWGet,
}

var kwFindCmd = &cobra.Command{
	Use:   "find",
	Short: "Find keywords with filters",
	RunE:  runKWFind,
}

var kwCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create targeting keywords (supports bulk)",
	RunE:  runKWCreate,
}

var kwUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update targeting keywords",
	RunE:  runKWUpdate,
}

var kwDeleteCmd = &cobra.Command{
	Use:   "delete <id,...>",
	Short: "Delete targeting keywords",
	Args:  cobra.ExactArgs(1),
	RunE:  runKWDelete,
}

var (
	kwCampaignID int64
	kwAdGroupID  int64
	kwLimit      int
	kwOffset     int
	kwFilters    []string
	kwSorts      []string
	kwAll        bool
	kwTexts      []string
	kwMatchType  string
	kwBid        string
	kwStatus     string
	kwID         int64
)

func init() {
	// Common flags
	for _, cmd := range []*cobra.Command{kwListCmd, kwGetCmd, kwFindCmd, kwCreateCmd, kwUpdateCmd, kwDeleteCmd} {
		cmd.Flags().Int64Var(&kwCampaignID, "campaign-id", 0, "Campaign ID (required)")
		cmd.Flags().Int64Var(&kwAdGroupID, "adgroup-id", 0, "Ad group ID (required)")
		cmd.MarkFlagRequired("campaign-id")
		cmd.MarkFlagRequired("adgroup-id")
	}

	// list
	kwListCmd.Flags().IntVar(&kwLimit, "limit", 20, "Number of results")
	kwListCmd.Flags().IntVar(&kwOffset, "offset", 0, "Results offset")

	// find
	kwFindCmd.Flags().StringSliceVar(&kwFilters, "filter", nil, "Filter conditions")
	kwFindCmd.Flags().StringSliceVar(&kwSorts, "sort", nil, "Sort order")
	kwFindCmd.Flags().IntVar(&kwLimit, "limit", 20, "Number of results")
	kwFindCmd.Flags().IntVar(&kwOffset, "offset", 0, "Results offset")
	kwFindCmd.Flags().BoolVar(&kwAll, "all", false, "Fetch all pages")

	// create
	kwCreateCmd.Flags().StringSliceVar(&kwTexts, "text", nil, "Keyword text(s) — repeatable for bulk")
	kwCreateCmd.Flags().StringVar(&kwMatchType, "match-type", "BROAD", "Match type: BROAD or EXACT")
	kwCreateCmd.Flags().StringVar(&kwBid, "bid", "", "Bid amount (e.g. 1.50)")
	kwCreateCmd.MarkFlagRequired("text")

	// update
	kwUpdateCmd.Flags().Int64Var(&kwID, "id", 0, "Keyword ID to update (required)")
	kwUpdateCmd.Flags().StringVar(&kwStatus, "status", "", "Status (ACTIVE/PAUSED)")
	kwUpdateCmd.Flags().StringVar(&kwBid, "bid", "", "Bid amount")
	kwUpdateCmd.MarkFlagRequired("id")

	keywordsCmd.AddCommand(kwListCmd, kwGetCmd, kwFindCmd, kwCreateCmd, kwUpdateCmd, kwDeleteCmd)
	rootCmd.AddCommand(keywordsCmd)
}

var keywordColumns = []output.Column{
	{Header: "ID", Field: "ID", Width: 12},
	{Header: "TEXT", Field: "Text", Width: 30},
	{Header: "MATCH TYPE", Field: "MatchType", Width: 12},
	{Header: "STATUS", Field: "Status", Width: 10},
	{Header: "BID", Field: "BidAmount", Width: 12},
}

func runKWList(cmd *cobra.Command, args []string) error {
	client, err := newAPIClient()
	if err != nil {
		return err
	}

	svc := services.NewKeywordService(client)
	keywords, _, err := svc.List(kwCampaignID, kwAdGroupID, kwLimit, kwOffset)
	if err != nil {
		return fmt.Errorf("listing keywords: %w", err)
	}

	output.Print(getFormat(), keywords, keywordColumns)
	return nil
}

func runKWGet(cmd *cobra.Command, args []string) error {
	id, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid keyword ID: %s", args[0])
	}

	client, err := newAPIClient()
	if err != nil {
		return err
	}

	svc := services.NewKeywordService(client)
	keyword, err := svc.Get(kwCampaignID, kwAdGroupID, id)
	if err != nil {
		return fmt.Errorf("getting keyword: %w", err)
	}

	output.Print(getFormat(), keyword, keywordColumns)
	return nil
}

func runKWFind(cmd *cobra.Command, args []string) error {
	client, err := newAPIClient()
	if err != nil {
		return err
	}

	selector := models.NewSelector(kwLimit, kwOffset)
	selector.Conditions = parseFilters(kwFilters)
	selector.OrderBy = parseSorts(kwSorts)

	svc := services.NewKeywordService(client)

	if kwAll {
		keywords, err := svc.FindAll(kwCampaignID, kwAdGroupID, selector)
		if err != nil {
			return fmt.Errorf("finding keywords: %w", err)
		}
		output.Print(getFormat(), keywords, keywordColumns)
	} else {
		keywords, _, err := svc.Find(kwCampaignID, kwAdGroupID, selector)
		if err != nil {
			return fmt.Errorf("finding keywords: %w", err)
		}
		output.Print(getFormat(), keywords, keywordColumns)
	}
	return nil
}

func runKWCreate(cmd *cobra.Command, args []string) error {
	client, err := newAPIClient()
	if err != nil {
		return err
	}

	currency, err := resolveOrgCurrency(client)
	if err != nil {
		return err
	}

	if kwBid != "" {
		if err := checkBidLimit(kwBid); err != nil {
			return err
		}
	}

	var keywords []models.Keyword
	for _, text := range kwTexts {
		kw := models.Keyword{
			Text:      text,
			MatchType: kwMatchType,
		}
		if kwBid != "" {
			kw.BidAmount = &models.Money{Amount: kwBid, Currency: currency}
		}
		keywords = append(keywords, kw)
	}

	svc := services.NewKeywordService(client)
	created, err := svc.Create(kwCampaignID, kwAdGroupID, keywords)
	if err != nil {
		return fmt.Errorf("creating keywords: %w", err)
	}

	output.Print(getFormat(), created, keywordColumns)
	return nil
}

func runKWUpdate(cmd *cobra.Command, args []string) error {
	client, err := newAPIClient()
	if err != nil {
		return err
	}

	update := models.KeywordUpdate{ID: kwID}
	if cmd.Flags().Changed("status") {
		update.Status = kwStatus
	}
	if cmd.Flags().Changed("bid") {
		if err := checkBidLimit(kwBid); err != nil {
			return err
		}
		currency, err := resolveOrgCurrency(client)
		if err != nil {
			return err
		}
		update.BidAmount = &models.Money{Amount: kwBid, Currency: currency}
	}

	svc := services.NewKeywordService(client)
	updated, err := svc.Update(kwCampaignID, kwAdGroupID, []models.KeywordUpdate{update})
	if err != nil {
		return fmt.Errorf("updating keyword: %w", err)
	}

	output.Print(getFormat(), updated, keywordColumns)
	return nil
}

func runKWDelete(cmd *cobra.Command, args []string) error {
	client, err := newAPIClient()
	if err != nil {
		return err
	}

	var ids []int64
	for _, s := range strings.Split(args[0], ",") {
		id, err := strconv.ParseInt(strings.TrimSpace(s), 10, 64)
		if err != nil {
			return fmt.Errorf("invalid keyword ID: %s", s)
		}
		ids = append(ids, id)
	}

	svc := services.NewKeywordService(client)
	if err := svc.Delete(kwCampaignID, kwAdGroupID, ids); err != nil {
		return fmt.Errorf("deleting keywords: %w", err)
	}

	fmt.Printf("Deleted %d keyword(s).\n", len(ids))
	return nil
}
