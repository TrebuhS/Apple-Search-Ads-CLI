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

var negKeywordsCmd = &cobra.Command{
	Use:   "negative-keywords",
	Short: "Manage negative keywords (campaign and ad-group level)",
}

// --- Campaign-level negative keywords ---

var nkCampaignListCmd = &cobra.Command{
	Use:   "campaign-list",
	Short: "List campaign-level negative keywords",
	RunE:  runNKCampaignList,
}

var nkCampaignCreateCmd = &cobra.Command{
	Use:   "campaign-create",
	Short: "Create campaign-level negative keywords",
	RunE:  runNKCampaignCreate,
}

var nkCampaignFindCmd = &cobra.Command{
	Use:   "campaign-find",
	Short: "Find campaign-level negative keywords",
	RunE:  runNKCampaignFind,
}

var nkCampaignDeleteCmd = &cobra.Command{
	Use:   "campaign-delete <id,...>",
	Short: "Delete campaign-level negative keywords",
	Args:  cobra.ExactArgs(1),
	RunE:  runNKCampaignDelete,
}

// --- Ad group-level negative keywords ---

var nkAdGroupListCmd = &cobra.Command{
	Use:   "adgroup-list",
	Short: "List ad-group-level negative keywords",
	RunE:  runNKAdGroupList,
}

var nkAdGroupCreateCmd = &cobra.Command{
	Use:   "adgroup-create",
	Short: "Create ad-group-level negative keywords",
	RunE:  runNKAdGroupCreate,
}

var nkAdGroupFindCmd = &cobra.Command{
	Use:   "adgroup-find",
	Short: "Find ad-group-level negative keywords",
	RunE:  runNKAdGroupFind,
}

var nkAdGroupDeleteCmd = &cobra.Command{
	Use:   "adgroup-delete <id,...>",
	Short: "Delete ad-group-level negative keywords",
	Args:  cobra.ExactArgs(1),
	RunE:  runNKAdGroupDelete,
}

var (
	nkCampaignID int64
	nkAdGroupID  int64
	nkLimit      int
	nkOffset     int
	nkTexts      []string
	nkMatchType  string
	nkFilters    []string
	nkSorts      []string
)

func init() {
	// Campaign-level commands
	for _, cmd := range []*cobra.Command{nkCampaignListCmd, nkCampaignCreateCmd, nkCampaignFindCmd, nkCampaignDeleteCmd} {
		cmd.Flags().Int64Var(&nkCampaignID, "campaign-id", 0, "Campaign ID (required)")
		cmd.MarkFlagRequired("campaign-id")
	}

	nkCampaignListCmd.Flags().IntVar(&nkLimit, "limit", 20, "Number of results")
	nkCampaignListCmd.Flags().IntVar(&nkOffset, "offset", 0, "Results offset")

	nkCampaignCreateCmd.Flags().StringSliceVar(&nkTexts, "text", nil, "Keyword text(s)")
	nkCampaignCreateCmd.Flags().StringVar(&nkMatchType, "match-type", "EXACT", "Match type: BROAD or EXACT")
	nkCampaignCreateCmd.MarkFlagRequired("text")

	nkCampaignFindCmd.Flags().StringSliceVar(&nkFilters, "filter", nil, "Filter conditions")
	nkCampaignFindCmd.Flags().StringSliceVar(&nkSorts, "sort", nil, "Sort order")
	nkCampaignFindCmd.Flags().IntVar(&nkLimit, "limit", 20, "Number of results")
	nkCampaignFindCmd.Flags().IntVar(&nkOffset, "offset", 0, "Results offset")

	// Ad group-level commands
	for _, cmd := range []*cobra.Command{nkAdGroupListCmd, nkAdGroupCreateCmd, nkAdGroupFindCmd, nkAdGroupDeleteCmd} {
		cmd.Flags().Int64Var(&nkCampaignID, "campaign-id", 0, "Campaign ID (required)")
		cmd.Flags().Int64Var(&nkAdGroupID, "adgroup-id", 0, "Ad group ID (required)")
		cmd.MarkFlagRequired("campaign-id")
		cmd.MarkFlagRequired("adgroup-id")
	}

	nkAdGroupListCmd.Flags().IntVar(&nkLimit, "limit", 20, "Number of results")
	nkAdGroupListCmd.Flags().IntVar(&nkOffset, "offset", 0, "Results offset")

	nkAdGroupCreateCmd.Flags().StringSliceVar(&nkTexts, "text", nil, "Keyword text(s)")
	nkAdGroupCreateCmd.Flags().StringVar(&nkMatchType, "match-type", "EXACT", "Match type: BROAD or EXACT")
	nkAdGroupCreateCmd.MarkFlagRequired("text")

	nkAdGroupFindCmd.Flags().StringSliceVar(&nkFilters, "filter", nil, "Filter conditions")
	nkAdGroupFindCmd.Flags().StringSliceVar(&nkSorts, "sort", nil, "Sort order")
	nkAdGroupFindCmd.Flags().IntVar(&nkLimit, "limit", 20, "Number of results")
	nkAdGroupFindCmd.Flags().IntVar(&nkOffset, "offset", 0, "Results offset")

	negKeywordsCmd.AddCommand(
		nkCampaignListCmd, nkCampaignCreateCmd, nkCampaignFindCmd, nkCampaignDeleteCmd,
		nkAdGroupListCmd, nkAdGroupCreateCmd, nkAdGroupFindCmd, nkAdGroupDeleteCmd,
	)
	rootCmd.AddCommand(negKeywordsCmd)
}

