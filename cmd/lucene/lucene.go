package lucene

import (
	"context"
	"errors"
	"escope/cmd/core"
	"escope/internal/connection"
	"escope/internal/constants"
	"escope/internal/elastic"
	"escope/internal/models"
	"escope/internal/services"
	"escope/internal/ui"
	"escope/internal/util"
	"fmt"
	"github.com/spf13/cobra"
	"sort"
)

var (
	indexName string
)

var luceneCmd = &cobra.Command{
	Use:           "lucene",
	Short:         "Show detailed Lucene segment and inverted index information",
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		client := elastic.NewClientWrapper(connection.GetClient())
		luceneService := services.NewLuceneService(client)

		luceneStats, err := util.ExecuteWithTimeout(func() ([]models.LuceneStats, error) {
			return luceneService.GetLuceneStats(context.Background())
		})
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				fmt.Printf("Lucene stats fetch failed: %s\n", constants.MsgTimeoutGeneric)
			} else {
				fmt.Printf("Lucene stats fetch failed: %v\n", err)
			}
			return
		}

		var filteredStats []models.LuceneStats
		if indexName != "" {
			for _, stat := range luceneStats {
				if stat.IndexName == indexName {
					filteredStats = append(filteredStats, stat)
				}
			}
			if len(filteredStats) == 0 {
				fmt.Printf("Index '%s' not found\n", indexName)
				return
			}
		} else {
			// Otherwise filter out system indices
			for _, stat := range luceneStats {
				if !util.IsSystemIndex(stat.IndexName) {
					filteredStats = append(filteredStats, stat)
				}
			}
		}

		sort.Slice(filteredStats, func(i, j int) bool {
			return filteredStats[i].SegmentMemoryBytes > filteredStats[j].SegmentMemoryBytes
		})

		headers := []string{"Segments", "Total Memory", "Terms Memory", "Stored Memory", "DocValues", "Index"}
		rows := make([][]string, 0, len(filteredStats))

		for _, stat := range filteredStats {
			row := []string{
				fmt.Sprintf("%d", stat.SegmentCount),
				stat.SegmentMemory,
				stat.TermsMemory,
				stat.StoredMemory,
				stat.DocValuesMemory,
				stat.IndexName,
			}
			rows = append(rows, row)
		}

		formatter := ui.NewGenericTableFormatter()
		fmt.Print(formatter.FormatTable(headers, rows))

		if indexName != "" && len(filteredStats) > 0 {
			for _, stat := range filteredStats {
				fmt.Printf("\n# Index: %s\n", stat.IndexName)
				fmt.Printf("   Segments: %d\n", stat.SegmentCount)
				fmt.Printf("   Total Memory: %s\n", stat.SegmentMemory)
				fmt.Printf("   Index Memory: %s\n", stat.IndexMemory)

				fmt.Println("   Memory Breakdown:")
				fmt.Printf("     • Terms (Inverted Index): %s\n", stat.TermsMemory)
				fmt.Printf("     • Stored Fields: %s\n", stat.StoredMemory)
				fmt.Printf("     • DocValues: %s\n", stat.DocValuesMemory)
				fmt.Printf("     • Points (Numeric): %s\n", stat.PointsMemory)
				fmt.Printf("     • Norms: %s\n", stat.NormsMemory)
				fmt.Printf("     • Fixed BitSet: %s\n", stat.FixedBitSetMemory)
				fmt.Printf("     • Version Map: %s\n", stat.VersionMapMemory)

				if stat.MaxUnsafeAutoIDTimestamp > 0 {
					fmt.Printf("   Max Auto-ID Timestamp: %d\n", stat.MaxUnsafeAutoIDTimestamp)
				}
			}
		}
	},
}

func init() {
	core.RootCmd.AddCommand(luceneCmd)
	luceneCmd.Flags().StringVarP(&indexName, "name", "n", "", "Show detailed breakdown for specific index")
}
