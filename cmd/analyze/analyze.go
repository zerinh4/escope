package analyze

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
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze [analyzer_name] [text]",
	Short: "Analyze text using specified analyzer or tokenizer",
	Long: `Analyze text using Elasticsearch analyzer or tokenizer.
	
Examples:
  # Analyze text with standard analyzer
  escope analyze standard "Hello World" --type analyzer
  
  # Analyze text with whitespace tokenizer
  escope analyze whitespace "Hello World" --type tokenizer
  
  # Default is analyzer type
  escope analyze standard "Hello World"`,
	SilenceErrors: true,
	Args:          cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		analyzerName := args[0]
		text := args[1]

		analyzeType, _ := cmd.Flags().GetString("type")
		if analyzeType == "" {
			analyzeType = "analyzer"
		}

		runAnalyze(analyzerName, text, analyzeType)
	},
}

func runAnalyze(analyzerName, text string, analyzeType string) {
	client := elastic.NewClientWrapper(connection.GetClient())
	analyzeService := services.NewAnalyzeService(client)

	result, err := util.ExecuteWithTimeout(func() (models.AnalyzeResult, error) {
		return analyzeService.AnalyzeText(context.Background(), analyzerName, text, analyzeType)
	})
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			fmt.Printf("Analyze failed: %s\n", constants.MsgTimeoutGeneric)
		} else {
			fmt.Printf("Analyze failed: %v\n", err)
		}
		return
	}

	// Display results using generic formatter
	if len(result.Tokens) == 0 {
		fmt.Println("No tokens generated")
		return
	}

	headers := []string{"Position", "Type", "Start", "End", "Token"}
	rows := make([][]string, 0, len(result.Tokens))

	for _, token := range result.Tokens {
		displayToken := token.Token
		if len(displayToken) > 40 {
			displayToken = token.Token[:37] + "..."
		}

		displayType := token.Type
		if len(displayType) > 15 {
			displayType = token.Type[:12] + "..."
		}

		row := []string{
			fmt.Sprintf("%d", token.Position),
			displayType,
			fmt.Sprintf("%d", token.StartOffset),
			fmt.Sprintf("%d", token.EndOffset),
			displayToken,
		}
		rows = append(rows, row)
	}

	formatter := ui.NewGenericTableFormatter()
	fmt.Print(formatter.FormatTable(headers, rows))
}

func init() {
	analyzeCmd.Flags().String("type", "analyzer", "Analyze type: 'analyzer' or 'tokenizer'")
	core.RootCmd.AddCommand(analyzeCmd)
}