var negKeywordColumns = []output.Column{
	{Header: "ID", Field: "ID", Width: 12},
	{Header: "TEXT", Field: "Text", Width: 30},
	{Header: "MATCH TYPE", Field: "MatchType", Width: 12},
	{Header: "STATUS", Field: "Status", Width: 10},
}

// --- Campaign-level implementations ---

func runNKCampaignList(cmd *cobra.Command, args []string) error {
	client, err := newAPIClient()
	if err != nil {
		return err
	}

	svc := services.NewKeywordService(client)
	keywords, _, err := svc.ListCampaignNegativeKeywords(nkCampaignID, nkLimit, nkOffset)
	if err != nil {
		return fmt.Errorf("listing negative keywords: %w", err)
	}

	output.Print(getFormat(), keywords, negKeywordColumns)
	return nil
}

func runNKCampaignCreate(cmd *cobra.Command, args []string) error {
	client, err := newAPIClient()
	if err != nil {
		return err
	}

	var keywords []models.NegativeKeyword
	for _, text := range nkTexts {
		keywords = append(keywords, models.NegativeKeyword{
			Text:      text,
			MatchType: nkMatchType,
		})
	}

	svc := services.NewKeywordService(client)
	created, err := svc.CreateCampaignNegativeKeywords(nkCampaignID, keywords)
	if err != nil {
		return fmt.Errorf("creating negative keywords: %w", err)
	}

	output.Print(getFormat(), created, negKeywordColumns)
	return nil
}

func runNKCampaignFind(cmd *cobra.Command, args []string) error {
	client, err := newAPIClient()
	if err != nil {
		return err
	}

	selector := models.NewSelector(nkLimit, nkOffset)
	selector.Conditions = parseFilters(nkFilters)
	selector.OrderBy = parseSorts(nkSorts)

	svc := services.NewKeywordService(client)
	keywords, _, err := svc.FindCampaignNegativeKeywords(nkCampaignID, selector)
	if err != nil {
		return fmt.Errorf("finding negative keywords: %w", err)
	}

	output.Print(getFormat(), keywords, negKeywordColumns)
	return nil
}

func runNKCampaignDelete(cmd *cobra.Command, args []string) error {
	client, err := newAPIClient()
	if err != nil {
		return err
	}

	ids, err := parseIDList(args[0])
	if err != nil {
		return err
	}

	svc := services.NewKeywordService(client)
	if err := svc.DeleteCampaignNegativeKeywords(nkCampaignID, ids); err != nil {
		return fmt.Errorf("deleting negative keywords: %w", err)
	}

	fmt.Printf("Deleted %d negative keyword(s).\n", len(ids))
	return nil
}

// --- Ad group-level implementations ---

func runNKAdGroupList(cmd *cobra.Command, args []string) error {
	client, err := newAPIClient()
	if err != nil {
		return err
	}

	svc := services.NewKeywordService(client)
	keywords, _, err := svc.ListAdGroupNegativeKeywords(nkCampaignID, nkAdGroupID, nkLimit, nkOffset)
	if err != nil {
		return fmt.Errorf("listing negative keywords: %w", err)
	}

	output.Print(getFormat(), keywords, negKeywordColumns)
	return nil
}

func runNKAdGroupCreate(cmd *cobra.Command, args []string) error {
	client, err := newAPIClient()
	if err != nil {
		return err
	}

	var keywords []models.NegativeKeyword
	for _, text := range nkTexts {
		keywords = append(keywords, models.NegativeKeyword{
			Text:      text,
			MatchType: nkMatchType,
		})
	}

	svc := services.NewKeywordService(client)
	created, err := svc.CreateAdGroupNegativeKeywords(nkCampaignID, nkAdGroupID, keywords)
	if err != nil {
		return fmt.Errorf("creating negative keywords: %w", err)
	}

	output.Print(getFormat(), created, negKeywordColumns)
	return nil
}

func runNKAdGroupFind(cmd *cobra.Command, args []string) error {
	client, err := newAPIClient()
	if err != nil {
		return err
	}

	selector := models.NewSelector(nkLimit, nkOffset)
	selector.Conditions = parseFilters(nkFilters)
	selector.OrderBy = parseSorts(nkSorts)

	svc := services.NewKeywordService(client)
	keywords, _, err := svc.FindAdGroupNegativeKeywords(nkCampaignID, nkAdGroupID, selector)
	if err != nil {
		return fmt.Errorf("finding negative keywords: %w", err)
	}

	output.Print(getFormat(), keywords, negKeywordColumns)
	return nil
}

func runNKAdGroupDelete(cmd *cobra.Command, args []string) error {
	client, err := newAPIClient()
	if err != nil {
		return err
	}

	ids, err := parseIDList(args[0])
	if err != nil {
		return err
	}

	svc := services.NewKeywordService(client)
	if err := svc.DeleteAdGroupNegativeKeywords(nkCampaignID, nkAdGroupID, ids); err != nil {
		return fmt.Errorf("deleting negative keywords: %w", err)
	}

	fmt.Printf("Deleted %d negative keyword(s).\n", len(ids))
	return nil
}

func parseIDList(s string) ([]int64, error) {
	var ids []int64
	for _, part := range strings.Split(s, ",") {
		id, err := strconv.ParseInt(strings.TrimSpace(part), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid ID: %s", part)
		}
		ids = append(ids, id)
	}
	return ids, nil
}
