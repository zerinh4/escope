package termvectors

import (
	"context"
	"errors"
	"fmt"
	"github.com/mertbahardogan/escope/cmd/core"
	"github.com/mertbahardogan/escope/internal/connection"
	"github.com/mertbahardogan/escope/internal/constants"
	"github.com/mertbahardogan/escope/internal/elastic"
	"github.com/mertbahardogan/escope/internal/models"
	"github.com/mertbahardogan/escope/internal/services"
	"github.com/mertbahardogan/escope/internal/ui"
	"github.com/mertbahardogan/escope/internal/util"
	"github.com/spf13/cobra"
	"strings"
)

var termvectorsCmd = &cobra.Command{
	Use:   "termvectors [index] [document_id] [term]",
	Short: "Analyze term vectors and search terms in fields",
	Long: `Analyze term vectors for documents and search for specific terms across multiple fields.
	
Examples:
  # Get document term vectors (2 args: index, doc_id)
  escope termvectors index_name 12345 --fields field1,field2
  
  # Search for term in document fields (3 args: index, doc_id, term)
  escope termvectors index_name 12345 term_example --fields field1,field2`,
	SilenceErrors: true,
	Args:          cobra.RangeArgs(2, 3),
	Run: func(cmd *cobra.Command, args []string) {
		indexName := args[0]
		documentID := args[1]

		fieldsFlag, _ := cmd.Flags().GetString("fields")
		var fields []string
		if fieldsFlag != "" {
			fields = strings.Split(fieldsFlag, ",")
		} else {
			_ = fmt.Errorf("error parsing fields flag")
		}

		if len(args) == 3 {
			searchTerm := args[2]
			runDocumentTermSearch(indexName, documentID, fields, searchTerm)
		} else {
			runDocumentTermvectors(indexName, documentID, fields)
		}
	},
}

func runDocumentTermvectors(indexName, documentID string, fields []string) {
	client := elastic.NewClientWrapper(connection.GetClient())
	termvectorsService := services.NewTermvectorsService(client)

	termInfos, err := util.ExecuteWithTimeout(func() ([]models.TermInfo, error) {
		return termvectorsService.GetDocumentTermvectors(context.Background(), indexName, documentID, fields)
	})
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			fmt.Printf("Document termvectors failed: %s\n", constants.MsgTimeoutGeneric)
		} else {
			fmt.Printf("Document termvectors failed: %v\n", err)
		}
		return
	}

	formatter := ui.NewTermvectorsFormatter()
	fmt.Print(formatter.FormatSummary(termInfos))
	fmt.Print(formatter.FormatTermvectorsTable(termInfos))
}

func runDocumentTermSearch(indexName, documentID string, fields []string, searchTerm string) {
	client := elastic.NewClientWrapper(connection.GetClient())
	termvectorsService := services.NewTermvectorsService(client)

	termInfos, err := util.ExecuteWithTimeout(func() ([]models.TermInfo, error) {
		return termvectorsService.GetDocumentTermvectors(context.Background(), indexName, documentID, fields)
	})
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			fmt.Printf("Document term search failed: %s\n", constants.MsgTimeoutGeneric)
		} else {
			fmt.Printf("Document term search failed: %v\n", err)
		}
		return
	}

	formatter := ui.NewTermvectorsFormatter()
	fmt.Print(formatter.FormatTermSearchResult(termInfos, searchTerm))
}

func init() {
	termvectorsCmd.Flags().String("fields", "", "Fields to analyze (comma-separated)")
	core.RootCmd.AddCommand(termvectorsCmd)
}
