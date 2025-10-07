package segments

import (
	"context"
	"errors"
	"fmt"
	"github.com/mertbahardogan/escope/cmd/core"
	"github.com/mertbahardogan/escope/cmd/util"
	"github.com/mertbahardogan/escope/internal/connection"
	"github.com/mertbahardogan/escope/internal/constants"
	"github.com/mertbahardogan/escope/internal/elastic"
	"github.com/mertbahardogan/escope/internal/models"
	"github.com/mertbahardogan/escope/internal/services"
	"github.com/mertbahardogan/escope/internal/ui"
	internalUtil "github.com/mertbahardogan/escope/internal/util"
	"github.com/spf13/cobra"
	"sort"
)

var segmentsCmd = &cobra.Command{
	Use:           "segments",
	Short:         "Show segment analysis and optimization recommendations",
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		client := elastic.NewClientWrapper(connection.GetClient())
		segmentsService := services.NewSegmentsService(client)

		segments, err := internalUtil.ExecuteWithTimeout(func() ([]models.SegmentInfo, error) {
			return segmentsService.GetSegmentsInfo(context.Background())
		})
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				fmt.Printf("Segments info fetch failed: %s\n", constants.MsgTimeoutGeneric)
			} else {
				fmt.Printf("Segments info fetch failed: %v\n", err)
			}
			return
		}

		var filteredSegments []models.SegmentInfo
		for _, seg := range segments {
			if !util.IsSystemIndex(seg.Index) {
				filteredSegments = append(filteredSegments, seg)
			}
		}

		if len(filteredSegments) == 0 {
			fmt.Println("No indices found with segments")
			return
		}

		headers := []string{"Segments", "Memory", "Avg Mem/Seg", "Index"}
		rows := make([][]string, 0, len(filteredSegments))

		var indexStats []struct {
			Name         string
			SegmentCount int
			TotalMemory  string
			AvgMemPerSeg string
		}

		for _, seg := range filteredSegments {
			avgMemPerSeg := "0b"
			if seg.SegmentCount > 0 {
				avgMemPerSeg = util.FormatBytes(seg.SizeBytes / int64(seg.SegmentCount))
			}

			indexStats = append(indexStats, struct {
				Name         string
				SegmentCount int
				TotalMemory  string
				AvgMemPerSeg string
			}{
				Name:         seg.Index,
				SegmentCount: seg.SegmentCount,
				TotalMemory:  util.FormatBytes(seg.SizeBytes),
				AvgMemPerSeg: avgMemPerSeg,
			})
		}

		sort.Slice(indexStats, func(i, j int) bool {
			return indexStats[i].SegmentCount > indexStats[j].SegmentCount
		})

		for _, stat := range indexStats {
			row := []string{
				fmt.Sprintf("%d", stat.SegmentCount),
				stat.TotalMemory,
				stat.AvgMemPerSeg,
				stat.Name,
			}
			rows = append(rows, row)
		}

		formatter := ui.NewGenericTableFormatter()
		fmt.Print(formatter.FormatTable(headers, rows))
	},
}

func init() {
	core.RootCmd.AddCommand(segmentsCmd)
}
