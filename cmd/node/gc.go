package node

import (
	"context"
	"errors"
	"escope/internal/connection"
	"escope/internal/constants"
	"escope/internal/elastic"
	"escope/internal/models"
	"escope/internal/services"
	"escope/internal/ui"
	"escope/internal/util"
	"fmt"
	"github.com/spf13/cobra"
)

var nodeGCCmd = &cobra.Command{
	Use:           "gc",
	Short:         "Show garbage collection information for nodes",
	SilenceErrors: true,
	Long: `Show garbage collection information for Elasticsearch nodes.
	
Examples:
  escope node gc                    # Show GC info for all nodes (table format)
  escope node gc --name=node1      # Show detailed GC info for specific node`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		client := elastic.NewClientWrapper(connection.GetClient())
		gcService := services.NewGCService(client)
		formatter := ui.NewGCFormatter()

		nodeName, _ := cmd.Flags().GetString("name")

		if nodeName != "" {
			gcInfo, err := util.ExecuteWithTimeout(func() (*models.GCInfo, error) {
				return gcService.GetGCInfoForNode(ctx, nodeName)
			})
			if err != nil {
				if errors.Is(err, context.DeadlineExceeded) {
					fmt.Printf("Failed to get GC info for node %s: %s\n", nodeName, constants.MsgTimeoutGeneric)
				} else {
					fmt.Printf("Failed to get GC info for node %s: %v\n", nodeName, err)
				}
				return
			}
			output := formatter.FormatGCDetails(*gcInfo)
			fmt.Print(output)
		} else {
			gcInfos, err := util.ExecuteWithTimeout(func() ([]models.GCInfo, error) {
				return gcService.GetGCInfo(ctx)
			})
			if err != nil {
				if errors.Is(err, context.DeadlineExceeded) {
					fmt.Printf("Failed to get GC info: %s\n", constants.MsgTimeoutGeneric)
				} else {
					fmt.Printf("Failed to get GC info: %v\n", err)
				}
				return
			}

			output := formatter.FormatGCTable(gcInfos)
			fmt.Print(output)
		}
	},
}

func init() {
	nodeCmd.AddCommand(nodeGCCmd)
	nodeGCCmd.Flags().String("name", "", "Show detailed GC info for specific node")
}
